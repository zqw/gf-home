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
    mdRoot := config.GetString("doc.path")
    ext    := gfile.Ext(path)
    if ext != "" && ext != "md" {
        r.Response.ServeFile(fmt.Sprintf("%s%s%s", mdRoot, gfile.Separator, path))
        return
    }
    baseTitle    := config.GetString("doc.title")
    title        := baseTitle
    menuMarkdown := doc.GetMarkdown("menus")
    match, _     := gregex.MatchString(fmt.Sprintf(`\[(.+)\]\(%s\)`, path), menuMarkdown)
    if len(match) > 1 {
        title = fmt.Sprintf("%s - %s", match[1], baseTitle)
    }
    r.Response.Template("doc.html", g.Map {
        "title"               : title,
        "baseTitle"           : baseTitle,
        "mdMenuContentParsed" : gview.HTML(doc.ParseMarkdown(menuMarkdown)),
        "mdMainContentParsed" : gview.HTML(doc.GetParsed(path)),
        "mdMainContent"       : gview.HTML(doc.GetMarkdown(path)),
    })
}

// 文档更新hook
func UpdateHook(r *ghttp.Request) {
    j := r.GetJson()
    if j != nil && j.GetString("password") == g.Config().GetString("doc.hook") {
        gproc.ShellRun(
            fmt.Sprintf(`cd %s && git pull origin master`, g.Config().GetString("doc.path")),
        )
    }
}

// 处理ajax请求
func serveMarkdownAjax(r *ghttp.Request) {
    r.Response.WriteJson(g.Map{
        "code" : 1,
        "msg"  : "",
        "data" : doc.GetMarkdown(r.Get("path", "index")),
    })
}