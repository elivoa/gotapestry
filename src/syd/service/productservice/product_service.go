package productservice

import (
	"fmt"
	"path/filepath"
	"strings"
	"syd/dal"
	"syd/dal/productdao"
	"syd/model"
	"syd/service/suggest"
	"syd/utils"
)

//
// CreateProduct create a new Product into database, including it's properties.
//
func CreateProduct(product *model.Product) (*model.Product, error) {
	if product == nil {
		panic("Product can't be null!")
	}

	product.Capital = getCapital(product.Name)
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

	// update suggest
	suggest.Add(suggest.Product, newProduct.Name, newProduct.Id)

	return newProduct, nil
}

func UpdateProduct(product *model.Product) {
	if product == nil {
		return
	}
	// update product information
	product.Capital = getCapital(product.Name)
	if _, err := productdao.UpdateProduct(product); err != nil {
		panic(err.Error())
	}

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

	// update suggest
	suggest.Update(suggest.Product, product.Name, product.Id)

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

func ListProductsByCapital(capital string) ([]*model.Product, error) {
	return productdao.ListByCapital(capital)
}

func RebuildProductPinyinCapital() {
	fmt.Println("________________________________________________________________________________")
	fmt.Println("Rebuild Product Capital")
	products, err := ListProducts()
	if err != nil {
		panic(err.Error())
	}
	for _, product := range products {
		product.Capital = getCapital(product.Name)
		if _, err := productdao.UpdateProduct(product); err != nil {
			panic(err.Error())
		}
		fmt.Printf("> processing %v capital is: %v\n", product.Name, product.Capital)
	}
	fmt.Println("all done")
}

func getCapital(text string) string {
	s := utils.ParsePinyin(text)
	if len(s) > 0 {
		return s[0:1]
	}
	return "-"
}

func DeleteProduct(id int) (affacted int64, err error) {
	if affacted, err = productdao.Delete(id); err == nil {
		suggest.Delete(suggest.Product, id)
	}
	return
}
