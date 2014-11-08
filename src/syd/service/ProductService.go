package service

import (
	"fmt"
	"github.com/elivoa/got/db"
	"path/filepath"
	"strings"
	"syd/dal/inventorydao"
	"syd/dal/productdao"
	"syd/model"
	"syd/service/suggest"
	"syd/utils"
)

type ProductService struct{}

func (s *ProductService) EntityManager() *db.Entity {
	return productdao.EntityManager()
}

//
// CreateProduct create a new Product into database, including it's properties.
//
func (s *ProductService) CreateProduct(product *model.Product) (*model.Product, error) {
	if product == nil {
		panic("Product can't be null!")
	}

	product.Capital = s.getCapital(product.Name)
	newProduct, err := productdao.Create(product)
	if err != nil {
		return nil, err
	}
	// newProduct := dal.CreateProduct(product)
	if product.Colors != nil {
		productdao.UpdateProductProperties(newProduct.Id, "color", product.Colors...)
	}
	if product.Sizes != nil {
		productdao.UpdateProductProperties(newProduct.Id, "size", product.Sizes...)
	}

	// update suggest
	suggest.Add(suggest.Product, newProduct.Name, newProduct.Id)

	return newProduct, nil
}

func (s *ProductService) UpdateProduct(product *model.Product) {
	if product == nil {
		return
	}
	// update product information
	product.Capital = s.getCapital(product.Name)
	if _, err := productdao.UpdateProduct(product); err != nil {
		panic(err.Error())
	}

	// update it's properties
	if product.Colors != nil {
		productdao.UpdateProductProperties(product.Id, "color", product.Colors...)
	}
	if product.Sizes != nil {
		productdao.UpdateProductProperties(product.Id, "size", product.Sizes...)
	}

	// update stock information
	if product.Stocks != nil {
		inventorydao.ClearProductStock(product.Id) // clear
		for _, stock := range product.Stocks {
			inventorydao.SetProductStock(product.Id, stock.Color, stock.Size, stock.Stock)
		}
	}

	// update suggest
	suggest.Update(suggest.Product, product.Name, product.Id)

}

func (s *ProductService) DeleteProduct(id int) (affacted int64, err error) {
	if affacted, err = productdao.Delete(id); err != nil {
		return -1, err
	} else {
		suggest.Delete(suggest.Product, id)
		return
	}
}

// --------------------------------------------------------------------------------
// The following is helper function to fill user to models.
func (s *ProductService) _batchFetchProduct(ids []int64) (map[int64]*model.Product, error) {
	return productdao.ListProductsByIdSet(ids...)
}

func (s *ProductService) BatchFetchProduct(ids ...int64) (map[int64]*model.Product, error) {
	return s._batchFetchProduct(ids)
}

func (s *ProductService) BatchFetchProductByIdMap(idset map[int64]bool) (map[int64]*model.Product, error) {
	var idarray = []int64{}
	if idset != nil {
		for id, _ := range idset {
			idarray = append(idarray, id)
		}
	}
	return s._batchFetchProduct(idarray)
}

func (s *ProductService) getCapital(text string) string {
	str := utils.ParsePinyin(text)
	if len(str) > 0 {
		return strings.ToLower(str[0:1])
	}
	return "-"
}

// TODO: multi with version;
// func (s *ProductService) ListProducts() ([]*model.Product, error) {
// 	return productdao.ListAll()
// }

func (s *ProductService) List(parser *db.QueryParser, withs Withs) ([]*model.Product, error) {
	if models, err := productdao.List(parser); err != nil {
		return nil, err
	} else {
		// TODO: Print warrning information when has unused withs.
		// fmt.Println("--------------------------------------------------------------------", withs)
		if withs&WITH_PRODUCT_DETAIL > 0 {
			if err := productdao.FillProductPropertiesByIdSet(models); err != nil {
				return nil, err
			}
		}
		if withs&WITH_PRODUCT_INVENTORY > 0 {
			if err := inventorydao.FillProductStocksByIdSet(models); err != nil {
				return nil, err
			}
		}
		return models, nil
	}
}

//
// Get Product, with product's size and color properties.
//
func (s *ProductService) GetProduct(id int, withs Withs) (*model.Product, error) {
	if product, err := productdao.Get(id); err != nil {
		return nil, err
	} else if nil != product {
		models := []*model.Product{product}
		if withs&WITH_PRODUCT_DETAIL > 0 {
			if err := productdao.FillProductPropertiesByIdSet(models); err != nil {
				return nil, err
			}
		}
		if withs&WITH_PRODUCT_INVENTORY > 0 {
			if err := inventorydao.FillProductStocksByIdSet(models); err != nil {
				return nil, err
			}
		}
		return product, nil
	}
	return nil, nil
}

func (s *ProductService) GetFullProduct(id int) (*model.Product, error) {
	return s.GetProduct(id, WITH_PRODUCT_DETAIL|WITH_PRODUCT_INVENTORY)
}

// No use // TODO: delete this;
// func (s *ProductService) ListProductsByCapital(capital string) ([]*model.Product, error) {
// 	return productdao.ListByCapital(capital)
// }

// Non-standard fill.
// func (s *ProductService) FillProductsWithDetails(models []*model.Product) error {
// 	var idset = map[int64]bool{}
// 	for _, model := range models {
// 		idset[int64(model.Id)] = true
// 	}
// 	productdao.FillProductPropertiesByIdSet(models)

// 	personmap, err := Product.BatchFetchProductDetailsByIdMap(idset)
// 	if err != nil {
// 		return err
// 	}
// 	if nil != personmap && len(personmap) > 0 {
// 		for _, order := range orders {
// 			if person, ok := personmap[int64(order.CustomerId)]; ok {
// 				order.Customer = person
// 			}
// 		}
// 	}
// 	return nil
// }

// --------------------------------------------------------------------------------
func (s *ProductService) RebuildProductPinyinCapital() {
	fmt.Println("________________________________________________________________________________")
	fmt.Println("Rebuild Product Capital")

	qp := db.NewQueryParser().Limit(10000).Where()
	products, err := s.List(qp, 0)
	if err != nil {
		panic(err.Error())
	}
	for _, product := range products {
		product.Capital = s.getCapital(product.Name)
		if _, err := productdao.UpdateProduct(product); err != nil {
			panic(err.Error())
		}
		fmt.Printf("> processing %v capital is: %v\n", product.Name, product.Capital)
	}
	fmt.Println("all done")
}

func (s *ProductService) ProductPictrues(product *model.Product) []string {
	if nil == product {
		return []string{}
	}
	pkeys := product.PictureKeys()
	for i := 0; i < len(pkeys); i++ {
		pkeys[i] = filepath.Join("/pictures", pkeys[i])
	}
	return pkeys
}
