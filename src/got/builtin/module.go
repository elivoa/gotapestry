/*
GOT Module: builtin package
*/

package builtin

import (

	// // pages import
	// p_root "got/builtin/pages"
	// p_builtin "got/builtin/pages/got"

	// // components import
	// c_builtin "got/builtin/components"

	"got/builtin/pages/got/fileupload"
	"got/register"
	"got/utils"
	"net/http"
)

var BuiltinModule = &register.Module{
	Name:        "got/builtin",
	BasePath:    utils.CurrentBasePath(),
	PackagePath: "got/builtin",
	Description: "GOT Framework Built-in pages and components etc.",
	// some special configuration.
	Register: func() {
		// // pages
		// p_root.Register()
		// p_builtin.Register()
		// p_fileupload.Register()

		// // components
		// c_builtin.Register()

		//
		// *** very special:: file upload *** TODO make this beautiful.
		// Special mapping, all file upload maps here
		//
		http.HandleFunc("/got/fileupload/", fileupload.FU)
	},
}
