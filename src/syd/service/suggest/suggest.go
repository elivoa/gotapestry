/**
  Time-stamp: <[suggest.go] Elivoa @ Thursday, 2013-08-01 01:21:16>
*/
package suggest

import (
	"errors"
	"fmt"
	"got/debug"
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
var cache map[string][]*Item
var loaded bool

// var suggestCache
type Item struct {
	Id          int
	Text        string
	QuickString string
	Type        string
}

func init() {
	cache = make(map[string][]*Item, 10)
}

func EnsureLoaded() {
	if loaded {
		return
	}
	println("lock")
	l.Lock()
	if !loaded {
		load()
		loaded = true
	}
	println("unlock")
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
			}
		}
	}

	products, err := productdao.ListAll()
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
			Text:        product.Name,
			QuickString: parseQuickText(product.Name), // TODO
		}
	}

	debug.Log("Loading suggest done. Use 0.00 ms.")
}

// convert to pinyin
func parseQuickText(text string) string {
	return utils.ParsePinyin(text)
}

func Add(category string, text string, id int) {
	EnsureLoaded()
	item := &Item{
		Id:          id,
		Text:        text,
		QuickString: parseQuickText(text),
	}

	println("lock")
	l.Lock()
	items, ok := cache[category]
	if !ok {
		cache[category] = []*Item{item}
	} else {
		items = append(items, item)
		cache[category] = items
	}
	println("unlock")
	l.Unlock()
}

func Delete(category string, id int) {
	EnsureLoaded()
	println("lock")
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
	println("unlock")
	l.Unlock()
}

func Update(category string, text string, id int) {
	Delete(category, id)
	Add(category, text, id)
}

func PrintAll() {
	fmt.Println("------ Print All Suggest Items ---------------")
	println("r-lock")
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
	println("r-unlock")
	l.RUnlock()
}

func Lookup(q string, category string) ([]*Item, error) {

	println("lookup r-lock")
	l.RLock()
	items, ok := cache[category]
	println("lookup r-unlock")
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
	println("r-lock")
	l.RLock()
	for _, item := range items {
		if item == nil {
			continue
		}
		if strings.HasPrefix(item.QuickString, q) {
			filtered[idx] = item
			found++
			if idx >= N-1 {
				break
			}
			idx++
		}
	}
	println("r-unlock")
	l.RUnlock()
	result := filtered[:found]
	return result, nil
}
