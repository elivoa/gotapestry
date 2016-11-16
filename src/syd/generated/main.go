// DO NOT EDIT THIS FILE -- GENERATED CODE
package main

import (
    "github.com/elivoa/got/config"
    "github.com/elivoa/got/parser"
    "github.com/elivoa/got/route"
    "fmt"
    _got "github.com/elivoa/got"
    "github.com/elivoa/got/register"
	"github.com/elivoa/got/cache"
	
    builtin "github.com/elivoa/got/builtin"
    components0 "github.com/elivoa/got/builtin/components"
    got "github.com/elivoa/got/builtin/components/got"
    layout0 "github.com/elivoa/got/builtin/components/layout"
    ui "github.com/elivoa/got/builtin/components/ui"
    got0 "github.com/elivoa/got/builtin/pages"
    got1 "github.com/elivoa/got/builtin/pages/got"
    fileupload "github.com/elivoa/got/builtin/pages/got/fileupload"
    syd "syd"
    components "syd/components"
    admin "syd/components/admin"
    client "syd/components/client"
    inventory "syd/components/inventory"
    layout "syd/components/layout"
    order "syd/components/order"
    person "syd/components/person"
    inventory0 "syd/components/placeorder"
    product "syd/components/product"
    stat "syd/components/stat"
    test "syd/components/test"
    index "syd/pages"
    account "syd/pages/account"
    account0 "syd/pages/account-old"
    test0 "syd/pages/accounting"
    admin0 "syd/pages/admin"
    preference "syd/pages/admin/preference"
    api "syd/pages/api"
    api0 "syd/pages/api/suggest"
    client0 "syd/pages/client"
    test1 "syd/pages/hidden"
    inventory1 "syd/pages/inventory"
    order0 "syd/pages/order"
    order1 "syd/pages/order/create"
    person0 "syd/pages/person"
    inventory2 "syd/pages/placeorder"
    product0 "syd/pages/product"
    api1 "syd/pages/service/suggest"
    stat0 "syd/pages/stat"
    test2 "syd/pages/test"
)

