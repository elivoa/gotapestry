package api

import (
	"encoding/json"
	"fmt"
	"got/core"
	"got/debug"
	"got/register"
	"syd/service/suggest"
)

func Register() {}

func init() { register.Page(Register, &Suggest{}) }

type Suggest struct {
	core.Page
	Type  string `path-param:"1"`
	Query string `query:"query"`
}

func (p *Suggest) Setup() (interface{}, interface{}) {
	suggest.EnsureLoaded()

	// search
	items, err := suggest.Lookup(p.Query, p.Type)
	if err != nil {
		return err, nil
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