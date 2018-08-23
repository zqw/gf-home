package ctlindex

import (
	"gitee.com/johng/gf/g/net/ghttp"
	"gitee.com/johng/gf/g"
    "fmt"
)

// 网站首页
func Index(r *ghttp.Request) {
    url := fmt.Sprintf("http://%s", g.Config().GetString("doc.domain"))
	r.Response.RedirectTo(url)
	//err := r.Response.Template("index.html", g.Map{
    //
	//})
	//fmt.Println(err)
}
