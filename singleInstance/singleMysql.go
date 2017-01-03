package singleInstance

import (
	"sync"
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var mysqldb *sql.DB
var mysqldb_once sync.Once

func GetMysqlInstance() *sql.DB {
	mysqldb_once.Do(func() {
		var err error
		mysqldb, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/gupiao?charset=utf8")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	})
	return mysqldb
}
