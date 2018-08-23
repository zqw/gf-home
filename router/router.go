package router

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf-home/app/ctl/doc"
    "gitee.com/johng/gf-home/app/ctl/index"
)

// 统一路由注册
func init() {
    config := g.Config()
    server := g.Server()

    // 开发文档
    server.Domain(config.GetString("doc.domain")).BindHandler("/*path", ctldoc.Index)

    // 官网首页
    server.Domain(config.GetString("index.domain")).BindHandler("/", ctlindex.Index)
}