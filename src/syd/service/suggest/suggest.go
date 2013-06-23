package suggest

import (
	"errors"
	"fmt"
	"got/debug"
	"strings"
	"syd/dal"
	"syd/utils"
)

var (
	loaded = false
)

// var suggestCache
type Item struct {
	Id          int
	Text        string
	QuickString string
	Type        string
}

const (
	Customer = "customer"
	Factory  = "factory"
	Product  = "product"
)

// a simple version, match
var SuggestCache map[string][]Item // type->[]Item

func EnsureLoaded() {
	if !loaded {
		Load()
		PrintAll()
		loaded = true
	}
}

func IsLoaded() bool {
	return loaded
}

func Load() {
	SuggestCache = make(map[string][]Item, 100)

	persons, err := dal.ListPerson("customer")
	if err != nil {
		debug.Error(err)
	} else {
		debug.Log("[suggest] load %v customers.", len(*persons))
		personItems := make([]Item, len(*persons))
		SuggestCache[Customer] = personItems
		for i, person := range *persons {
			personItems[i] = Item{
				Id:          person.Id,
				Text:        person.Name,
				QuickString: parseQuickText(person.Name),
			}
		}
	}

	factories, err := dal.ListPerson("factory")
	if err != nil {
	} else {
		debug.Log("[suggest] load %v factories.", len(*factories))
		factoryItems := make([]Item, len(*factories))
		SuggestCache[Factory] = factoryItems
		for i, factory := range *factories {
			factoryItems[i] = Item{
				Id:          factory.Id,
				Text:        factory.Name,
				QuickString: parseQuickText(factory.Name), // TODO
			}
		}
	}

	products := dal.ListProduct()
	debug.Log("[suggest] load %v products.", len(*products))
	productItems := make([]Item, len(*products))
	SuggestCache[Product] = productItems
	for i, product := range *products {
		productItems[i] = Item{
			Id:          product.Id,
			Text:        product.Name,
			QuickString: parseQuickText(product.Name), // TODO
		}
	}

	debug.Log("Loading suggest done. Use 0.00 ms.")
}

// convert to pinyin
func parseQuickText(text string) string {
	s := utils.ParsePinyin(text)
	fmt.Println(s)
	return s
	//return text
}

func Add(category string, text string, id int) {
	item := Item{
		Id:          id,
		Text:        text,
		QuickString: parseQuickText(text),
	}

	items, ok := SuggestCache[category]
	if !ok {
		items = []Item{}
		SuggestCache[category] = items
	}
	items = append(items, item) // Performance?
}

func Delete() {} // TODO
func Update() {} // TODO

func PrintAll() {
	fmt.Println("------ Print All Suggest Items ---------------")
	for key, value := range SuggestCache {
		fmt.Printf("> %v\n", key)
		for i, item := range value {
			fmt.Printf("\t%v: %v | %v\n", i, item.Text, item.QuickString)
		}
	}
}

func LookupAll(q string) *[]Item {
	return nil
}

func Lookup(q string, category string) (*[]Item, error) {
	items, ok := SuggestCache[category]
	if !ok {
		err := errors.New(fmt.Sprintf("Category '%v' not found.", category))
		return nil, err
	}

	var (
		N        = 50
		idx      = 0
		filtered = make([]Item, N, N)
		found    = 0
	)
	for _, item := range items {
		if strings.HasPrefix(item.QuickString, q) {
			filtered[idx] = item
			found++
			if idx >= N-1 {
				break
			}
			idx++
		}
	}
	result := filtered[:found]
	return &result, nil
}









