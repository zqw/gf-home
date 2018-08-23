package ctlindex

import (
	"gitee.com/johng/gf/g"
	"gitee.com/johng/gf/g/net/ghttp"
	"fmt"
)


// 网站首页
func Index(r *ghttp.Request) {
	err := r.Response.Template("index.html", g.Map{

	})
	fmt.Println(err)
}
