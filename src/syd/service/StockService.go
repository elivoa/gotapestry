package service

import (
	"syd/dal/inventorydao"
	"sync"
)

// 库存数量处理的 Service; Table: product_pku

type StockService struct {
	stockLock sync.Mutex
}

// TODO Batch this
// return old stock and new stock;
func (s *StockService) UpdateStockDelta(productId int64, color, size string, stockDelta int) (
	oldStock int, newStock int, err error) {
	// TODO lock
	s.stockLock.Lock()
	defer s.stockLock.Unlock()
	// get old stock
	if oldStock, err = inventorydao.GetProductStocks(productId, color, size); err != nil {
		return
	}

	// modify stock
	if err = inventorydao.UpdateProductStockWithDelta(productId, color, size, stockDelta); err != nil {
		return
	}
	// get new stock
	if newStock, err = inventorydao.GetProductStocks(productId, color, size); err != nil {
		return
	}
	return
}
