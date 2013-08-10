package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// TODO kill this
var DB *sql.DB

func Connect() (*sql.DB, error) {
	var err error
	conn, err := sql.Open("mysql", "root:eserver409$)(@/syd?charset=utf8&parseTime=true&loc=Local")
	if err != nil {
		panic(err.Error())
		return nil, err
	}
	DB = conn // kill this
	return conn, nil
}

func CloseConn(conn *sql.DB) {
	conn.Close()
	DB.Close()
}

// TODO delete this
func Close2(conn *sql.DB) {
	conn.Close()
	DB.Close()
}

func Close() {
	DB.Close()
}

func CloseStmt(stmt *sql.Stmt) {
	if stmt != nil {
		stmt.Close()
	}
}

func CloseRows(rows *sql.Rows) {
	if rows != nil {
		rows.Close()
	}
}

/*
   params
*/

type Filter struct {
	// TODO design a filter/parameter
}

/*
  Error handling
*/
func Err(err error) bool {
	if err != nil {
		panic(err.Error())
		return true
	}
	return false
}
