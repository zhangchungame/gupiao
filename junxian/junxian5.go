package junxian

import (
	"gupiao/allcode"
	"gupiao/rikxian"
	"github.com/astaxie/beego/orm"
	"gupiao/common"
	"fmt"
)



func CalculateAll_5()  {
	codes := allcode.Getallcodes()
	process_count:=0
	total_count:=0
	chprocess:=make(chan int)
	for _, val := range codes {
		//junxian30(val.Code)
		go junxian5(val.Code,chprocess)
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

func junxian5(code string,chprocess chan int)  {
	fmt.Println(code+"start")
	db:=orm.NewOrm()
	var rikxians []rikxian.Rikxian
	_,err:=db.QueryTable("rikxian").Filter("code",code).Filter("zuigao__gt",0).OrderBy("date_int").Limit(-1).All(&rikxians)
	common.Checkerr(err)
	for index,_:=range rikxians{
		if index<4{
			continue
		}
		sum:=float64(0.0)
		for j:=0;j<5;j++{
			sum+=rikxians[index-4+j].Shoupan
		}
		rikxians[index].Rijun5=sum/5
		rikxians[index].Rijun5_cha=rikxians[index].Shoupan/rikxians[index].Rijun5-1
		db.QueryTable("rikxian").Filter("id",rikxians[index].Id).Update(orm.Params{
			"rijun5":rikxians[index].Rijun5,
			"rijun5_cha":rikxians[index].Rijun5_cha,
		})
	}
	fmt.Println(code+"finish")
	chprocess<-1
}

func Jiaocuo5()  {
	codes := allcode.Getallcodes()
	process_count:=0
	total_count:=0
	chprocess:=make(chan int)
	for _, val := range codes {
		//junxian30(val.Code)
		go jiaocuo5(val.Code,chprocess)
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
func jiaocuo5(code string,ch chan int)  {
	total:=0
	hit:=0
	db:=orm.NewOrm()
	rikxians:=[]rikxian.Rikxian{}
	db.QueryTable("rikxian").Filter("code",code).OrderBy("date_int").Limit(-1).All(&rikxians)
	riklen:=len(rikxians)
	for index,rikxian:=range rikxians{
		if index<4{
			continue
		}
		if index==riklen-1{
			continue
		}
		if rikxian.Rijun5_cha>0{
			check:=false
			for i:=1;i<5;i++{
				if rikxians[index-i].Rijun5_cha>0{
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

