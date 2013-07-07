package got

import (
	"fmt"
	"github.com/gorilla/mux"
	"got/builtin"
	"got/register"
	"got/route"
	"net/http"
)

type GOT struct {
	r *mux.Router
}

func Init() *GOT {
	g := GOT{
		r: mux.NewRouter(),
	}
	g.Module(builtin.GotBuiltinModule) // init tapestry builtin module.
	return &g
}

func (g *GOT) Module(modules ...func(*mux.Router)) *GOT {
	for _, module := range modules {
		module(g.r) // call register module
	}
	return g
}

// Start the server
func (g *GOT) StartServer(config *GotConfig) {

	// mapping static files
	if config.StaticResources != nil {
		for _, staticResource := range config.StaticResources {
			g.r.PathPrefix(staticResource[0]).
				Handler(http.FileServer(http.Dir(staticResource[1])))
		}
	}

	// got url matcher
	g.r.HandleFunc("/{url:.*}", route.RouteHandler)

	// bind port and start server
	http.Handle("/", g.r)
	http.ListenAndServe(":8080", nil)
}

// ________________________________________________________________________________
// Init GOT framework builtin module.
// TODO: automatically get this address
//
func PrintRegistry() {
	register.Apps.PrintALL()

	fmt.Println("\n---- Pages ---------------------")
	register.Pages.PrintALL()

	fmt.Println("\n---- Components ---------------------")
	register.Components.PrintALL()

	fmt.Println("\n---- Mixins ---------------------")
	fmt.Println("... no mixins avaliable ...")

	fmt.Println("--------------------------------------------------------------------------------\n")
}

// ________________________________________________________________________________
// Got start configs

type GotConfig struct {
	StaticResources [][]string // e.g.: [["/static/", "../"], ...]
}
