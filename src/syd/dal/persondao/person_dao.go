package persondao

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"got/db"
	"log"
	"syd/model"
	"time"
)

var logdebug = true
var em = &db.Entity{
	Table: "person",
	PK:    "id",
	Fields: []string{"id", "name", "type", "phone", "city", "address",
		"postalcode", "qq", "website", "note", "createtime", "updatetime"},
	CreateFields: []string{"name", "type", "phone", "city", "address",
		"postalcode", "qq", "website", "note", "createtime", "updatetime"}, // CreateFields
}

func init() {
	db.RegisterEntity("person", em)
}

// ________________________________________________________________________________
// Get person by person id
//
func Get(id int) (*model.Person, error) {
	p := new(model.Person)
	err := em.Select().Where("id", id).QueryOne(
		func(row *sql.Row) error {
			return row.Scan(
				&p.Id, &p.Name, &p.Type, &p.Phone, &p.City, &p.Address, &p.PostalCode, &p.QQ,
				&p.Website, &p.Note, &p.CreateTime, &p.UpdateTime,
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
				&p.Website, &p.Note, &p.CreateTime, &p.UpdateTime,
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
func Create(person *model.Person) (*model.Person, error) {
	res, err := em.Insert().Exec(
		person.Name, person.Type, person.Phone, person.City, person.Address, person.PostalCode,
		person.QQ, person.Website, person.Note, time.Now(), time.Now(),
	)
	if err != nil {
		return nil, err
	}
	liid, err := res.LastInsertId()
	person.Id = int(liid)
	return person, nil
}

/* update an existing item */
func UpdatePerson(person *model.Person) {
	db.Connect()
	defer db.Close()

	if logdebug {
		log.Printf("[dal] Edit person: %v", person)
	}

	stmt, err := db.DB.Prepare("update person set name=?, type=?, phone=?, city=?, address=?, postalcode=?, qq=?, website=?, note=?, updatetime=? where id=?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	stmt.Exec(person.Name, person.Type, person.Phone, person.City, person.Address, person.PostalCode, person.QQ, person.Website, person.Note, time.Now(), person.Id)
}

var EmptyPersonList = &[]model.Person{}

/*
  List person with type @the latest example of go-dao
*/
func ListPerson(personType string) (persons *[]model.Person, err error) {
	if logdebug {
		log.Printf("[dal] List person with type:%v", personType)
	}

	conn, _ := db.Connect()
	defer db.CloseConn(conn)

	stmt, err := conn.Prepare("select * from person where type=?")
	defer db.CloseStmt(stmt)
	if db.Err(err) {
		return
	}

	rows, err := stmt.Query(personType)
	defer db.CloseRows(rows)
	if db.Err(err) {
		return
	}

	// big performance issue, maybe.
	ps := []model.Person{}
	for rows.Next() {
		p := new(model.Person)
		rows.Scan(&p.Id, &p.Name, &p.Type, &p.Phone, &p.City, &p.Address, &p.PostalCode, &p.QQ, &p.Website, &p.Note, &p.CreateTime, &p.UpdateTime)
		ps = append(ps, *p)
	}
	persons = &ps
	return
}

func DeletePerson(id int) {
	if logdebug {
		log.Printf("[dal] delete person %n", id)
	}

	db.Connect()
	defer db.Close()

	stmt, err := db.DB.Prepare("delete from person where id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	stmt.Exec(id)
}
