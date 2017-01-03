package main

import (
	"database/sql"
	"gupiao/common"
	"gupiao/rikxian"
)

var db *sql.DB
var dberr error
func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/gupiao?charset=utf8")
	common.Checkerr(err)
	rikxian.Getallkxian(db)
	//rikxian.Rikxianget(db)
	//allcode.DownloadCodeSh(db)
	//allcode.LoadCodeSz(db)
}