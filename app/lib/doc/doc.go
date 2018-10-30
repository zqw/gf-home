package libDoc

import (
    "gitee.com/johng/gf/g/os/gfile"
    "gopkg.in/russross/blackfriday.v2"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/util/gregex"
    "fmt"
    "strings"
    "gitee.com/johng/gf/g/util/gstr"
    "gitee.com/johng/gf/g/os/gfcache"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/os/gproc"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/util/gconv"
)

var (
    // 文档缓存
    cache = gcache.New()
)

// 更新doc版本库
func UpdateDocGit() {
    err := gproc.ShellRun(
        fmt.Sprintf(`cd %s && git pull origin master`, g.Config().GetString("doc.path")),
    )
    if err == nil {
        // 每次文档的更新都要清除缓存对象数据
        cache.Clear()

        glog.Cat("doc-hook").Printfln("doc hook updates")
    } else {
        glog.Cat("doc-hook").Printfln("doc hook updates error: %v",  err)
    }
}

// 根据关键字进行markdown文档搜索，返回文档path列表
func SearchMdByKey(key string) []string {
    glog.Cat("search").Println(key)
    v := cache.GetOrSetFunc("doc_search_result_" + key, func() interface{} {
        // 当该key的检索缓存不存在时，执行检索
        array    := garray.NewStringArray(0, 0, false)
        docPath  := g.Config().GetString("doc.path")
        paths    := cache.GetOrSetFunc("doc_files_recursive", func() interface{} {
            // 当目录列表不存在时，执行检索
            paths, _ := gfile.ScanDir(docPath, "*.md", true)
            return paths
        }, 0)
        // 遍历markdown文件列表，执行字符串搜索
        for _, path := range gconv.Strings(paths) {
            content := gfcache.GetContents(path)
            if len(content) > 0 {
                if strings.Index(content, key) != -1 {
                    index := gstr.Replace(path, ".md", "")
                    index  = gstr.Replace(index, docPath, "")
                    array.Append(index)
                }
            }
        }
        return array.Slice()
    }, 0)

    return gconv.Strings(v)
}

// 根据path参数获得层级显示的title
func GetTitleByPath(path string) string {
    v := cache.GetOrSetFunc("title_by_path_" + path, func() interface{} {
        type lineItem struct {
            indent int
            name   string
        }
        path        = strings.TrimLeft(path, "/")
        array      := make([]lineItem, 0)
        mdContent  := GetMarkdown("menus")
        lines      := strings.Split(mdContent, "\n")
        indent     := 0
        for _, line := range lines {
            match, _ := gregex.MatchString(`(\s*)\*\s+\[(.+)\]\((.+)\)`, line)
            if len(match) == 4 {
                item := lineItem{
                    indent : len(match[1]),
                    name   : match[2],
                }
                mdPath := gstr.Replace(match[3], ".md", "")
                if item.indent > indent || len(array) == 0 {
                    array = append(array, item)
                } else if len(match[1]) == indent {
                    array[len(array) - 1] = item
                } else {
                    newArray := make([]lineItem, 0)
                    for _, v := range array {
                        if v.indent < item.indent {
                            newArray = append(newArray, v)
                        }
                    }
                    newArray = append(newArray, item)
                    array    = newArray
                }
                indent = item.indent
                if mdPath == path {
                    break
                }
            }
        }
        if len(array) > 0 {
            title := ""
            for i := len(array) - 1; i >= 0; i-- {
                if len(title) > 0 {
                    title += " - " + array[i].name
                } else {
                    title  = array[i].name
                }
            }
            return title
        }
        return nil
    }, 0)
    if v != nil {
        return v.(string)
    }
    return ""
}

// 获得指定uri路径的markdown文件内容
func GetMarkdown(path string) string {
    mdRoot    := g.Config().GetString("doc.path")
    content   := gfcache.GetContents(mdRoot + gfile.Separator + path + ".md")
    pattern   := `\[(.*)\]\((.+?)\)`
    content, _ = gregex.ReplaceStringFunc(pattern, content, func(s string) string {
        match, _ := gregex.MatchString(pattern, s)
        if len(match) > 1 {
            url := match[2]
            // 替换为绝对路径
            if url[0] != '/' && url[0] != '#' && !strings.Contains(url, "://") {
                url = fmt.Sprintf(`/%s`, url)
            }
            // 去掉markdown连接的后缀名称
            if strings.EqualFold(gfile.Ext(url), ".md") {
                url = gstr.Replace(url, ".md", "")
            }
            return fmt.Sprintf(`[%s](%s)`, match[1], url)
        }
        return s
    })
    return content
}

// 获得解析为html的markdown文件内容
func GetParsed(path string) string {
    return ParseMarkdown(GetMarkdown(path))
}

// 解析markdown为html
func ParseMarkdown(content string) string {
    if content == "" {
        return ""
    }
    return string(blackfriday.Run([]byte(content)))
}
