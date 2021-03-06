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



func Jiaocuo30()  {
	codes := allcode.Getallcodes()
	process_count:=0
	total_count:=0
	chprocess:=make(chan int)
	for _, val := range codes {
		//junxian30(val.Code)
		go jiaocuo30(val.Code,chprocess)
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
func jiaocuo30(code string,ch chan int)  {
	total:=0
	hit:=0
	db:=orm.NewOrm()
	rikxians:=[]rikxian.Rikxian{}
	db.QueryTable("rikxian").Filter("code",code).OrderBy("date_int").Limit(-1).All(&rikxians)
	riklen:=len(rikxians)
	for index,rikxian:=range rikxians{
		if index<29{
			continue
		}
		if index==riklen-1{
			continue
		}
		if rikxian.Rijun30_cha>0{
			check:=false
			for i:=1;i<10;i++{
				if rikxians[index-i].Rijun30_cha>0{
					check=true
				}
			}
			if check==true{
				continue
			}
			total+=1
			if rikxians[index+1].Zhangdiefu>0 {
				hit+=1
			}
		}
	}
	fmt.Print(hit)
	fmt.Print("--")
	fmt.Print(total)
	fmt.Print("--")
	fmt.Println(float64(hit)/float64(total))
	ch<-1
}

