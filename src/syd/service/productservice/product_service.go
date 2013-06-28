package productservice

import (
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
	}
	return product
}
