package ctldoc

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g"
    "fmt"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gview"
)

type Controller struct { }

func (c *Controller) Index(r *ghttp.Request) {
    path   := r.Get("path")
    mdRoot := g.Config().GetString("gf-doc.path")
    ext    := gfile.Ext(path)
    if ext != "" {
        r.Response.ServeFile(fmt.Sprintf("%s%s%s", mdRoot, gfile.Separator, path))
        return
    }
    mdMenuPath    := fmt.Sprintf("%s%smenus.md", mdRoot, gfile.Separator)
    mdMainPath    := fmt.Sprintf("%s%s%s.md",    mdRoot, gfile.Separator, path)
    mdMenuContent := gfile.GetContents(mdMenuPath)
    mdMainContent := ""
    if gfile.Exists(mdMainPath) {
        mdMainContent = gfile.GetContents(mdMainPath)
    }
    r.Response.Template("index.html", g.Map {
        "mdMenuContent" : gview.HTML(mdMenuContent),
        "mdMainContent" : gview.HTML(mdMainContent),
    })
}