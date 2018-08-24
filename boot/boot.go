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
    c.AddPath("config")
    v.AddPath("static/template")

    logpath := c.GetString("setting.logpath")

    glog.SetPath(logpath)
    glog.SetStdPrint(true)

    s.SetDenyRoutes([]string{
        "/config/*",
    })
    s.SetLogPath(logpath)
    s.SetErrorLogEnabled(true)
    s.SetAccessLogEnabled(true)
    s.SetPort(9999)
}

