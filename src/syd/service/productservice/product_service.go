package productservice

import (
	"path/filepath"
	"strings"
	"syd/dal"
	"syd/dal/productdao"
	"syd/model"
)

//
// CreateProduct create a new Product into database, including it's properties.
//
func CreateProduct(product *model.Product) (*model.Product, error) {
	if product == nil {
		panic("Product can't be null!")
	}
	newProduct, err := productdao.Create(product)
	if err != nil {
		return nil, err
	}
	// newProduct := dal.CreateProduct(product)
	if product.Colors != nil {
		dal.UpdateProductProperties(newProduct.Id, "color", product.Colors...)
	}
	if product.Sizes != nil {
		dal.UpdateProductProperties(newProduct.Id, "size", product.Sizes...)
	}
	return newProduct, nil
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
	product, err := productdao.Get(id)
	if err == nil && product != nil {
		product.Colors = dal.GetProductProperties(id, "color")
		product.Sizes = dal.GetProductProperties(id, "size")
		product.Stocks = *dal.ListProductStocks(id)
	}
	return product
}

func ProductPictrues(product *model.Product) []string {
	if nil == product {
		return []string{}
	}
	pkeys := product.PictureKeys()
	for i := 0; i < len(pkeys); i++ {
		pkeys[i] = filepath.Join("/pictures", pkeys[i])
	}
	return pkeys
}

func ListProducts() ([]*model.Product, error) {
	return productdao.ListAll()
}
