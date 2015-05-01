package service

import (
	"syd/dal/constdao"
)

var Const = new(ConstService)

type ConstService struct{}

func (s *ConstService) Set(name string, key string, value interface{}, floatValue float64) error {
	defer handleClear(name)
	return constdao.Set(name, key, value, floatValue)
}

func (s *ConstService) Update(name string, key string, value interface{}, floatValue float64, id int64) error {
	defer handleClear(name)
	return constdao.Update(name, key, value, floatValue, id)
}

func (s *ConstService) DeleteById(id int64) (int64, error) {
	c, err := constdao.GetById(id)
	if err != nil {
		panic(err)
	}
	handleClear(c.Name) // update cache
	return constdao.DeleteById(id)
}

func handleClear(name string) {
	// this is a fix for another system.

	// clear nothing

	// Carfilm.ClearConsts(name)
	// switch name {
	// case "car_made":
	// 	Carfilm.ClearCarMadeOptions()
	// case "window_tint":
	// 	Carfilm.ClearWindowTintOptions()
	// case "window_tint_optional":
	// 	Carfilm.ClearWindowTintOptOptions()
	// }
}