func main() {
    fmt.Println("\n=============== STARTING GENERATED CODE ================================================")

    // setup config.ModulePath
    
    config.Config.RegisterModule(syd.SYDModule)
    config.Config.RegisterModule(builtin.BuiltinModule)

    // parse source again.
    sourceInfo, compileError := parser.ParseSource(config.Config.Modules, false) // deep parse
    if compileError != nil {
        panic(compileError.Error())
    }

    // The first cache is runtime system's cache.
    // Important things are put into sourceInfo.
    // TODO: Put everything into SourceCache
    cache.SourceCache = sourceInfo

    // register real module
    register.RegisterModule(syd.SYDModule,builtin.BuiltinModule,)

    // register pages & components
    
    route.RegisterProton("syd/components", "OrderDetailsEditor", "syd", &components.OrderDetailsEditor{})
    route.RegisterProton("syd/components", "SuggestControl", "syd", &components.SuggestControl{})
    route.RegisterProton("syd/components/admin", "ConstList", "syd", &admin.ConstList{})
    route.RegisterProton("syd/components/client", "ClientOrderDetail", "syd", &client.ClientOrderDetail{})
    route.RegisterProton("syd/components/client", "ClientOrderList", "syd", &client.ClientOrderList{})
    route.RegisterProton("syd/components/inventory", "InventoryGroupList", "syd", &inventory.InventoryGroupList{})
    route.RegisterProton("syd/components/inventory", "InventoryList", "syd", &inventory.InventoryList{})
    route.RegisterProton("syd/components/inventory", "InventoryProductSelector", "syd", &inventory.InventoryProductSelector{})
    route.RegisterProton("syd/components/layout", "Header", "syd", &layout.Header{})
    route.RegisterProton("syd/components/layout", "HeaderNav", "syd", &layout.HeaderNav{})
    route.RegisterProton("syd/components/layout", "LeftNav", "syd", &layout.LeftNav{})
    route.RegisterProton("syd/components/order", "OrderCloseButton", "syd", &order.OrderCloseButton{})
    route.RegisterProton("syd/components/order", "OrderDeliverButton", "syd", &order.OrderDeliverButton{})
    route.RegisterProton("syd/components/order", "OrderDetailsForm", "syd", &order.OrderDetailsForm{})
    route.RegisterProton("syd/components/order", "OrderList", "syd", &order.OrderList{})
    route.RegisterProton("syd/components/order", "OrderProductSelector", "syd", &order.OrderProductSelector{})
    route.RegisterProton("syd/components/order", "BatchCloseOrder", "syd", &order.BatchCloseOrder{})
    route.RegisterProton("syd/components/order", "CustomerInfo", "syd", &order.CustomerInfo{})
    route.RegisterProton("syd/components/person", "CustomerProfileCard", "syd", &person.CustomerProfileCard{})
    route.RegisterProton("syd/components/person", "DebtCustomerList", "syd", &person.DebtCustomerList{})
    route.RegisterProton("syd/components/placeorder", "PlaceOrderList", "syd", &inventory0.PlaceOrderList{})
    route.RegisterProton("syd/components/product", "ProductList", "syd", &product.ProductList{})
    route.RegisterProton("syd/components/product", "ProductSalesChart", "syd", &product.ProductSalesChart{})
    route.RegisterProton("syd/components/product", "ProductCalcTable", "syd", &product.ProductCalcTable{})
    route.RegisterProton("syd/components/product", "ProductColorSizeTable", "syd", &product.ProductColorSizeTable{})
    route.RegisterProton("syd/components/stat", "HotSaleProduct2", "syd", &stat.HotSaleProduct2{})
    route.RegisterProton("syd/components/stat", "TodayStat", "syd", &stat.TodayStat{})
    route.RegisterProton("syd/components/stat", "TrendStat", "syd", &stat.TrendStat{})
    route.RegisterProton("syd/components/stat", "HotSaleProduct", "syd", &stat.HotSaleProduct{})
    route.RegisterProton("syd/components/test", "TestPager", "syd", &test.TestPager{})
    route.RegisterProton("syd/pages", "Index", "syd", &index.Index{})
    route.RegisterProton("syd/pages/account", "AccountChangePassword", "syd", &account.AccountChangePassword{})
    route.RegisterProton("syd/pages/account", "AccountIndex", "syd", &account.AccountIndex{})
    route.RegisterProton("syd/pages/account", "AccountLogin", "syd", &account.AccountLogin{})
    route.RegisterProton("syd/pages/account", "AccountLogout", "syd", &account.AccountLogout{})
    route.RegisterProton("syd/pages/account", "AccountRegister", "syd", &account.AccountRegister{})
    route.RegisterProton("syd/pages/account-old", "AccountIndex", "syd", &account0.AccountIndex{})
    route.RegisterProton("syd/pages/account-old", "AccountLogin", "syd", &account0.AccountLogin{})
    route.RegisterProton("syd/pages/account-old", "AccountLogout", "syd", &account0.AccountLogout{})
    route.RegisterProton("syd/pages/account-old", "AccountRegister", "syd", &account0.AccountRegister{})
    route.RegisterProton("syd/pages/accounting", "FactorySettleAccount", "syd", &test0.FactorySettleAccount{})
    route.RegisterProton("syd/pages/admin", "AdminIndex", "syd", &admin0.AdminIndex{})
    route.RegisterProton("syd/pages/admin/preference", "PreferenceIndex", "syd", &preference.PreferenceIndex{})
    route.RegisterProton("syd/pages/api", "Api", "syd", &api.Api{})
    route.RegisterProton("syd/pages/api/suggest", "Suggest", "syd", &api0.Suggest{})
    route.RegisterProton("syd/pages/client", "ClientOrderIndex", "syd", &client0.ClientOrderIndex{})
    route.RegisterProton("syd/pages/hidden", "Hidden", "syd", &test1.Hidden{})
    route.RegisterProton("syd/pages/inventory", "InventoryEdit", "syd", &inventory1.InventoryEdit{})
    route.RegisterProton("syd/pages/inventory", "InventoryIndex", "syd", &inventory1.InventoryIndex{})
    route.RegisterProton("syd/pages/inventory", "InventoryList", "syd", &inventory1.InventoryList{})
    route.RegisterProton("syd/pages/inventory", "InventoryStat", "syd", &inventory1.InventoryStat{})
    route.RegisterProton("syd/pages/order", "OrderIndex", "syd", &order0.OrderIndex{})
    route.RegisterProton("syd/pages/order", "ShippingInsteadList", "syd", &order0.ShippingInsteadList{})
    route.RegisterProton("syd/pages/order", "ButtonSubmitHere", "syd", &order0.ButtonSubmitHere{})
    route.RegisterProton("syd/pages/order", "OrderList", "syd", &order0.OrderList{})
    route.RegisterProton("syd/pages/order", "OrderPrintNoPrice", "syd", &order0.OrderPrintNoPrice{})
    route.RegisterProton("syd/pages/order", "OrderQuery", "syd", &order0.OrderQuery{})
    route.RegisterProton("syd/pages/order", "PrintExpressYTO", "syd", &order0.PrintExpressYTO{})
    route.RegisterProton("syd/pages/order", "ShippingInsteadPrint", "syd", &order0.ShippingInsteadPrint{})
    route.RegisterProton("syd/pages/order", "DeliveringUnclosedOrders", "syd", &order0.DeliveringUnclosedOrders{})
    route.RegisterProton("syd/pages/order", "OrderPrint", "syd", &order0.OrderPrint{})
    route.RegisterProton("syd/pages/order", "ViewOrder", "syd", &order0.ViewOrder{})
    route.RegisterProton("syd/pages/order/create", "OrderCreateDetail", "syd", &order1.OrderCreateDetail{})
    route.RegisterProton("syd/pages/order/create", "OrderCreateIndex", "syd", &order1.OrderCreateIndex{})
    route.RegisterProton("syd/pages/order/create", "ShippingInstead", "syd", &order1.ShippingInstead{})
    route.RegisterProton("syd/pages/person", "PersonDetail", "syd", &person0.PersonDetail{})
    route.RegisterProton("syd/pages/person", "PersonEdit", "syd", &person0.PersonEdit{})
    route.RegisterProton("syd/pages/person", "EditAccountBallance", "syd", &person0.EditAccountBallance{})
    route.RegisterProton("syd/pages/person", "PersonIndex", "syd", &person0.PersonIndex{})
    route.RegisterProton("syd/pages/person", "PersonList", "syd", &person0.PersonList{})
    route.RegisterProton("syd/pages/placeorder", "PlaceOrderIndex", "syd", &inventory2.PlaceOrderIndex{})
    route.RegisterProton("syd/pages/product", "ProductList", "syd", &product0.ProductList{})
    route.RegisterProton("syd/pages/product", "ProductSendNewProduct", "syd", &product0.ProductSendNewProduct{})
    route.RegisterProton("syd/pages/product", "ProductIndex", "syd", &product0.ProductIndex{})
    route.RegisterProton("syd/pages/product", "ProductCreate", "syd", &product0.ProductCreate{})
    route.RegisterProton("syd/pages/product", "ProductDetail", "syd", &product0.ProductDetail{})
    route.RegisterProton("syd/pages/product", "ProductEdit", "syd", &product0.ProductEdit{})
    route.RegisterProton("syd/pages/service/suggest", "Suggest", "syd", &api1.Suggest{})
    route.RegisterProton("syd/pages/stat", "StatProductSold", "syd", &stat0.StatProductSold{})
    route.RegisterProton("syd/pages/stat", "StatToday", "syd", &stat0.StatToday{})
    route.RegisterProton("syd/pages/stat", "StatTrend", "syd", &stat0.StatTrend{})
    route.RegisterProton("syd/pages/test", "Test", "syd", &test2.Test{})
    route.RegisterProton("github.com/elivoa/got/builtin/components", "PageFinalBootstrap", "github.com/elivoa/got/builtin", &components0.PageFinalBootstrap{})
    route.RegisterProton("github.com/elivoa/got/builtin/components", "PageHeadBootstrap", "github.com/elivoa/got/builtin", &components0.PageHeadBootstrap{})
    route.RegisterProton("github.com/elivoa/got/builtin/components", "ProvinceSelect", "github.com/elivoa/got/builtin", &components0.ProvinceSelect{})
    route.RegisterProton("github.com/elivoa/got/builtin/components", "Select", "github.com/elivoa/got/builtin", &components0.Select{})
    route.RegisterProton("github.com/elivoa/got/builtin/components", "A", "github.com/elivoa/got/builtin", &components0.A{})
    route.RegisterProton("github.com/elivoa/got/builtin/components", "Delegate", "github.com/elivoa/got/builtin", &components0.Delegate{})
    route.RegisterProton("github.com/elivoa/got/builtin/components", "FileUpload", "github.com/elivoa/got/builtin", &components0.FileUpload{})
    route.RegisterProton("github.com/elivoa/got/builtin/components", "Output", "github.com/elivoa/got/builtin", &components0.Output{})
    route.RegisterProton("github.com/elivoa/got/builtin/components/got", "TemplateStatus", "github.com/elivoa/got/builtin", &got.TemplateStatus{})
    route.RegisterProton("github.com/elivoa/got/builtin/components/layout", "GOTHeader", "github.com/elivoa/got/builtin", &layout0.GOTHeader{})
    route.RegisterProton("github.com/elivoa/got/builtin/components/ui", "Pager", "github.com/elivoa/got/builtin", &ui.Pager{})
    route.RegisterProton("github.com/elivoa/got/builtin/pages", "PermissionDenied", "github.com/elivoa/got/builtin", &got0.PermissionDenied{})
    route.RegisterProton("github.com/elivoa/got/builtin/pages", "Error404", "github.com/elivoa/got/builtin", &got0.Error404{})
    route.RegisterProton("github.com/elivoa/got/builtin/pages", "Error500", "github.com/elivoa/got/builtin", &got0.Error500{})
    route.RegisterProton("github.com/elivoa/got/builtin/pages", "Errors", "github.com/elivoa/got/builtin", &got0.Errors{})
    route.RegisterProton("github.com/elivoa/got/builtin/pages/got", "Status", "github.com/elivoa/got/builtin", &got1.Status{})
    route.RegisterProton("github.com/elivoa/got/builtin/pages/got", "TestIndex", "github.com/elivoa/got/builtin", &got1.TestIndex{})
    route.RegisterProton("github.com/elivoa/got/builtin/pages/got/fileupload", "FileUploadTest", "github.com/elivoa/got/builtin", &fileupload.FileUploadTest{})
    route.RegisterProton("github.com/elivoa/got/builtin/pages/got/fileupload", "FileUploadIndex", "github.com/elivoa/got/builtin", &fileupload.FileUploadIndex{})

    // start the server
    _got.Start()
}
