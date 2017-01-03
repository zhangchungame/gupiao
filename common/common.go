package common

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"gupiao/allcode"
	_ "github.com/go-sql-driver/mysql"
)

func Checkerr(err error) {
	if err!=nil {
		fmt.Println(err)
	}
}

func OrmInit()  {
	orm.RegisterModel(new(allcode.Baseinfo))

	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:root@tcp(127.0.0.1:3306)/gupiao?charset=utf8")
}