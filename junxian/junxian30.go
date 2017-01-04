package junxian

import (
	"gupiao/allcode"
	"gupiao/rikxian"
	"github.com/astaxie/beego/orm"
	"gupiao/common"
	"fmt"
)



func CalculateAll_30()  {
	codes := allcode.Getallcodes()
	process_count:=0
	total_count:=0
	chprocess:=make(chan int)
	for _, val := range codes {
		//junxian30(val.Code)
		go junxian30(val.Code,chprocess)
		process_count++
		if process_count>100{
			<-chprocess
			process_count--
			total_count++
		}
	}
	for i := 0; i < process_count; i++ {
		<-chprocess
		total_count++
	}
	fmt.Println(total_count)
}

func junxian30(code string,chprocess chan int)  {
	fmt.Println(code+"start")
	db:=orm.NewOrm()
	var rikxians []rikxian.Rikxian
	_,err:=db.QueryTable("rikxian").Filter("code",code).Filter("zuigao__gt",0).OrderBy("date_int").Limit(-1).All(&rikxians)
	common.Checkerr(err)
	for index,_:=range rikxians{
		if index<29{
			continue
		}
		sum:=float64(0.0)
		for j:=0;j<30;j++{
			sum+=rikxians[index-29+j].Shoupan
		}
		rikxians[index].Rijun30=sum/30
		rikxians[index].Rijun30_cha=rikxians[index].Shoupan/rikxians[index].Rijun30-1
		db.QueryTable("rikxian").Filter("id",rikxians[index].Id).Update(orm.Params{
			"rijun30":rikxians[index].Rijun30,
			"rijun30_cha":rikxians[index].Rijun30_cha,
		})

	}
	fmt.Println(code+"finish")
	chprocess<-1
}

