package main

import (
    _ "gitee.com/johng/gf-home/boot"
    _ "gitee.com/johng/gf-home/router"
    "gitee.com/johng/gf/g"
)

func main() {
    g.Server().Run()
}