package doc

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
        cache.Clear()
        glog.Cat("doc-hook").Printfln("doc hook updates")
    } else {
        glog.Cat("doc-hook").Printfln("doc hook updates error: %v",  err)
    }
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
    mdRoot  := g.Config().GetString("doc.path")
    content := gfcache.GetContents(mdRoot + gfile.Separator + path + ".md")
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
    content    = string(blackfriday.Run([]byte(content)))
    pattern   := `href="(.+?)"`
    content, _ = gregex.ReplaceStringFunc(pattern, content, func(s string) string {
        match, _ := gregex.MatchString(pattern, gstr.Replace(s, ".md", ""))
        if len(match) > 1 {
            if match[1][0] != '/' && match[1][0] != '#' && !strings.Contains(match[1], "://") {
                return fmt.Sprintf(`href="/%s"`, match[1])
            }
        }
        return s
    })
    return content
}
