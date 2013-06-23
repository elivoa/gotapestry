//package dal
package main

import (
	// "database/sql"
	"fmt"
	// "github.com/ziutek/mymysql/autorc"
	"github.com/ziutek/mymysql/mysql"
	//	_ "github.com/ziutek/mymysql/thrsafe"
	"log"
	"os"
	// "strconv"					
	// "syd/model"
"reflect"
	//	"time"
)

/*
   Connection Related
*/

const (
	db_proto = "tcp"
	db_addr  = "localhost:3306"
	db_user  = "root"
	db_pass  = "root"
	db_name  = "syd"
)

func init() {
	// Initialisation command
	// connect()
	// db.Raw.Register("SET NAMES utf8")
}

var db *mysql.Conn

func connect() {
	fmt.Println("-00---------")
	db := mysql.New("tcp", "", "127.0.0.1:3306", db_user, db_pass, db_name)
	fmt.Println(db)
	fmt.Println(reflect.TypeOf(db))
	fmt.Println("===9999999999")
	//fmt.Println(err)
	// db = autorc.New(db_proto, "", db_addr, db_user, db_pass, db_name)
}

func close() {
	//	defer db.Close()
}

func mysqlError(err error) (ret bool) {
	ret = (err != nil)
	if ret {
		log.Println("MySQL error:", err)
	}
	return
}

func mysqlErrExit(err error) {
	if mysqlError(err) {
		os.Exit(1)
	}
}

/*
   params
*/

type Filter struct {
	// TODO design a filter/parameter
}

// =============================================
func main() {
	fmt.Println("---------------")
	connect()
	//ListPerson("customer")
	fmt.Println("done")
}
