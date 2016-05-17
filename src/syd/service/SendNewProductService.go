package service

import (
	"database/sql"
	"fmt"
	"github.com/elivoa/got/db"
	"syd/dal/productdao"
	"syd/model"
)

type SendNewProductService struct{}

var KEY_SendNewProduct = "SendNewProduct"
var KEY_RecentProductItems = "RecentProductItems"

// 返回发样衣表中的所有数据。
// 默认显示1年内的数据；
func (s *SendNewProductService) GetSendNewProductTableData() (*model.ProductSalesTable, error) {
	ids := s.GetNewProductIds(s.GetSettingRecentProductItems())
	fmt.Println(">>>>>>>>>>> ids: ", ids)
	// TOBE CONTINUED....
	
	return nil, nil
}

func (s *SendNewProductService) GetSettingRecentProductItems() int {
	recentProductItems, err := Const.Get2ndIntValue(KEY_SendNewProduct, KEY_RecentProductItems)
	if err != nil {
		panic(err)
	}
	return int(recentProductItems)
}

func (s *SendNewProductService) GetNewProductIds(n int) []int64 {
	query := productdao.EntityManager().Select("id").OrderBy("CreateTime", db.DESC).Limit(10)

	var results = []int64{}
	var value int64
	err := query.Query(
		func(rows *sql.Rows) (bool, error) {
			if err2 := rows.Scan(&value); err2 != nil {
				panic(err2) // TODO a better way to solve this?
			}
			results = append(results, value)
			return true, nil
		},
	)
	if err != nil {
		panic(err)
	}
	return results
}
