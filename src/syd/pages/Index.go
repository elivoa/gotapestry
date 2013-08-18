package index

import (
	"got/core"
	"syd/model"
	"syd/service/personservice"
)

// _______________________________________________________________________________
//  ROOT Page
//
type Index struct {
	core.Page
	Title     string
	Customers []*model.Person
}

func (p *Index) SetupRender() {
	p.Title = "圣衣蝶服饰销售管理系统"
	customers, err := personservice.ListCustomer()
	if err != nil {
		panic(err.Error())
	}
	personservice.SortByAccumulated(customers)

	// TODO: sort customers by accoumulated.
	// sort.Float64s

	// for _, c := range customers {
	// 	fmt.Println(c.AccountBallance)
	// }
	p.Customers = customers
}
