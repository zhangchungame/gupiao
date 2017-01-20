package screen

import (
	"github.com/astaxie/beego/orm"
	"gupiao/rikxian"
	"fmt"
	"strconv"
	"gupiao/allcode"
	"strings"
)

type Screen struct {
	Id int64
	Zhangdiefu float64
	Rijun30_cha float64
	Hit_num int64
	Total_num int64
	Percent float64
}

func Calculator()  {
	db:=orm.NewOrm()
	db.Raw("TRUNCATE `gupiao`.`screen`").Exec()
	codes := allcode.Getallcodes()
	process_count:=0
	result:=make(map[string]Screen)
	ch:=make(chan map[string]Screen,10)
	for _, val := range codes {
		//junxian30(val.Code)
		go codeCalculator(val.Code,ch)
		process_count++
		if process_count>100{
			output:=<-ch
			for key,val:=range output{
				if _,ok:=result[key];!ok{
					result[key]=Screen{}
				}
				screen:=result[key]
				screen.Zhangdiefu=val.Zhangdiefu
				screen.Rijun30_cha=val.Rijun30_cha
				screen.Hit_num+=val.Hit_num
				screen.Total_num+=val.Total_num
				screen.Percent=float64(screen.Hit_num)/float64(screen.Total_num)
				result[key]=screen
			}
			process_count--
		}
	}
	for i := 0; i < process_count; i++ {
		output:=<-ch
		for key,val:=range output{
			if _,ok:=result[key];!ok{
				result[key]=Screen{}
			}
			screen:=result[key]
			screen.Zhangdiefu=val.Zhangdiefu
			screen.Rijun30_cha=val.Rijun30_cha
			screen.Hit_num+=val.Hit_num
			screen.Total_num+=val.Total_num
			screen.Percent=float64(screen.Hit_num)/float64(screen.Total_num)
			result[key]=screen
		}
	}
	screens := make([]Screen,0)
	count:=0;
	for key,val:=range result{
		//_,err:=db.Insert(&val)
		//if err!=nil{
		//	fmt.Println(err)
		//	fmt.Println(key)
		//	fmt.Println(val)
		//}
		if count>100{
			_,err:=db.InsertMulti(100,&screens)
			if err!=nil{
				fmt.Println(err)
				fmt.Println(key)
				fmt.Println(val)
			}
			count=0
			screens=screens[(len(screens)-1):]
		}
		screens=append(screens,val)
		count++
	}
}

func codeCalculator(code string,ch chan map[string]Screen){
	db:=orm.NewOrm()

	tmpmap :=make(map[string]map[string]int64)
	var rikxians []rikxian.Rikxian
	_,err:=db.QueryTable("rikxian").Filter("code",code).Filter("rijun30_cha__gt",0).Limit(-1).All(&rikxians)
	if err!=nil{
		fmt.Println("db_err="+err.Error())
		return
	}
	for _,rikxian :=range rikxians{
		rijun30:=strconv.FormatFloat(quzheng_bai(rikxian.Rijun30_cha),'f',3,64)
		zhangdifu:=strconv.FormatFloat(quzheng_shi(rikxian.Zhangdiefu),'f',3,64)
		if _,ok:=tmpmap[rijun30];!ok{
			tmpmap[rijun30]=make(map[string]int64)
		}
		if _,ok:=tmpmap[rijun30][zhangdifu];!ok{
			tmpmap[rijun30][zhangdifu]=1
		}else{
			tmpmap[rijun30][zhangdifu]+=1
		}
		if _,ok:=tmpmap[rijun30][zhangdifu];!ok{
			tmpmap[rijun30]["total"]=1
		}else{
			tmpmap[rijun30]["total"]+=1
		}
	}
	result:=make(map[string]Screen)
	for rijun_key,val:=range tmpmap{
		for zhangdi_key,val2:=range val{
			if strings.EqualFold(zhangdi_key,"total") {
				fmt.Println("相同")
				continue
			}
			screen:=Screen{0,0,0,val2,val["total"],float64(val2)/float64(val["total"])}
			screen.Rijun30_cha,_=strconv.ParseFloat(rijun_key,64)
			screen.Zhangdiefu,_=strconv.ParseFloat(zhangdi_key,64)
			result[rijun_key+"=="+zhangdi_key]=screen
		}
	}
	fmt.Println(code+"finish")
	ch<-result
}



func quzheng_qian(input float64) float64{
	input=input*1000
	input=float64(int64(input))/1000
	return input
}

func quzheng_bai(input float64) float64{
	input=input*100
	input=float64(int64(input))/100
	return input
}

func quzheng_shi(input float64) float64{
	input=input*10
	input=float64(int64(input))/10
	return input
}
