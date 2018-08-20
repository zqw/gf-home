package ctldoc

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g"
    "fmt"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf-home/app/lib/doc"
)

type Controller struct { }

func (c *Controller) Index(r *ghttp.Request) {
    path := r.Get("path")
    if path == "" {
        r.Response.RedirectTo("/doc/index")
        return
    }
    mdRoot := g.Config().GetString("gf-doc.path")
    ext    := gfile.Ext(path)
    if ext != "" && ext != "md" {
        r.Response.ServeFile(fmt.Sprintf("%s%s%s", mdRoot, gfile.Separator, path))
        return
    }

    r.Response.Template("index.html", g.Map {
        "mdMenuContentParsed" : gview.HTML(doc.GetParsed("menus")),
        "mdMainContentParsed" : gview.HTML(doc.GetParsed(path)),
        "mdMainContent"       : gview.HTML(doc.GetMarkdown(path)),
    })
}