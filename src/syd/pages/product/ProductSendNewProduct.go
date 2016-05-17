package product

import (
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"syd/model"
	"syd/service"
)

/*
   Product List page
   -------------------------------------------------------------------------------
*/
type ProductSendNewProduct struct {
	core.Page

	// parameters
	RecentProductItems int64

	// table data
	Data *model.ProductSalesTable

	// others.

	// Products []*model.Product
	// Capital string `path-param:"1"`
	// ShowAll bool   `query:"showall"`
	Referer string `query:"referer"` // return here if non-empty
}

func (p *ProductSendNewProduct) Setup() *exit.Exit {
	// // get items. TODO: Performance Issue.
	// recentProductItems, err := service.Const.Get2ndIntValue(SendNewProduct, RecentProductItems)
	// if err != nil {
	// 	panic(err)
	// }
	// p.RecentProductItems = recentProductItems

	// get data
	pst, err := service.SendNewProduct.GetSendNewProductTableData()
	if err != nil {
		panic(err)
	}
	p.Data = pst

	return nil
}
