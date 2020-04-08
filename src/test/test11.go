package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func main() {
	module_basePash := "/Users/bogao/develop/syd/v1/got/builtin"
	module_PackageName := "github.com/elivoa/got/builtin"
	ip := "github.com/elivoa/got/builtin/pages/got"
	a := ""
	if strings.HasPrefix(ip, module_PackageName) {
		a = ip[len(module_PackageName):]
	}
	fmt.Println(a)
	fmt.Println(filepath.Join(
		module_basePash,
		a,
		"Status.html",
		// s.StructInfo.ImportPath,
		// fmt.Sprintf("%v%v", s.StructInfo.StructName, conf.TemplateFileExtension),
	))
}
