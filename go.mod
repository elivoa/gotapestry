module gotapestry

go 1.14

require (
	github.com/elivoa/got v0.0.0
	github.com/elivoa/gxl v0.0.0
	syd v0.0.0
)

replace (
	github.com/elivoa/got => ../got
	github.com/elivoa/gxl => ../gxl
	syd => ./src/syd
)
