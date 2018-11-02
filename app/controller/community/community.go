package community

import (
    "gitee.com/johng/gf/g/net/ghttp"
)

type Community struct {

}

func (c *Community) Index(r *ghttp.Request) {
    r.Response.WriteTpl("community/index.html", nil)
}