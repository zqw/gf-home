package ctldoc

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g"
    "fmt"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf-home/app/lib/doc"
    "gitee.com/johng/gf/g/util/gregex"
    "gitee.com/johng/gf/g/os/gproc"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/encoding/gjson"
    "net/http"
)

// 文档首页
func Index(r *ghttp.Request) {
    if r.IsAjaxRequest() {
        serveMarkdownAjax(r)
        return
    }
    path := r.Get("path")
    if path == "" {
        r.Response.RedirectTo("/index")
        return
    }
    config := g.Config()
    // 如果是静态文件请求，那么表示Web Server没有找到该文件，那么直接404，本接口不支持待后缀的静态文件处理。
    // 由于路由规则比较宽，这里也会有未存在的静态文件请求匹配进来。
    if gfile.Ext(path) != "" {
        r.Response.WriteStatus(http.StatusNotFound)
        return
    }
    // 菜单内容
    baseTitle    := config.GetString("doc.title")
    title        := baseTitle
    menuMarkdown := doc.GetMarkdown("menus")
    fmt.Println(path)
    match, _     := gregex.MatchString(fmt.Sprintf(`\[(.+)\]\(%s\.md\)`, path), menuMarkdown)
    if len(match) > 1 {
        title = fmt.Sprintf("%s - %s", match[1], baseTitle)
    } else {
        title = fmt.Sprintf("404 NOT FOUND - %s", baseTitle)
    }
    // markdown内容
    mdMainContent       := doc.GetMarkdown(path)
    mdMainContentParsed := doc.ParseMarkdown(mdMainContent)
    r.Response.Template("doc/index.html", g.Map {
        "title"               : title,
        "baseTitle"           : baseTitle,
        "mdMenuContentParsed" : gview.HTML(doc.ParseMarkdown(menuMarkdown)),
        "mdMainContentParsed" : gview.HTML(mdMainContentParsed),
        "mdMainContent"       : gview.HTML(mdMainContent),
    })
}

// 文档更新hook
func UpdateHook(r *ghttp.Request) {
    raw    := r.GetRaw()
    j, err := gjson.DecodeToJson(raw)
    if j != nil && j.GetString("password") == g.Config().GetString("doc.hook") {
        err = gproc.ShellRun(
            fmt.Sprintf(`cd %s && git pull origin master`, g.Config().GetString("doc.path")),
        )
    }
    glog.Cat("doc-hook").Printfln("doc hook update from: %s, error: %v, content: %s", r.URL.String(), err, string(raw))
    r.Response.Write("ok")
}

// 处理ajax请求
func serveMarkdownAjax(r *ghttp.Request) {
    r.Response.WriteJson(g.Map{
        "code" : 1,
        "msg"  : "",
        "data" : doc.GetMarkdown(r.Get("path", "index")),
    })
}