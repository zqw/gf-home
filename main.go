package main

import (
    "gitee.com/johng/gf/g"
    _ "gitee.com/johng/gf-home/app/ctl/doc"
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

    s.SetPort(8199)
    s.Run()
}