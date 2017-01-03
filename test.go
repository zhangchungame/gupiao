package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"gupiao/common"
)

func init() {
	common.OrmInit()
}

func main() {
	o := orm.NewOrm()
	o.Using("default") // 默认使用 default，你可以指定为其他数据库

	user := new(Baseinfo)
	user.Code="123123"
	user.Name="zc"
	id,err:=o.Insert(user)
	fmt.Println(id)
	fmt.Println(err)
	var infos []*Baseinfo
	o.QueryTable("baseinfo").Filter("id__lt", 5).All(&infos)
	fmt.Println(infos[2])
}
