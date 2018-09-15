package ctldoc

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf-home/app/lib/doc"
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
    baseTitle := config.GetString("doc.title")
    title     := doc.GetTitleByPath(path)
    if title == "" {
        title = "404 NOT FOUND"
    }
    title += " - " + config.GetString("doc.title")
    // markdown内容
    mdMainContent       := doc.GetMarkdown(path)
    mdMainContentParsed := doc.ParseMarkdown(mdMainContent)
    r.Response.WriteTpl("doc/index.html", g.Map {
        "title"               : title,
        "baseTitle"           : baseTitle,
        "mdMenuContentParsed" : gview.HTML(doc.GetParsed("menus")),
        "mdMainContentParsed" : gview.HTML(mdMainContentParsed),
        "mdMainContent"       : gview.HTML(mdMainContent),
    })
}

// 文档更新hook
func UpdateHook(r *ghttp.Request) {
    raw    := r.GetRaw()
    j, err := gjson.DecodeToJson(raw)
    if err != nil {
        panic(err)
    }
    if j != nil && j.GetString("password") == g.Config().GetString("doc.hook") {
        doc.UpdateDocGit()
    }
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