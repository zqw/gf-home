package boot

import (
    "gitee.com/johng/gf/g"
)

// 初始化
func init() {
    g.Config().AddPath("config")

    g.View().AddPath("static/template")
}

