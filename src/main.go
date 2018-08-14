package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func init() {
    g.View().AddPath("static/template")
}

func main() {
    s := g.Server()
    s.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Template("index.html")
    })
    s.SetPort(8199)
    s.Run()
}