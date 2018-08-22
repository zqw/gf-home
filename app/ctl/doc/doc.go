package ctldoc

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g"
    "fmt"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf-home/app/lib/doc"
    "gitee.com/johng/gf/g/util/gregex"
)

type Controller struct { }


func init() {
    g.Server().BindObjectMethod("/*path", new(Controller), "Index")
}

// 文档首页
func (c *Controller) Index(r *ghttp.Request) {
    if r.IsAjaxRequest() {
        c.serveMarkdownAjax(r)
        return
    }

    path := r.Get("path")
    if path == "" {
        r.Response.RedirectTo("/index")
        return
    }
    config := g.Config()
    mdRoot := config.GetString("gf-doc.path")
    ext    := gfile.Ext(path)
    if ext != "" && ext != "md" {
        r.Response.ServeFile(fmt.Sprintf("%s%s%s", mdRoot, gfile.Separator, path))
        return
    }
    baseTitle    := config.GetString("gf-doc.title")
    title        := baseTitle
    menuMarkdown := doc.GetMarkdown("menus")
    match, _     := gregex.MatchString(fmt.Sprintf(`\[(.+)\]\(%s\)`, path), menuMarkdown)
    if len(match) > 1 {
        title = fmt.Sprintf("%s - %s", match[1], baseTitle)
    }
    r.Response.Template("index.html", g.Map {
        "title"               : title,
        "baseTitle"           : baseTitle,
        "mdMenuContentParsed" : gview.HTML(doc.ParseMarkdown(menuMarkdown)),
        "mdMainContentParsed" : gview.HTML(doc.GetParsed(path)),
        "mdMainContent"       : gview.HTML(doc.GetMarkdown(path)),
    })
}

// 处理ajax请求
func (c *Controller) serveMarkdownAjax(r *ghttp.Request) {
    r.Response.WriteJson(g.Map{
        "code" : 1,
        "msg"  : "",
        "data" : doc.GetMarkdown(r.Get("path", "index")),
    })
}