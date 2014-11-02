// --- extends product service.
package service

import ()

// --------------------------------------------------------------------------------
// The following is helper function to fill user to models.
// func (s *ProductService) _batchFetchProductDetails(ids []int64) (map[int64]*model.Product, error) {
// 	return productdao.ListProductDetailsByIdSet(ids...)
// }

// func (s *ProductService) BatchFetchPerson(ids ...int64) (map[int64]*model.Product, error) {
// 	return s._batchFetchProduct(ids)
// }

// func (s *ProductService) BatchFetchProductDetailsByIdMap(idset map[int64]bool) (map[int64]*model.Product, error) {
// 	var idarray = []int64{}
// 	if idset != nil {
// 		for id, _ := range idset {
// 			idarray = append(idarray, id)
// 		}
// 	}
// 	return s._batchFetchProduct(idarray)
// }
