/*
  Data Access Object for person module.

  Time-stamp: <[person_dao.go] Elivoa @ Wednesday, 2013-07-17 13:52:09>

  Note: This is the latest Template for dao functions.


*/
package persondao

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"got/db"
	"syd/model"
	"time"
)

var logdebug = true
var em = &db.Entity{
	Table: "person",
	PK:    "id",
	Fields: []string{"id", "name", "type", "phone", "city", "address",
		"postalcode", "qq", "website", "note", "AccountBallance", "createtime", "updatetime"},
	CreateFields: []string{"name", "type", "phone", "city", "address",
		"postalcode", "qq", "website", "note", "AccountBallance", "createtime", "updatetime"},
	UpdateFields: []string{"name", "type", "phone", "city", "address",
		"postalcode", "qq", "website", "note", "AccountBallance", "updatetime"},
}

func init() {
	db.RegisterEntity("person", em)
}

// ________________________________________________________________________________
// Get person by person id
//
func Get(id int) (*model.Person, error) {
	p := new(model.Person)
	err := em.Select().Where("id", id).Query(
		func(rows *sql.Rows) (bool, error) {
			return false, rows.Scan(
				&p.Id, &p.Name, &p.Type, &p.Phone, &p.City, &p.Address, &p.PostalCode, &p.QQ,
				&p.Website, &p.Note, &p.AccountBallance, &p.CreateTime, &p.UpdateTime,
			)
		},
	)
	if err != nil {
		return nil, err
	}
	if p.Id > 0 {
		return p, nil
	}
	return nil, errors.New("Person not found!")
}

// personType: customer, factory
func ListAll(personType string) ([]*model.Person, error) {
	persons := make([]*model.Person, 0)
	err := em.Select().Where("type", personType).Query(
		func(rows *sql.Rows) (bool, error) {
			p := new(model.Person)
			err := rows.Scan(
				&p.Id, &p.Name, &p.Type, &p.Phone, &p.City, &p.Address, &p.PostalCode, &p.QQ,
				&p.Website, &p.Note, &p.AccountBallance, &p.CreateTime, &p.UpdateTime,
			)
			persons = append(persons, p)
			return true, err
		},
	)
	if err != nil {
		return nil, err
	}
	return persons, nil
}

// ________________________________________________________________________________
// Create person
//
func Create(person *model.Person) error {
	res, err := em.Insert().Exec(
		person.Name, person.Type, person.Phone, person.City, person.Address, person.PostalCode,
		person.QQ, person.Website, person.Note, person.AccountBallance, time.Now(), time.Now(),
	)
	if err != nil {
		return err
	}
	liid, err := res.LastInsertId()
	person.Id = int(liid)
	return nil
}

// ________________________________________________________________________________
// Update returns RowsAffacted, error
//
func Update(person *model.Person) (int64, error) {
	res, err := em.Update().Exec(person.Name, person.Type, person.Phone, person.City, person.Address, person.PostalCode, person.QQ, person.Website, person.Note, person.AccountBallance, time.Now(), person.Id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

var EmptyPersonList = &[]model.Person{}

// ________________________________________________________________________________
// Delete a pesron
//
func Delete(id int) (int64, error) {
	res, err := em.Delete().Exec(id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
