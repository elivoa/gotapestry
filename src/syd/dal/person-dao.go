package dal

// package main

import (
	_ "github.com/go-sql-driver/mysql"
	"got/db"
	"log"
	"syd/model"
	"time"
)

var logdebug = true

/*
   syd :: person
*/
/* create new item in db */
func CreatePerson(person *model.Person) {
	db.Connect()
	defer db.Close()

	if logdebug {
		log.Printf("[dal] Create person: %v", person)
	}

	stmt, err := db.DB.Prepare("insert into person(name, type, phone, city, address, postalcode, qq,website, note, updatetime) values(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	stmt.Exec(person.Name, person.Type, person.Phone, person.City, person.Address, person.PostalCode, person.QQ, person.Website, person.Note, time.Now())
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

func GetPerson(id int) *model.Person {
	if logdebug {
		log.Printf("[dal] Get Person with id %v", id)
	}

	db.Connect()
	defer db.Close()

	stmt, err := db.DB.Prepare("select * from person where id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	p := new(model.Person)
	row.Scan(&p.Id, &p.Name, &p.Type, &p.Phone, &p.City, &p.Address, &p.PostalCode, &p.QQ, &p.Website, &p.Note, &p.CreateTime, &p.UpdateTime)
	return p
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
