package main

import (
    "gitee.com/johng/gf/g"
    _ "gitee.com/johng/gf-home/boot"
    _ "gitee.com/johng/gf-home/router"
)

func main() {
    g.Server("doc").Run()
}