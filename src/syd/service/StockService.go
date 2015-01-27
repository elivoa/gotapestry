package service

import (
	"syd/dal/inventorydao"
	"sync"
)

// 库存数量处理的 Service; Table: product_pku

type StockService struct {
	stockLock sync.Mutex
}

// TODO batch this;
func (s *StockService) UpdateStockDelta(productId int64, color, size string, stockDelta int) error {
	// TODO lock
	s.stockLock.Lock()
	defer s.stockLock.Unlock()
	return inventorydao.UpdateProductStockWithDelta(productId, color, size, stockDelta)
}
