package main

import (
	"gupiao/allcode"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"gupiao/rikxian"
	"gupiao/junxian"
)

func OrmInit()  {

	orm.RegisterDataBase("default", "mysql", "root:root@tcp(127.0.0.1:3306)/gupiao?charset=utf8")
	orm.RegisterModel(new(allcode.Baseinfo))
	orm.RegisterModel(new(rikxian.Rikxian))

	orm.RegisterDriver("mysql", orm.DRMySQL)
}

func main() {
	OrmInit()
	rikxian.Getallkxian()
	junxian.CalculateAll_30()
	//o := orm.NewOrm()
	//o.Using("default") // 默认使用 default，你可以指定为其他数据库
	//
	//user := new(Baseinfo)
	//user.Code="123123"
	//user.Name="zc"
	//id,err:=o.Insert(user)
	//fmt.Println(id)
	//fmt.Println(err)
	//var infos []*Baseinfo
	//o.QueryTable("baseinfo").Filter("id__lt", 5).All(&infos)
	//fmt.Println(infos[2])
}
