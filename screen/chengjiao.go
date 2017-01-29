package screen

import (
	"github.com/astaxie/beego/orm"
	"gupiao/allcode"
	"time"
	"gupiao/common"
	"strings"
	"gupiao/rimingxi"
	"math"
	"gupiao/rikxian"
	"fmt"
	"strconv"
)

type Chengjiaoliang struct {
	Id int64
	Code string
	Date string
	Chaoda_buy int64
	Chaoda_sall int64
	Da_buy int64
	Da_sall int64
	Zhong_buy int64
	Zhong_sall int64
	Xiao_buy int64
	Xiao_sall int64
	Chaoda_buy_shou int64
	Chaoda_sall_shou int64
	Da_buy_shou int64
	Da_sall_shou int64
	Zhong_buy_shou int64
	Zhong_sall_shou int64
	Xiao_buy_shou int64
	Xiao_sall_shou int64
	Zhangdiefu float64
}

type fenbu struct {
	total int64
	hit int64
	persent float64
	chaoda_buy int64
	chaoda_sall int64
	Da_buy int64
	Da_sall int64
}

const(
	Chaoda float64   = 0.01
	Da float64       =0.005
	Zhong float64 =0.001
)


func ChengjiaoScreen()  {
	db:=orm.NewOrm()
	db.Raw("TRUNCATE `gupiao`.`chengjiaoliang`").Exec()
	codes:=allcode.Getallcodes()
	process_count:=0
	total_count:=0
	chprocess:=make(chan int)

	for _,val:=range codes{
		go ChengjiaoOneCode(val,chprocess)
		process_count++
		if process_count>20{
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

func ChengjiaoOneCode(baseinfo allcode.Baseinfo,ch chan int)  {
	startStr:="2016-01-01"
	startDate,err:=time.ParseInLocation("2006-01-02",startStr,time.Local)
	common.Checkerr(err)
	//startInt:=startDate.Unix()
	nowDate:=time.Now()
	for startDate.Unix()<=nowDate.Unix(){
		ChengjiaoOneCodeOneDay(startDate.Format("2006-01-02"),baseinfo)
		startDate=startDate.AddDate(0,0,1)
	}
	fmt.Println(baseinfo.Code+"++++++++++++++++++++++++++++++=finish")
	ch<-1
}

func ChengjiaoOneCodeOneDay(date string,baseinfo allcode.Baseinfo)  {
	db:=orm.NewOrm()
	rimingxis:=make([]rimingxi.Rimingxi,0)
	sql:="select * from rikxian where code='"+baseinfo.Code+"' and date='"+date+"'"
	rikxian:=rikxian.Rikxian{}
	db.Raw(sql).QueryRow(&rikxian)
	if rikxian.Chengjiaoliang<=0{
		return
	}
	sql="select * from rimingxi_"+baseinfo.Code+" where date='"+date+"' order by date_int asc"
	db.Raw(sql).QueryRows(&rimingxis)
	result:=Chengjiaoliang{}
	for index,rimingxi:=range rimingxis {
		if index==0{
			continue
		}
		abs:=math.Abs(float64(rimingxi.Chengjiaojia)/float64(rimingxis[index-1].Chengjiaojia)-1.0)
		if abs > Chaoda {
			if strings.EqualFold(rimingxi.Buy_sall, "b") {
				result.Chaoda_buy++
				result.Chaoda_buy_shou += rimingxi.Chengjiaoshou
			} else {
				result.Chaoda_sall++
				result.Chaoda_sall_shou += rimingxi.Chengjiaoshou
			}
		} else if abs >  Da {
			if strings.EqualFold(rimingxi.Buy_sall, "b") {
				result.Da_buy++
				result.Da_buy_shou += rimingxi.Chengjiaoshou
			} else {
				result.Da_sall++
				result.Da_sall_shou += rimingxi.Chengjiaoshou
			}
		} else if abs > Zhong {
			if strings.EqualFold(rimingxi.Buy_sall, "b") {
				result.Zhong_buy++
				result.Zhong_buy_shou += rimingxi.Chengjiaoshou
			} else {
				result.Zhong_sall++
				result.Zhong_sall_shou += rimingxi.Chengjiaoshou
			}
		} else {
			if strings.EqualFold(rimingxi.Buy_sall, "b") {
				result.Xiao_buy++
				result.Xiao_buy_shou += rimingxi.Chengjiaoshou
			} else {
				result.Xiao_sall++
				result.Xiao_sall_shou += rimingxi.Chengjiaoshou
			}
		}
	}
	result.Code=baseinfo.Code
	result.Date=date
	result.Zhangdiefu=rikxian.Zhangdiefu
	_,err:=db.Insert(&result)
	common.Checkerr(err)
	fmt.Println(date+"#####"+baseinfo.Code)
}

func Showfenbu()  {
	codes:=allcode.Getallcodes()
	process_count:=0
	total_count:=0
	chprocess:=make(chan int)

	for _,val:=range codes{
		go showfenbu(val,chprocess)
		process_count++
		if process_count>20{
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
func showfenbu(baseinfo allcode.Baseinfo,ch chan int)  {
	db:=orm.NewOrm()
	chengjiaoliangs:=make([]Chengjiaoliang,0)
	_,err:=db.QueryTable("chengjiaoliang").Filter("code",baseinfo.Code).OrderBy("-id").All(&chengjiaoliangs)
	common.Checkerr(err)
	chengjiaomap:=make(map[string]fenbu)
	clen:=len(chengjiaoliangs)
	for index,chengjiaoliang:=range chengjiaoliangs{
		if index==clen-1{
			continue
		}
		key:="da_buy"+strconv.FormatInt(chengjiaoliang.Da_buy,10)+"da_sall"+strconv.FormatInt(chengjiaoliang.Da_sall,10)
		if _,ok:=chengjiaomap[key];!ok{
			chengjiaomap[key]=fenbu{total:1,hit:0,chaoda_buy:chengjiaoliang.Chaoda_buy,chaoda_sall:chengjiaoliang.Chaoda_sall,Da_buy:chengjiaoliang.Da_buy,Da_sall:chengjiaoliang.Da_sall}
			if chengjiaoliangs[index+1].Zhangdiefu>0{
				tmp:=chengjiaomap[key]
				tmp.hit+=1
				chengjiaomap[key]=tmp
			}
		}else{
			tmp:=chengjiaomap[key]
			tmp.total+=1
			if chengjiaoliangs[index+1].Zhangdiefu>0{
				tmp.hit+=1
			}
			chengjiaomap[key]=tmp
		}
	}
	result:=make([]fenbu,0)
	for key,_:=range chengjiaomap{
		tmp:=chengjiaomap[key]
		if tmp.total>3{
			tmp.persent=float64(tmp.hit)/float64(tmp.total)
			result=append(result,tmp)
		}
	}
	reslen:=len(result)
	for i:=0;i<reslen;i++{
		for j:=reslen-1;j>i;j--{
			if result[j].persent>result[j-1].persent{
				tmp:=result[j]
				result[j]=result[j-1]
				result[j-1]=tmp
			}
		}
	}
	if reslen>3{
		if result[0].persent>0.9{
			result=result[:3]
			fmt.Print(baseinfo.Code+"----")
			fmt.Println(result)
		}
	}
	ch<-1
}