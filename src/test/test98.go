package main

import (
	"encoding/json"
	"fmt"
	"syd/model"
)

func main() {
	jsonstr := "[{\"Id\":1,\"GroupId\":1,\"ProductId\":274,\"Color\":\"RED\",\"Size\":\"F\",\"Stock\":1500,\"ProviderId\":17,\"OperatorId\":3,\"Status\":0,\"Type\":0,\"Note\":\"this is a note\",\"SendTime\":\"2011-01-01T00:00:00Z\",\"ReceiveTime\":\"2011-05-01T00:00:00Z\",\"CreateTime\":\"2011-01-01T00:00:00Z\",\"UpdateTime\":\"2015-01-09T23:37:42Z\",\"Product\":null,\"Provider\":null,\"Operator\":null,\"Stocks\":{\"RED\":{\"F\":1500,\"M\":1},\"Green\":{\"S\":20}}},{\"Id\":2,\"GroupId\":1,\"ProductId\":23,\"Color\":\"BLACK\",\"Size\":\"F\",\"Stock\":90,\"ProviderId\":18,\"OperatorId\":3,\"Status\":3,\"Type\":0,\"Note\":\"note2\",\"SendTime\":\"0001-01-01T00:00:00Z\",\"ReceiveTime\":\"0001-01-01T00:00:00Z\",\"CreateTime\":\"0001-01-01T00:00:00Z\",\"UpdateTime\":\"2015-01-09T23:42:25Z\",\"Product\":null,\"Provider\":null,\"Operator\":null,\"Stocks\":{\"BLACK\":{\"F\":90}}}]"

	invs := []model.Inventory{}
	if err := json.Unmarshal([]byte(jsonstr), &invs); err == nil {
		if invs != nil {
			fmt.Println("invs is not nil")
			for idx, a := range invs {
				fmt.Println(idx, " : ", a)
				fmt.Println(" id: ", a.Id, "; productId: ", a.ProductId)
				fmt.Println("Stocks is", a.Stocks)
			}
		}
	} else {
		panic(err)
	}

}
