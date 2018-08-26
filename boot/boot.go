package boot

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/glog"
)

// 用于应用初始化。
func init() {
    v := g.View()
    c := g.Config()
    s := g.Server("doc")

    // 配置对象及视图对象配置
    c.AddPath("config")
    v.AddPath("static/template")

    // glog配置
    logpath := c.GetString("setting.logpath")
    glog.SetPath(logpath)
    glog.SetStdPrint(true)

    // Web Server配置
    s.AddSearchPath(c.GetString("doc.path"))
    s.SetDenyRoutes([]string{
        "/config/*",
    })
    s.SetLogPath(logpath)
    s.SetErrorLogEnabled(true)
    s.SetAccessLogEnabled(true)
    s.SetPort(9999)
}

