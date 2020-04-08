module syd

go 1.14

require (
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394
	github.com/elivoa/got v0.0.0
	github.com/elivoa/gxl v0.0.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/mozillazg/go-pinyin v0.16.0
	github.com/stretchr/testify v1.5.1
)

replace (
	github.com/elivoa/got => ../../../got
	github.com/elivoa/gxl => ../../../gxl
)
