/**
  Time-stamp: <[suggest.go] Elivoa @ 星期三, 2016-11-23 11:34:16>
*/
package suggest

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"syd/base/product"
	"syd/dal/persondao"
	"syd/dal/productdao"
	"syd/utils"
	"sync"

	"github.com/elivoa/got/debug"
)

const (
	Customer = "customer"
	Factory  = "factory"
	Product  = "product"
)

var l sync.RWMutex
var cache map[string]Items // inner suggest cache.
var loaded bool

// var suggestCache
type Item struct {
	Id          int            // for Product
	SN          string         // Product Id
	Text        string         // Product Name
	QuickString string         // capital of pinyin
	Type        string         // 1-customer,2-factory,3-product | category
	Status      product.Status // 0-normal, 1-hide. for product. TODO 暂时用了product的隐藏和显示.
}

type Items []*Item

func NewItems(len int) Items { return Items(make([]*Item, len)) }
func (p Items) Len() int     { return len(p) }
func (p Items) Less(i, j int) bool {
	if p[i] == nil || p[j] == nil {
		return false
	}
	return p[i].SN > p[j].SN // reverse order
}
func (p Items) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func init() {
	cache = make(map[string]Items, 10)
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
	cache = make(map[string]Items, 100)

	persons, err := persondao.ListAll("customer")
	// persons, err := personservice.ListCustomer()
	if err != nil {
		panic(err.Error())
	} else {
		debug.Log("[suggest] load %v customers.", len(persons))
		personItems := NewItems(len(persons))
		for i, person := range persons {
			personItems[i] = &Item{
				Id:          person.Id,
				Text:        person.Name,
				QuickString: parseQuickText(person.Name),
				Type:        "1",
			}
		}
		sort.Sort(personItems)
		cache[Customer] = personItems
	}

	factories, err := persondao.ListAll("factory")
	// factories, err := personservice.ListFactory()
	if err != nil {
		panic(err.Error())
	} else {
		debug.Log("[suggest] load %v factories.", len(factories))
		factoryItems := NewItems(len(factories))
		for i, factory := range factories {
			factoryItems[i] = &Item{
				Id:          factory.Id,
				Text:        factory.Name,
				QuickString: parseQuickText(factory.Name), // TODO
				Type:        "2",
			}
		}
		sort.Sort(factoryItems)
		cache[Factory] = factoryItems
	}

	// TODO: chagne to step load.
	parser := productdao.EntityManager().NewQueryParser().Where().Limit(10000) // all limit to 1w
	products, err := productdao.List(parser)
	if err != nil {
		panic(err.Error())
	}
	// products := dal.ListProduct()
	debug.Log("[suggest] load %v products.", len(products))
	productItems := NewItems(len(products))
	for i, p := range products {
		productItems[i] = &Item{
			Id:          p.Id,
			SN:          p.ProductId,
			Text:        p.Name,
			QuickString: parseQuickText(p.Name), // TODO
			Type:        "3",
			Status:      p.Status,
		}
	}
	sort.Sort(productItems)
	cache[Product] = productItems

	debug.Log("Loading suggest done. Use 0.00 ms.")
}

// convert to pinyin
func parseQuickText(text string) string {
	return utils.ParsePinyin2(text)
}

// Add add a suggest item into suggest-cache.
func Add(category string, text string, id int, sn string, status product.Status) {
	EnsureLoaded()

	item := &Item{
		Id:          id,
		Text:        text,
		SN:          sn,
		QuickString: parseQuickText(text),
		Type:        _categoryToType(category),
		Status:      status,
	}

	l.Lock()
	items, ok := cache[category]
	if !ok {
		cache[category] = Items([]*Item{item}) // newitems
	} else {
		items := append(items, item)
		sort.Sort(items)
		cache[category] = items
	}
	l.Unlock()
}

// Delete deletes item from suggest cache.
func Delete(category string, id int) {
	EnsureLoaded()
	l.Lock()
	items, ok := cache[category]
	if !ok {
		items = NewItems(0)
		cache[category] = items
	}
	for i := 0; i < len(items); i++ {
		if items[i] != nil && items[i].Id == id {
			items[i] = nil
			break
		}
	}
	l.Unlock()
}

// Update use item{Id,Text,SN,Type} to update an suggest item.
func Update(category string, item *Item /* text string, id int, sn string*/) {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>> Update category", item, item.Status)
	Delete(category, item.Id)
	Add(category, item.Text, item.Id, item.SN, item.Status)
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
		N        = 100
		idx      = 0
		filtered = make([]*Item, N, N)
		found    = 0
	)

	l.RLock()
	for _, item := range items {
		if item == nil {
			continue
		}
		fmt.Println("LOOKUP:>", item.Text, item.Id, item.SN, item.Status)
		fmt.Println("LOOKUP query:>", q, item.Text)

		var matched = false

		if strings.HasPrefix(item.QuickString, q) {
			matched = true
		}
		if strings.HasPrefix(item.SN, q) {
			matched = true
		}
		// 支持模糊搜索
		if strings.Contains(item.Text, q) {
			matched = true
		}
		if matched && item.Status == product.StatusNormal {
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
