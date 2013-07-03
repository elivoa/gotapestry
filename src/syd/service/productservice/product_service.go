package productservice

import (
	"strings"
	"syd/dal"
	"syd/model"
)

//
// CreateProduct create a new Product into database, including it's properties.
//
func CreateProduct(product *model.Product) *model.Product {
	if product == nil {
		return nil
	}
	newProduct := dal.CreateProduct(product)
	if product.Colors != nil {
		dal.UpdateProductProperties(newProduct.Id, "color", product.Colors...)
	}
	if product.Sizes != nil {
		dal.UpdateProductProperties(newProduct.Id, "size", product.Sizes...)
	}
	return newProduct
}

func UpdateProduct(product *model.Product) {
	if product == nil {
		return
	}
	// update product information
	dal.UpdateProduct(product)

	// update it's properties
	if product.Colors != nil {
		dal.UpdateProductProperties(product.Id, "color", product.Colors...)
	}
	if product.Sizes != nil {
		dal.UpdateProductProperties(product.Id, "size", product.Sizes...)
	}

	// update stock information
	if product.Stocks != nil {
		dal.ClearProductStock(product.Id) // clear
		for key, stock := range product.Stocks {
			ps := strings.Split(key, "__")
			if len(ps) != 2 {
				panic("Key format not correct!" + key)
			}
			dal.SetProductStock(product.Id, ps[0], ps[1], stock)
		}
	}
}

//
// Get Product, with product's size and color properties.
// TODO get all properties.
//
func GetProduct(id int) *model.Product {
	product := dal.GetProduct(id)
	if product != nil {
		product.Colors = dal.GetProductProperties(id, "color")
		product.Sizes = dal.GetProductProperties(id, "size")
		product.Stocks = *dal.ListProductStocks(id)
	}
	return product
}
