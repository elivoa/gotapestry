package service

import ()

type FactorySettleAccountService struct{}

// func (s *FactorySettleAccount) CreateAccountChangeLog(acl *model.AccountChangeLog) (
// 	*model.AccountChangeLog, error) {
// 	return accountdao.CreateAccountChangeLog(acl)
// }

// func (s *FactorySettleAccount) UpdateAccountBalance(personId int, delta float64,
// 	reason string, relatedOrderTrackingNo int64) {

// 	if person, err := Person.GetPersonById(personId); err != nil {
// 		panic(err)
// 	} else if person == nil {
// 		panic(fmt.Sprintf("Person %d not found!", personId))
// 	} else {
// 		person.AccountBallance += delta        // p.Order.SumOrderPrice()
// 		_, err := personservice.Update(person) // TODO
// 		if err != nil {
// 			panic(err.Error())
// 		}

// 		// create chagne log at the same time:
// 		accountdao.CreateAccountChangeLog(&model.AccountChangeLog{
// 			CustomerId:     person.Id,
// 			Delta:          delta,                  // -p.Order.SumOrderPrice(),
// 			Account:        person.AccountBallance, //
// 			Type:           2,                      // create takeaway order
// 			RelatedOrderTN: relatedOrderTrackingNo, // p.Order.TrackNumber,
// 			Reason:         reason,
// 		})

// 	}
// }
