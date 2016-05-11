package test

import (
	"fmt"
	"github.com/elivoa/got/builtin/services"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"github.com/elivoa/got/utils"
	"syd/dal/settleaccountdao"
	"syd/model"
	"time"
)

type FactorySettleAccount struct {
	core.Page

	Data *model.ProductSalesTable

	// searc form or filter
	Provider int64     `query:"provider"` // filter by provider Id.
	TimeFrom time.Time `query:"from"`
	TimeTo   time.Time `query:"to"`

	// products map[int64]*model.Product

	Referer string // return to this place
}

func (p *FactorySettleAccount) SetupRender() *exit.Exit {

	// parameter time
	if !utils.IsValidTime(p.TimeFrom) {
		p.TimeFrom = time.Date(2015, time.January, 1, 0, 0, 0, 0, time.Local)
	}
	if !utils.IsValidTime(p.TimeTo) {
		p.TimeTo = time.Now()
	}

	fmt.Println("Time is from / to: ", p.TimeFrom, p.TimeTo)

	// p.Data = model.NewTestProductSalesTable() //(startTime, endTime)
	pst, err := settleaccountdao.SettleAccount(p.TimeFrom, p.TimeTo, p.Provider)
	if err != nil {
		fmt.Println(err)
		return exit.Error(err)
	}
	p.Data = pst

	return nil
}

// --------------------------------------------------------------------------------
// Search Form

func (p *FactorySettleAccount) OnSuccessFromSearchForm() *exit.Exit {
	// time is injected and then return linkpage.
	return exit.Redirect(p.ThisPageLink())
}

func (p *FactorySettleAccount) OnClearForm() *exit.Exit {
	p.TimeFrom = time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
	p.TimeTo = p.TimeFrom
	return exit.Redirect(p.ThisPageLink())
}

func (p *FactorySettleAccount) ThisPageLink() string {
	// 一个普通的SearchBox实现。所有东西都放到url里面。直接redirect到本页面。
	var parameters = map[string]interface{}{
		"provider": p.Provider,
		"from":     p.TimeFrom,
		"to":       p.TimeTo,
	}

	url := services.Link.GeneratePageUrlWithContextAndQueryParameters("accounting/factorysettleaccount",
		parameters)
	return url

}
