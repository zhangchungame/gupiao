package main

import (
	"gupiao/screen"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"gupiao/allcode"
	"gupiao/rikxian"
	"gupiao/junxian"
)


func OrmInit()  {

	orm.RegisterDataBase("default", "mysql", "root:root@tcp(127.0.0.1:3306)/gupiao?charset=utf8")
	orm.RegisterModel(new(allcode.Baseinfo))
	orm.RegisterModel(new(rikxian.Rikxian))
	orm.RegisterModel(new(screen.Screen))
	orm.RegisterModel(new(screen.Chengjiaoliang))

	orm.RegisterDriver("mysql", orm.DRMySQL)
}
func main() {
	OrmInit()
	//rimingxi.RimingxigetAll()
	//allcode.Setallcodes()
	//rikxian.Getallkxian()
	//junxian.CalculateAll_30()
	//junxian.CalculateAll_10()
	junxian.Jiaocuo30()

	//rimingxi.RimingxigetAll()
	//screen.ChengjiaoScreen()
	//screen.Showfenbu()
}


