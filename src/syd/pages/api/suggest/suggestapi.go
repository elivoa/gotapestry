package api

import (
	"encoding/json"
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"strings"
	"syd/service/suggest"
)

type Suggest struct {
	core.Page
	Type string `path-param:"1"`
	// PathQuery string `path-param:"2"`
	Query string `query:"query"`
}

func (p *Suggest) Setup() *exit.Exit {
	if strings.TrimSpace(p.Query) == "" {
		return exit.Json("{service:'/service'}")
	}
	suggest.EnsureLoaded()

	// search
	items, err := suggest.Lookup(p.Query, p.Type)
	if err != nil {
		return exit.Error(err)
	}
	// translate
	sj := &SuggestJson{Query: p.Query}
	sj.Suggestions = make([]SuggestionJsonItem, len(items))
	for idx, item := range items {
		sj.Suggestions[idx] = SuggestionJsonItem{
			Value: item.QuickString + "||" + item.Text,
			Data:  fmt.Sprint(item.Id),
		}
	}

	// marshal // use auto marshal.
	jsonbytes, err := json.Marshal(sj)
	if err != nil {
		return exit.Error(err)
	}
	jsonstr := string(jsonbytes)
	return exit.Json(jsonstr)
}

// Query product
func (p *Suggest) Onproduct() *exit.Exit {
	var query = strings.TrimSpace(p.Query)
	if query == "" {
		return exit.Json("{service:'no suggestions'}")
	}

	// search
	suggest.EnsureLoaded()
	items, err := suggest.Lookup(query, suggest.Product)
	if err != nil {
		return exit.Error(err)
	}
	// translate
	sj := &ProductSuggestions{Query: query}
	sj.Suggestions = make([]ProductSuggestionItem, len(items))
	for idx, item := range items {
		sj.Suggestions[idx] = ProductSuggestionItem{
			Id:          item.Id,
			ProductId:   item.SN,
			Name:        item.Text,
			QueryString: item.QuickString,
			// Data: fmt.Sprint(item.Id),
			// Value: item.QuickString + "||" + item.Text,
		}
	}

	// marshal // use auto marshal.
	jsonbytes, err := json.Marshal(sj)
	if err != nil {
		return exit.Error(err)
	}
	jsonstr := string(jsonbytes)
	return exit.Json(jsonstr)
}

// product suggest json
type ProductSuggestions struct {
	Query       string                  `json:"query"`
	Suggestions []ProductSuggestionItem `json:"suggestions"`
}

type ProductSuggestionItem struct {
	Id          int    `json:"id"`
	ProductId   string `json:"productId"`
	Name        string `json:"name"`
	QueryString string `json:"value"`
	// Value       string `json:"value"`
	// Data        string `json:"data"`
}

// 	// search
// 	items, err := suggest.Lookup(p.Query, p.Type)
// 	if err != nil {
// 		return err, nil
// 	}
// 	// translate
// 	sj := &SuggestJson{Query: p.Query}
// 	sj.Suggestions = make([]SuggestionJsonItem, len(items))
// 	for idx, item := range items {
// 		sj.Suggestions[idx] = SuggestionJsonItem{
// 			Value: item.QuickString + "||" + item.Text,
// 			Data:  fmt.Sprint(item.Id),
// 		}
// 	}

// 	// marshal
// 	jsonbytes, err := json.Marshal(sj)
// 	if err != nil {
// 		debug.Error(err)
// 		return err, nil
// 	}
// 	jsonstr := string(jsonbytes)
// 	return "json", jsonstr
// }

// // On query
// func (p *Suggest) OnQuery() (interface{}, interface{}) {
// 	suggest.EnsureLoaded()

// 	// search
// 	items, err := suggest.Lookup(p.Query, p.Type)
// 	if err != nil {
// 		return err, nil
// 	}
// 	// translate
// 	sj := &SuggestJson{Query: p.Query}
// 	sj.Suggestions = make([]SuggestionJsonItem, len(items))
// 	for idx, item := range items {
// 		sj.Suggestions[idx] = SuggestionJsonItem{
// 			Value: item.QuickString + "||" + item.Text,
// 			Data:  fmt.Sprint(item.Id),
// 		}
// 	}

// 	// marshal
// 	jsonbytes, err := json.Marshal(sj)
// 	if err != nil {
// 		debug.Error(err)
// 		return err, nil
// 	}
// 	jsonstr := string(jsonbytes)
// 	return "json", jsonstr
// }

/* struct to json */
type SuggestJson struct {
	Query       string               `json:"query"`
	Suggestions []SuggestionJsonItem `json:"suggestions"`
}

type SuggestionJsonItem struct {
	Value string `json:"value"`
	Data  string `json:"data"`
}
