/**
  Time-stamp: <[suggest.go] Elivoa @ Saturday, 2016-11-12 23:10:24>
*/
package suggest

import (
	"errors"
	"fmt"
	"github.com/elivoa/got/debug"
	"strings"
	"syd/dal/persondao"
	"syd/dal/productdao"
	"syd/utils"
	"sync"
)

const (
	Customer = "customer"
	Factory  = "factory"
	Product  = "product"
)

var l sync.RWMutex
var cache map[string][]*Item //
var loaded bool

// var suggestCache
type Item struct {
	Id          int    // for Product
	SN          string // Product Id
	Text        string // Product Name
	QuickString string // capital of pinyin
	Type        string // 1-customer,2-factory,3-product
}

func init() {
	cache = make(map[string][]*Item, 10)
}

func EnsureLoaded() {
	if loaded {
		return
	}

	l.Lock()
	if !loaded {
		load()
		loaded = true
	}

	l.Unlock()
	PrintAll()
}

func IsLoaded() bool {
	return loaded
}

func load() {
	cache = make(map[string][]*Item, 100)

	persons, err := persondao.ListAll("customer")
	// persons, err := personservice.ListCustomer()
	if err != nil {
		panic(err.Error())
	} else {
		debug.Log("[suggest] load %v customers.", len(persons))
		personItems := make([]*Item, len(persons))
		cache[Customer] = personItems
		for i, person := range persons {
			personItems[i] = &Item{
				Id:          person.Id,
				Text:        person.Name,
				QuickString: parseQuickText(person.Name),
				Type:        "1",
			}
		}
	}

	factories, err := persondao.ListAll("factory")
	// factories, err := personservice.ListFactory()
	if err != nil {
		panic(err.Error())
	} else {
		debug.Log("[suggest] load %v factories.", len(factories))
		factoryItems := make([]*Item, len(factories))
		cache[Factory] = factoryItems
		for i, factory := range factories {
			factoryItems[i] = &Item{
				Id:          factory.Id,
				Text:        factory.Name,
				QuickString: parseQuickText(factory.Name), // TODO
				Type:        "2",
			}
		}
	}

	// TODO: chagne to step load.
	parser := productdao.EntityManager().NewQueryParser().Where().Limit(10000) // all limit to 1w
	products, err := productdao.List(parser)
	if err != nil {
		panic(err.Error())
	}
	// products := dal.ListProduct()
	debug.Log("[suggest] load %v products.", len(products))
	productItems := make([]*Item, len(products))
	cache[Product] = productItems
	for i, product := range products {
		productItems[i] = &Item{
			Id:          product.Id,
			SN:          product.ProductId,
			Text:        product.Name,
			QuickString: parseQuickText(product.Name), // TODO
			Type:        "3",
		}
	}

	debug.Log("Loading suggest done. Use 0.00 ms.")
}

// convert to pinyin
func parseQuickText(text string) string {
	return utils.ParsePinyin(text)
}

func Add(category string, text string, id int, sn string) {
	EnsureLoaded()

	item := &Item{
		Id:          id,
		Text:        text,
		SN:          sn,
		QuickString: parseQuickText(text),
		Type:        _categoryToType(category),
	}

	l.Lock()
	items, ok := cache[category]
	if !ok {
		cache[category] = []*Item{item}
	} else {
		items = append(items, item)
		cache[category] = items
	}
	l.Unlock()
}

func Delete(category string, id int) {
	EnsureLoaded()
	//
	l.Lock()
	items, ok := cache[category]
	if !ok {
		items = []*Item{}
		cache[category] = items
	}
	for i := 0; i < len(items); i++ {
		if items[i] != nil && items[i].Id == id {
			items[i] = nil
			break
		}
	}
	//
	l.Unlock()
}

func Update(category string, text string, id int, sn string) {
	Delete(category, id)
	Add(category, text, id, sn)
}

func PrintAll() {
	fmt.Println("------ Print All Suggest Items ---------------")
	//
	l.RLock()
	for key, value := range cache {
		fmt.Printf("> %v\n", key)
		for i, item := range value {
			if item == nil {
				fmt.Printf("\t%v: deleted\n", i)
			} else {
				fmt.Printf("\t%v: %v | %v\n", i, item.Text, item.QuickString)
			}
		}
	}

	l.RUnlock()
}

// 遍历法查找匹配项。
func Lookup(q string, category string) ([]*Item, error) {
	l.RLock()
	items, ok := cache[category]
	l.RUnlock()
	if !ok {
		err := errors.New(fmt.Sprintf("Category '%v' not found.", category))
		return nil, err
	}

	var (
		N        = 50
		idx      = 0
		filtered = make([]*Item, N, N)
		found    = 0
	)

	l.RLock()
	for _, item := range items {
		if item == nil {
			continue
		}
		// fmt.Println("LOOKUP:>", item.Text, item.Id, item.SN)

		var matched bool = false

		if strings.HasPrefix(item.QuickString, q) {
			matched = true
		}
		if strings.HasPrefix(item.SN, q) {
			matched = true
		}

		if matched {
			filtered[idx] = item
			found++
			if idx >= N-1 {
				break
			}
			idx++
		}
	}

	l.RUnlock()
	result := filtered[:found]
	return result, nil
}

// ---------- small functions --------------

func _categoryToType(category string) string {
	var Type string
	switch category {
	case Customer:
		Type = "1"
	case Factory:
		Type = "2"
	case Product:
		Type = "3"
	}
	return Type
}
