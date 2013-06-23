package syd

import (
	"fmt"
	"github.com/gorilla/mux"
	"got"
	"got/register"
	"got/templates"
	"syd/pages/ajax"
	"syd/pages/order"

	// pages & components' import
	index_pages "syd/pages"
	person_pages "syd/pages/person"
	product_pages "syd/pages/product"

	syd_components "syd/components"
)

func Start() {
	// welcome message
	fmt.Println("syd > SYD Sales Manage System Starting...")

	// prepare include templates.
	// Move to gotframework got.init
	prepareIncludeTemplates()
	// registerComponents()

	g := got.Init()
	g.Module(
		simpleModule,
		sydModule,
	)

	welcome()

	// start server
	g.StartServer(&got.GotConfig{
		StaticResources: []string{"/static/", "../"},
	})
}

// ________________________________________________________________________________
// Register SYD GOT style pages.
// TODO: How to automatically do this?
//
func sydModule(r *mux.Router) {
	register.RegisterApp(
		"syd",
		"SYD Module",
		"/Users/bogao/sync/sydPage/src/syd",
	)

	index_pages.Register()
	person_pages.Register()
	product_pages.Register()

	syd_components.Register()
}

// register simple router module into GOT.
func simpleModule(r *mux.Router) {
	// person.New().Mapping(r)
	// product.New().Mapping(r)
	order.New().Mapping(r)
	ajax.New().Mapping(r)
}

func welcome() {
	fmt.Println("\n")
	fmt.Println("``````````````````````````````````````````````````")
	fmt.Println("`  SYD Sale Management System (ALPHA 1)          `")
	fmt.Println("`                                                `")
	fmt.Println("``````````````````````````````````````````````````")
	fmt.Printf("Server Started, Listen localhost:8080\n\n")
	got.PrintRegistry()
}

func prepareIncludeTemplates() {
	templates.Add("include/header")
	templates.Add("include/left-nav")
}
