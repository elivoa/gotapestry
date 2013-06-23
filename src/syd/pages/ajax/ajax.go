package ajax

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"got/core"
	"got/debug"
	"got/route"
	"syd/service/suggest"
)

type AjaxPage struct{}

func New() *AjaxPage {
	fmt.Println("Creating Ajax Page Module.")
	// templates.Add("person-list")
	return &AjaxPage{}
}

func (p *AjaxPage) Mapping(r *mux.Router) {
	r.HandleFunc("/ajax/suggest/{type}",
		route.PageHandler(&AjaxSuggest{}))
}

type AjaxSuggest struct {
	core.Page

	Query       string `query:"query"`
	SuggestType string `param:"type"`
}

func (p *AjaxSuggest) SetupRender() (interface{}, interface{}) {
	suggest.EnsureLoaded()
	// suggest.PrintAll()

	// search
	items, err := suggest.Lookup(p.Query, p.SuggestType)
	if err != nil {
		debug.Error(err)
		return err, nil
	}
	// translate
	sj := &SuggestJson{Query: p.Query}
	sj.Suggestions = make([]SuggestionJsonItem, len(*items))
	for idx, item := range *items {
		sj.Suggestions[idx] = SuggestionJsonItem{
			Value: item.QuickString + "||" + item.Text,
			Data:  fmt.Sprint(item.Id),
		}
	}

	// marshal
	jsonbytes, err := json.Marshal(sj)
	if err != nil {
		debug.Error(err)
		return err, nil
	}
	jsonstr := string(jsonbytes)
	fmt.Println(jsonstr)
	return "json", jsonbytes
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
