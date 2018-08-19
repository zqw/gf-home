package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf-home/app/ctl/doc"
)

func init() {
    g.View().AddPath("static/template")
    g.Config().AddPath("config")
}

func main() {
    s := g.Server()
    s.SetDenyRoutes([]string{
        "/config/*",
        // "/static/template/*",
    })

    // 我们可以将所有的路由注册放到这里执行
    s.BindObjectMethod("/doc/*path", new(ctldoc.Controller), "Index")

    s.SetPort(8199)
    s.Run()
}