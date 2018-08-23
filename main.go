package main

import (
    "gitee.com/johng/gf/g"
    _ "gitee.com/johng/gf-home/boot"
    _ "gitee.com/johng/gf-home/router"
)

func main() {
    s := g.Server()
    s.SetDenyRoutes([]string{
        "/config/*",
        // "/static/template/*",
    })

    s.SetPort(8199)
    s.Run()
}