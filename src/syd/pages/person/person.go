package person

import (
	"bytes"
	"fmt"
	"got/core"
	"gxl"
	"strconv"
	"syd/dal/orderdao"
	"syd/dal/persondao"
	"syd/model"
	"syd/service/orderservice"
	"syd/service/personservice"
	"time"
)

var (
	listTypeLabel = map[string]string{"customer": "客户", "factory": "厂商"}
)

/*_______________________________________________________________________________
  Person List Page
*/
type PersonIndex struct{ core.Page }

func (p *PersonIndex) Setup() (string, string) {
	return "redirect", "/person/list/customer"
}

/*_______________________________________________________________________________
  Person List Page
*/
type PersonList struct {
	core.Page

	PersonType string `path-param:"1" param:"type"`
	Persons    []*model.Person
	SubTitle   string
}

func (p *PersonList) Setup() interface{} {
	persons, err := personservice.List(p.PersonType)
	if err != nil {
		return err
	}
	p.Persons = persons
	p.SubTitle = listTypeLabel[p.PersonType]
	return nil
	// return "template", "person-list"
}

func (p *PersonList) Ondelete(personId int, personType string) (string, string) {
	personservice.Delete(personId)
	// TODO make this default redirect.
	return "redirect", fmt.Sprintf("/person/list/%v", personType)
}

/*_______________________________________________________________________________
  Person Edit Page
*/
type PersonEdit struct {
	core.Page

	Id     *gxl.Int `path-param:"1" required:"true" param:"id"`
	Person *model.Person

	Title    string
	SubTitle string

	TypeData interface{} // for type select
}

func (p *PersonEdit) Activate() {
	// here is some lightweight init.
	p.TypeData = &listTypeLabel
}

func (p *PersonEdit) Setup() {
	p.Title = "create/edit Person"

	if p.Id != nil {
		person, err := persondao.Get(p.Id.Int)
		if err != nil {
			// TODO how to handle error on page object?
			panic(err.Error())
		}
		p.Person = person
		p.SubTitle = "编辑"
	} else {
		p.Person = model.NewPerson()
		p.SubTitle = "新建"
	}
}

func (p *PersonEdit) OnSubmit() {
	if p.Id != nil {
		p.Person = personservice.GetPerson(p.Id.Int)
	} else {
		// No Need to edit.
	}
}

func (p *PersonEdit) OnSuccess() (string, string) {
	if p.Id != nil {
		personservice.Update(p.Person)
	} else {
		personservice.Create(p.Person)
	}
	return "redirect", fmt.Sprintf("/person/list/%v", p.Person.Type)
}

/* ________________________________________________________________________________
   PersonEdit
*/
type PersonDetail struct {
	core.Page

	Id *gxl.Int `path-param:"1"`

	Person         *model.Person
	Orders         []*model.Order
	LeavingMessage string
	// TodayOrders []*model.Order
}

func (p *PersonDetail) Setup() {
	if p.Id == nil {
		return
	}
	// performance issue: here we load all orders, this has an performance issue.
	p.Person = personservice.GetPerson(p.Id.Int)
	if p.Person != nil {
		orders, err := orderdao.ListOrderByCustomer(p.Person.Id, "all")
		if err != nil {
			panic(err.Error())
		}
		p.Orders = orders
	}

	// loop orders to find today's order and generate leaving message.
	// TODO support more than one order
	// TODO support shipping instead order.
	var msg bytes.Buffer
	for _, order := range p.Orders {
		// reload order with full details
		order, _ := orderservice.GetOrder(order.Id)
		y, m, d := order.CreateTime.Date()
		y2, m2, d2 := time.Now().Date()
		if y == y2 && m == m2 && d == d2 { // today's order
			jo := orderservice.OrderDetailsJson(order)
			var sumTotal float64
			var sumQuantity int
			for _, id := range jo.Orders {
				productJson := jo.Products[strconv.Itoa(id)]
				// 例如：奢华宝石
				msg.WriteString(productJson.Name)
				totalQuantity := 0
				for _, q := range productJson.Quantity {
					totalQuantity += q[2].(int)
					sumTotal += float64(q[2].(int)) * productJson.SellingPrice
				}
				sumQuantity += totalQuantity
				// eg: 1件
				msg.WriteString(strconv.Itoa(totalQuantity))
				msg.WriteString("件")

				// details
				if len(productJson.Quantity) >= 1 {
					msg.WriteString("(")
					i := 0
					for _, q := range productJson.Quantity {
						if i > 0 {
							msg.WriteString(", ")
						}
						i += 1
						_color := q[0].(string)
						_size := q[1].(string)
						if _color != "默认颜色" {
							msg.WriteString(_color)
						}
						if _size != "均码" {
							msg.WriteString(_size)
						}
						msg.WriteString(" ")
						msg.WriteString(strconv.Itoa(q[2].(int)))
					}
					msg.WriteString(")")
				}
				msg.WriteString("，")

				// price eg: xxx元
				msg.WriteString(fmt.Sprint(productJson.SellingPrice * float64(totalQuantity)))
				// msg.WriteString(gxl.FormatCurrency(productJson.SellingPrice*float64(totalQuantity), 2))
				msg.WriteString("元")
				msg.WriteString("；")
			}

			// 共计 n件 x元
			msg.WriteString("共计")
			msg.WriteString(strconv.Itoa(sumQuantity))
			msg.WriteString("件")
			msg.WriteString(gxl.FormatCurrency(sumTotal, 2))
			msg.WriteString("元")
			msg.WriteString("；")

			// shipping
			if order.DeliveryMethodIs("SF") {
				msg.WriteString("顺风运费")
			} else if order.DeliveryMethodIs("YTO") {
				msg.WriteString("圆通运费")
			} else {
				msg.WriteString("【快递类型错误】 运费")
			}
			if order.ExpressFee > 0 {
				msg.WriteString(fmt.Sprint(order.ExpressFee))
				msg.WriteString("元；")
			}
			msg.WriteString("单号")
			msg.WriteString(order.DeliveryTrackingNumber)
			msg.WriteString("; ")

			// 总计
			msg.WriteString("总计")
			msg.WriteString(gxl.FormatCurrency(sumTotal+float64(order.ExpressFee), 2))
			msg.WriteString("元")
			msg.WriteString("；")

			// 累计欠款
			if order.Accumulated > 0 {
				msg.WriteString("累计欠款：")
				msg.WriteString(fmt.Sprint(order.Accumulated))
				msg.WriteString(" + ")
				msg.WriteString(fmt.Sprint(int64(sumTotal) + order.ExpressFee))
				msg.WriteString(" = ")
				msg.WriteString(gxl.FormatCurrency(float64(int64(sumTotal)+order.ExpressFee)+order.Accumulated, 2))
				msg.WriteString("元")
				msg.WriteString("；")
			}
		}
	}
	p.LeavingMessage = msg.String()
}
