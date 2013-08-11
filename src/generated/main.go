// DO NOT EDIT THIS FILE -- GENERATED CODE
package main

import (
    "fmt"
    _got "got"
    "got/config"
    "got/register"
    "got/cache"
    "got/parser"
    "got/route"
	
    builtin "got/builtin"
    components0 "got/builtin/components"
    got "got/builtin/pages"
    got0 "got/builtin/pages/got"
    fileupload "got/builtin/pages/got/fileupload"
    syd "syd"
    components "syd/components"
    layout "syd/components/layout"
    order "syd/components/order"
    product "syd/components/product"
    index "syd/pages"
    admin "syd/pages/admin"
    api "syd/pages/api"
    api0 "syd/pages/api/suggest"
    order0 "syd/pages/order"
    order1 "syd/pages/order/create"
    person "syd/pages/person"
    product0 "syd/pages/product"
)

func main() {
    fmt.Println("=============== STARTING ================================================")

    // restore config.ModulePath
    
    config.Config.RegisterModulePath("/Users/bogao/develop/gitme/gotapestry/src/syd", "SYDModule")
    config.Config.RegisterModulePath("/Users/bogao/develop/gitme/gotapestry/src/got/builtin", "BuiltinModule")

    // parse source again.
    sourceInfo, compileError := parser.ParseSource(config.Config.ModulePath, false) // deep parse
    if compileError != nil {
        panic(compileError.Error())
    }

	// cache source info into cache.
    cache.SourceCache = sourceInfo

    // register real module
    register.RegisterModule(
    
      syd.SYDModule,
      builtin.BuiltinModule,
    )

    // register pages & components
    
    route.RegisterProton("syd/components", "OrderDetailsEditor", "syd", &components.OrderDetailsEditor{})
    route.RegisterProton("syd/components", "SuggestControl", "syd", &components.SuggestControl{})
    route.RegisterProton("syd/components/layout", "Header", "syd", &layout.Header{})
    route.RegisterProton("syd/components/layout", "HeaderNav", "syd", &layout.HeaderNav{})
    route.RegisterProton("syd/components/layout", "LeftNav", "syd", &layout.LeftNav{})
    route.RegisterProton("syd/components/order", "BatchCloseOrder", "syd", &order.BatchCloseOrder{})
    route.RegisterProton("syd/components/order", "CustomerInfo", "syd", &order.CustomerInfo{})
    route.RegisterProton("syd/components/order", "OrderCloseButton", "syd", &order.OrderCloseButton{})
    route.RegisterProton("syd/components/order", "OrderDeliverButton", "syd", &order.OrderDeliverButton{})
    route.RegisterProton("syd/components/order", "OrderDetailsForm", "syd", &order.OrderDetailsForm{})
    route.RegisterProton("syd/components/order", "OrderList", "syd", &order.OrderList{})
    route.RegisterProton("syd/components/order", "OrderProductSelector", "syd", &order.OrderProductSelector{})
    route.RegisterProton("syd/components/product", "ProductColorSizeTable", "syd", &product.ProductColorSizeTable{})
    route.RegisterProton("syd/components/product", "ProductList", "syd", &product.ProductList{})
    route.RegisterProton("syd/pages", "Index", "syd", &index.Index{})
    route.RegisterProton("syd/pages/admin", "AdminIndex", "syd", &admin.AdminIndex{})
    route.RegisterProton("syd/pages/api", "Api", "syd", &api.Api{})
    route.RegisterProton("syd/pages/api/suggest", "Suggest", "syd", &api0.Suggest{})
    route.RegisterProton("syd/pages/order", "DeliveringUnclosedOrders", "syd", &order0.DeliveringUnclosedOrders{})
    route.RegisterProton("syd/pages/order", "OrderIndex", "syd", &order0.OrderIndex{})
    route.RegisterProton("syd/pages/order", "OrderList", "syd", &order0.OrderList{})
    route.RegisterProton("syd/pages/order", "ShippingInsteadList", "syd", &order0.ShippingInsteadList{})
    route.RegisterProton("syd/pages/order", "ButtonSubmitHere", "syd", &order0.ButtonSubmitHere{})
    route.RegisterProton("syd/pages/order", "ViewOrder", "syd", &order0.ViewOrder{})
    route.RegisterProton("syd/pages/order", "OrderPrint", "syd", &order0.OrderPrint{})
    route.RegisterProton("syd/pages/order", "PrintExpressYTO", "syd", &order0.PrintExpressYTO{})
    route.RegisterProton("syd/pages/order", "ShippingInsteadPrint", "syd", &order0.ShippingInsteadPrint{})
    route.RegisterProton("syd/pages/order/create", "OrderCreateDetail", "syd", &order1.OrderCreateDetail{})
    route.RegisterProton("syd/pages/order/create", "OrderCreateIndex", "syd", &order1.OrderCreateIndex{})
    route.RegisterProton("syd/pages/order/create", "ShippingInstead", "syd", &order1.ShippingInstead{})
    route.RegisterProton("syd/pages/person", "PersonIndex", "syd", &person.PersonIndex{})
    route.RegisterProton("syd/pages/person", "PersonList", "syd", &person.PersonList{})
    route.RegisterProton("syd/pages/person", "PersonEdit", "syd", &person.PersonEdit{})
    route.RegisterProton("syd/pages/person", "PersonDetail", "syd", &person.PersonDetail{})
    route.RegisterProton("syd/pages/product", "ProductIndex", "syd", &product0.ProductIndex{})
    route.RegisterProton("syd/pages/product", "ProductCreate", "syd", &product0.ProductCreate{})
    route.RegisterProton("syd/pages/product", "ProductEdit", "syd", &product0.ProductEdit{})
    route.RegisterProton("syd/pages/product", "ProductList", "syd", &product0.ProductList{})
    route.RegisterProton("syd/pages/product", "ProductDetail", "syd", &product0.ProductDetail{})
    route.RegisterProton("got/builtin/components", "FileUpload", "got/builtin", &components0.FileUpload{})
    route.RegisterProton("got/builtin/components", "Output", "got/builtin", &components0.Output{})
    route.RegisterProton("got/builtin/components", "ProvinceSelect", "got/builtin", &components0.ProvinceSelect{})
    route.RegisterProton("got/builtin/components", "Select", "got/builtin", &components0.Select{})
    route.RegisterProton("got/builtin/pages", "Errors", "got/builtin", &got.Errors{})
    route.RegisterProton("got/builtin/pages/got", "Status", "got/builtin", &got0.Status{})
    route.RegisterProton("got/builtin/pages/got", "TestIndex", "got/builtin", &got0.TestIndex{})
    route.RegisterProton("got/builtin/pages/got/fileupload", "FileUploadTest", "got/builtin", &fileupload.FileUploadTest{})
    route.RegisterProton("got/builtin/pages/got/fileupload", "FileUploadIndex", "got/builtin", &fileupload.FileUploadIndex{})

    // start the server
    _got.Start()
}
