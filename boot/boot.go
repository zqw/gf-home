package boot

import (
    "gitee.com/johng/gf/g"
    "fmt"
)

// 初始化
func init() {
    g.Config().AddPath("config")

    g.View().AddPath("static/template")
    g.View().BindFunc("Config", funcConfig)
}

// 模板内置方法：include
func funcConfig(pattern string, file...string) interface{} {
    fmt.Println(pattern)
    return g.Config().Get(pattern, file...)
}

