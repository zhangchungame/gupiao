package main

import (
	"gupiao/screen"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"gupiao/allcode"
	"gupiao/rikxian"
)


func OrmInit()  {

	orm.RegisterDataBase("default", "mysql", "root:root@tcp(127.0.0.1:3306)/gupiao?charset=utf8")
	orm.RegisterModel(new(allcode.Baseinfo))
	orm.RegisterModel(new(rikxian.Rikxian))
	orm.RegisterModel(new(screen.Screen))

	orm.RegisterDriver("mysql", orm.DRMySQL)
}
func main() {
	OrmInit()
	screen.Calculator()
}
