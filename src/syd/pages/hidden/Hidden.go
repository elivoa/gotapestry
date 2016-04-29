package test

import (
	"fmt"
	"github.com/elivoa/got/core"
	"github.com/elivoa/got/route/exit"
	"syd/dal/settleaccountdao"
	"syd/model"
	"time"
)

type Hidden struct {
	core.Page
	Data *model.ProductSalesTable
}

func (p *Hidden) SetupRender() *exit.Exit {
	fmt.Println("=============================================")
	starttime, err := time.Parse("2006-01-02", "2016-03-25")
	if err != nil {
		panic(err)
	}
	endtime, err := time.Parse("2006-01-02", "2016-03-28")
	if err != nil {
		panic(err)
	}
	fmt.Println("Time is from / to: ", starttime, endtime)

	// startTime, err := time.Parse("2006-01-02", starttime)
	// if nil != err {
	// 	panic(err)
	// }
	// endTime, err2 := time.Parse("2006-01-02", endtime)
	// if nil != err2 {
	// 	panic(err2)
	// }

	p.Data = model.NewTestProductSalesTable() //(startTime, endTime)

	pst, err := settleaccountdao.SettleAccount(starttime, endtime, 17)
	if err != nil {
		fmt.Println(err)
		return exit.Error(err)
	}
	p.Data = pst
	fmt.Println("m : ", p.Data)
	return nil
}

func (p *Hidden) GetData(date string, productId int64) string {
	value := p.Data.Get(date, productId)
	if value == 0 {
		return ""
	} else {
		return fmt.Sprint(value)
	}
}
