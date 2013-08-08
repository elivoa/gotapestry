package syd

import (
	// pages & components' import
	// index_pages "syd/pages"
	// p_admin "syd/pages/admin"
	// p_api "syd/pages/api"
	// p_api_suggest "syd/pages/api/suggest"
	// person_pages "syd/pages/person"
	// product_pages "syd/pages/product"

	// pages_order "syd/pages/order"
	// pages_order_create "syd/pages/order/create"

	"got/register"
	// syd_components "syd/components"
	// layout_components "syd/components/layout"
	// order_components "syd/components/order"
	// c_product "syd/components/product"
	"got/config"
	"got/utils"
)

var SYDModule = &register.Module{
	Name:        "syd",
	BasePath:    utils.CurrentBasePath(),
	PackagePath: "syd",
	Description: "SYD Selling System Main module.",
	Register: func() {
		c := config.Config
		c.AddStaticResource("/pictures/", "/var/site/data/syd/pictures/")
		c.AddStaticResource("/static/", "../static/")

		// index_pages.Register()
		// person_pages.Register()
		// product_pages.Register()
		// p_api.Register()
		// p_api_suggest.Register()
		// p_admin.Register()
		// pages_order.Register()
		// pages_order_create.Register()

		// syd_components.Register()
		// layout_components.Register()
		// order_components.Register()
		// c_product.Register()
	},
}
