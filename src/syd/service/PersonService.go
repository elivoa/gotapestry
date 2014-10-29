package service

import (
	"github.com/elivoa/got/db"
	"syd/dal/persondao"
	"syd/model"
)

type PersonService struct{}

func (s *PersonService) EntityManager() *db.Entity {
	return persondao.EntityManager()
}

// --------------------------------------------------------------------------------
// The following is helper function to fill user to models.
func (s *PersonService) _batchFetchPerson(ids []int64) (map[int64]*model.Person, error) {
	return persondao.ListPersonByIdSet(ids...)
}

func (s *PersonService) BatchFetchPerson(ids ...int64) (map[int64]*model.Person, error) {
	return s._batchFetchPerson(ids)
}

func (s *PersonService) BatchFetchPersonByIdMap(idset map[int64]bool) (map[int64]*model.Person, error) {
	var idarray = []int64{}
	if idset != nil {
		for id, _ := range idset {
			idarray = append(idarray, id)
		}
	}
	return s._batchFetchPerson(idarray)
}
