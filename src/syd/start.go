package syd

import (
	// pages & components' import
	index_pages "syd/pages"
	api_pages "syd/pages/api"
	person_pages "syd/pages/person"
	product_pages "syd/pages/product"

	pages_order "syd/pages/order"
	pages_order_create "syd/pages/order/create"

	syd_components "syd/components"
	layout_components "syd/components/layout"
	order_components "syd/components/order"

	"fmt"
	"github.com/gorilla/mux"
	"got"
	"got/config"
	"got/register"
	"syd/pages/ajax"
)

func Start() {
	// welcome message
	fmt.Println("syd > SYD Sales Manage System Starting...")

	g := got.Init()
	g.Module(
		simpleModule,
		sydModule,
	)

	welcome()

	// start server
	config.Config.ResourcePath = "/var/site/data/syd/pictures/"
	g.StartServer(&got.GotConfig{
		StaticResources: [][]string{
			[]string{"/pictures/", "/var/site/data/syd/pictures/"},
			[]string{"/static/", "../static/"},
		},
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
		"syd",
	)

	index_pages.Register()
	person_pages.Register()
	product_pages.Register()
	api_pages.Register()

	pages_order.Register()
	pages_order_create.Register()

	syd_components.Register()
	layout_components.Register()
	order_components.Register()
}

// register simple router module into GOT.
func simpleModule(r *mux.Router) {
	// person.New().Mapping(r)
	// product.New().Mapping(r)
	// order.New().Mapping(r)
	ajax.New().Mapping(r)
}

func welcome() {
	fmt.Println("\n")
	fmt.Println("``````````````````````````````````````````````````")
	fmt.Println("`  SYD Sale Management System (ALPHA 1)          `")
	fmt.Println("`                                                `")
	fmt.Println("``````````````````````````````````````````````````")
	fmt.Printf("Server Started, Listen localhost:8080\n\n")
	// got.PrintRegistry()
}
