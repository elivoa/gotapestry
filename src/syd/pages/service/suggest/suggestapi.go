package api

import (
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
	return exit.MarshalJson(sj)
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
	return exit.MarshalJson(sj)
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
