module gf-home-main

go 1.12

require (
	blackfriday v0.0.0
	github.com/gogf/gf v1.8.0
	github.com/gogf/gf-home v0.0.0
	github.com/shurcooL/sanitized_anchor_name v1.0.0
)

replace github.com/gogf/gf-home => ../gf-home

replace blackfriday => ../gf-home/jar_package/gopkg.in/russross/blackfriday.v2
