package router

import (
    "github.com/gogf/gf-home/app/controller/document"
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/net/ghttp"
)

// 统一路由注册.
func init() {
    // 开发文档
    g.Server().BindHandler("/*path",    ctl_document.Index)
    g.Server().BindHandler("/hook",     ctl_document.UpdateHook)
    g.Server().BindHandler("/search",   ctl_document.Search)

    // 管理接口
    g.Server().EnableAdmin("/admin")

    // 某些浏览器会直接请求/favicon.ico文件，会产生404
    g.Server().SetRewrite("/favicon.ico", "/resource/image/favicon.ico")

    // 为平滑重启管理页面设置HTTP Basic账号密码
    g.Server().BindHookHandler("/admin/*", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
        user := g.Config().GetString("admin.user")
        pass := g.Config().GetString("admin.pass")
        if !r.BasicAuth(user, pass) {
            r.Exit()
        }
    })

    // 强制跳转到HTTPS访问
    //g.Server().BindHookHandler("/*", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
    //  if !r.IsFileRequest() && r.TLS == nil {
    //      r.Response.RedirectTo(fmt.Sprintf("https://%s%s", r.Host, r.URL.String()))
    //  }
    //})
}
