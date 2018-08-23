package router

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf-home/app/ctl/doc"
)

// 统一路由注册
func init() {
    // 开发文档
    g.Server("doc").BindHandler("/*path", ctldoc.Index)
}
