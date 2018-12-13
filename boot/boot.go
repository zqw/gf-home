package boot

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/glog"
)

// 用于配置初始化.
func init() {
    v := g.View()
    c := g.Config()
    s := g.Server()

    // 配置及视图对象
    c.AddPath("config")
    v.AddPath("template")

    // 日志模块配置
    logpath := c.GetString("logpath")
    glog.SetPath(logpath)
    glog.SetStdPrint(true)

    // Web Server配置
    s.SetLogPath(logpath)
    s.SetErrorLogEnabled(true)
    s.SetAccessLogEnabled(true)
    if c.Get("ssl") != nil {
        s.EnableHTTPS(c.GetString("ssl.crt"), c.GetString("ssl.key"))
        s.SetHTTPSPort(c.GetInt("https-port"))
    }
    s.SetServerRoot("./public")
    s.AddSearchPath(c.GetString("document.path"))
    s.SetPort(c.GetInt("http-port"))
}

