package router

import (
    "fmt"
    "gitee.com/johng/gf-home/app/ctl/doc"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/util/gregex"
)

// 统一路由注册.
func init() {
    // 开发文档
    g.Server("doc").BindHandler("/*path",   ctlDoc.Index)
    g.Server("doc").BindHandler("/hook",    ctlDoc.UpdateHook)
    g.Server("doc").BindHandler("/search",  ctlDoc.Search)
    g.Server("doc").EnableAdmin("/admin")
    // 某些浏览器会直接请求/favicon.ico文件，会产生404
    g.Server("doc").BindHandler("/favicon.ico", func(r *ghttp.Request) {
        r.Response.ServeFile("/static/resource/image/favicon.ico")
    })
    // 为平滑重启管理页面设置HTTP Basic账号密码
    g.Server("doc").BindHookHandler("/admin/*", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
        user := g.Config().GetString("doc.admin.user")
        pass := g.Config().GetString("doc.admin.pass")
        if !r.BasicAuth(user, pass) {
            r.Exit()
        }
    })
    // 强制跳转到HTTPS访问
    g.Server("doc").BindHookHandler("/*", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
        if !r.IsFileServe() && r.TLS == nil {
            r.Response.RedirectTo(fmt.Sprintf("https://%s%s", r.Host, r.URL.String()))
            r.Exit()
        }
    })
    // 所有静态文件使用CDN加速
    g.Server("doc").BindHookHandler("/*", ghttp.HOOK_BEFORE_OUTPUT, func(r *ghttp.Request) {
        // 对所有动态内容执行替换
        if !r.IsFileServe() {
            pattern := `(src|href)=["'](\/.+\.(js|css|png|jpg|jpeg|gif|font|ico).*?)["']`
            b, _    := gregex.Replace(pattern,
                []byte(fmt.Sprintf(`$1="https://%s$2"`, g.Config().GetString("cdn.url"))),
                r.Response.Buffer(),
            )
            r.Response.SetBuffer(b)
        }
    })
}
