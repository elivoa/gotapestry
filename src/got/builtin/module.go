/*
GOT Module: builtin package
*/

package builtin

import (
	"github.com/gorilla/mux"
	"got/register"

	// pages import
	p_root "got/builtin/pages"
	p_builtin "got/builtin/pages/got"
	p_fileupload "got/builtin/pages/got/fileupload"

	// components import
	c_builtin "got/builtin/components"
	"got/core"
)

// todo change this.
func GotBuiltinModule(r *mux.Router) {

	//
	// register core builtin components and pages
	//
	register.RegisterApp(
		"got/builtin",      // app name
		"GOT Core Modules", // app description
		"got/builtin",      // app path related
	)

	// pages
	p_root.Register()
	p_builtin.Register()
	p_fileupload.Register()

	// components
	c_builtin.Register()

	//
	// *** very special:: file upload *** TODO make this beautiful.
	// Special mapping, all file upload maps here
	//
	r.HandleFunc("/got/fileupload/", p_fileupload.FU)

}

// --------------------------------------------------------------------------------
// TODO Make this presents like this:

type BuiltinModule struct {
	core.Module
}

func (m *BuiltinModule) Pages() {
	// TODO ...
}

func (m *BuiltinModule) Components() {
	// TODO ...
}

func (m *BuiltinModule) Mixins() {
	// TODO ...
}
