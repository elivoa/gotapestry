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
	Type  string `path-param:"1"`
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

// On query
func (p *Suggest) OnQuery() *exit.Exit {
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

/* struct to json */
type SuggestJson struct {
	Query       string               `json:"query"`
	Suggestions []SuggestionJsonItem `json:"suggestions"`
}

type SuggestionJsonItem struct {
	Value string `json:"value"`
	Data  string `json:"data"`
}
