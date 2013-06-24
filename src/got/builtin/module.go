/*
GOT Module: builtin package
*/

package builtin

import (
	"github.com/gorilla/mux"
	"got/register"

	// pages import
	root_pages "got/builtin/pages"
	builtin_pages "got/builtin/pages/got"

	// components import
	builtin_components "got/builtin/components"
)

import ()

func GotBuiltinModule(r *mux.Router) {

	register.RegisterApp(
		"got/builtin",
		"GOT Core Modules",
		"/Users/bogao/develop/gitme/gotapestry/src/got/builtin",
	)

	// pages
	root_pages.Register()
	builtin_pages.Register()

	// components
	builtin_components.Register()
}
