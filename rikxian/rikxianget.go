package rikxian

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gupiao/allcode"
	"gupiao/common"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"gupiao/singleInstance"
	"strconv"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Rikxian struct {
	Code string
	Date string
	Date_int int64
	Kaipan float64
	Shoupan float64
	Zuigao float64
	Zuidi float64
	Zhangdiee float64
	Zhangdiefu float64
	Chengjiaoliang float64
	Chengjiaoe float64
	Huanshoulv float64
	Zongshizhi float64
	Liutongshizhi float64
}
var collection *mgo.Collection
func Getallkxian() {
	codes := allcode.Getallcodes()
	collection=singleInstance.GetMongoInstance().C("rikxian")
	var code string
	process_count := 0
	total_count := 0
	chprocess := make(chan int)
	for _, val := range codes {
		if val.Jiaoyisuo == "sh" {
			code = "0" + val.Code
		} else {
			code = "1" + val.Code
		}
		go Rikxianget(code,  chprocess)
		//code="0600000"
		//Rikxianget(code,  chprocess)
		process_count++
		if process_count > 100 {
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
func Rikxianget(code string,chprocess chan int) {
	//db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/gupiao?charset=utf8")
	fmt.Println(code, "---start")
	nowstamp := time.Now().Unix()
	beginstr := "20140101000000"
	begintime, err := time.ParseInLocation("20060102030405", beginstr, time.Local)
	common.Checkerr(err)
	beginstamp := begintime.Unix()
	rikxianget(code, beginstamp,nowstamp)
	fmt.Println(code, "---end")
	chprocess <- 1
}

func rikxianget(code string,startstamp,endstamp int64)  {
	if(startstamp>endstamp){
		fmt.Println("开始时间大于结束时间")
		return;
	}
	startdate := time.Unix(startstamp, 0).Format("20060102")
	enddate := time.Unix(endstamp, 0).Format("20060102")
	url := "http://quotes.money.163.com/service/chddata.html?code=" + code + "&start=" + startdate + "&end=" + enddate + "&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;VOTURNOVER;VATURNOVER;TURNOVER;TCAP;MCAP"
	resp, err := http.Get(url)
	common.Checkerr(err)
	robots, err := ioutil.ReadAll(resp.Body)
	common.Checkerr(err)
	str := string(robots)
	resp.Body.Close()
	strs := strings.Split(str, "\r\n")
	for i,str:=range strs{
		if i==0 {
			continue;
		}
		if str!=""{
			res := strings.Split(str, ",")
			rik:=Rikxian{}
			//去掉第一个'号
			rs := []rune(res[1])
			rik.Code=string(rs[1:])

			datetime, err := time.ParseInLocation("2006-01-02", res[0], time.Local)
			var exist_baseinfo allcode.Baseinfo
			err=collection.Find(bson.M{"code":rik.Code,"date_int":datetime.Unix()}).One(&exist_baseinfo)
			if(exist_baseinfo.Id!=""){
				return
			}
			rik.Date=res[0]
			rik.Date_int=datetime.Unix()
			rik.Kaipan,err=strconv.ParseFloat(res[6],64)
			if err!=nil{
				return
			}
			rik.Shoupan,err=strconv.ParseFloat(res[3],64)
			if err!=nil{
				return
			}
			rik.Zuigao,err=strconv.ParseFloat(res[4],64)
			if err!=nil{
				return
			}
			rik.Zuidi,err=strconv.ParseFloat(res[5],64)
			if err!=nil{
				return
			}
			rik.Zhangdiee,err=strconv.ParseFloat(res[8],64)
			if err!=nil{
				return
			}
			rik.Zhangdiefu,err=strconv.ParseFloat(res[9],64)
			if err!=nil{
				return
			}
			rik.Chengjiaoliang,err=strconv.ParseFloat(res[10],64)
			if err!=nil{
				return
			}
			rik.Chengjiaoe,err=strconv.ParseFloat(res[11],64)
			if err!=nil{
				return
			}
			rik.Huanshoulv,err=strconv.ParseFloat(res[12],64)
			if err!=nil{
				return
			}
			rik.Zongshizhi,err=strconv.ParseFloat(res[13],64)
			if err!=nil{
				return
			}
			rik.Liutongshizhi,err=strconv.ParseFloat(res[14],64)
			if err!=nil{
				return
			}
			rik.Chengjiaoliang,err=strconv.ParseFloat(res[10],64)
			if err!=nil{
				return
			}
			collection.Insert(rik)
		}
	}
}

func rikxianget_oneday(code string, stamp int64) {
	date := time.Unix(stamp, 0).Format("20060102")
	url := "http://quotes.money.163.com/service/chddata.html?code=" + code + "&start=" + date + "&end=" + date + "&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;VOTURNOVER;VATURNOVER;TURNOVER;TCAP;MCAP"
	resp, err := http.Get(url)

	common.Checkerr(err)
	robots, err := ioutil.ReadAll(resp.Body)
	common.Checkerr(err)
	str := string(robots)
	resp.Body.Close()
	strs := strings.Split(str, "\r\n")
	if strs[1] != "" {
		var exist_baseinfo allcode.Baseinfo
		err=collection.Find(bson.M{"code":code,"date_int":stamp}).One(&exist_baseinfo)
		if(exist_baseinfo.Id!=""){
			return
		}
		res := strings.Split(strs[1], ",")
		rik:=Rikxian{}
		rik.Code=code
		rik.Date=res[0]
		rik.Date_int=stamp
		rik.Kaipan,err=strconv.ParseFloat(res[6],64)
		if err!=nil{
			return
		}
		rik.Shoupan,err=strconv.ParseFloat(res[3],64)
		if err!=nil{
			return
		}
		rik.Zuigao,err=strconv.ParseFloat(res[4],64)
		if err!=nil{
			return
		}
		rik.Zuidi,err=strconv.ParseFloat(res[5],64)
		if err!=nil{
			return
		}
		rik.Zhangdiee,err=strconv.ParseFloat(res[8],64)
		if err!=nil{
			return
		}
		rik.Zhangdiefu,err=strconv.ParseFloat(res[9],64)
		if err!=nil{
			return
		}
		rik.Chengjiaoliang,err=strconv.ParseFloat(res[10],64)
		if err!=nil{
			return
		}
		rik.Chengjiaoe,err=strconv.ParseFloat(res[11],64)
		if err!=nil{
			return
		}
		rik.Huanshoulv,err=strconv.ParseFloat(res[12],64)
		if err!=nil{
			return
		}
		rik.Zongshizhi,err=strconv.ParseFloat(res[13],64)
		if err!=nil{
			return
		}
		rik.Liutongshizhi,err=strconv.ParseFloat(res[14],64)
		if err!=nil{
			return
		}
		rik.Chengjiaoliang,err=strconv.ParseFloat(res[10],64)
		if err!=nil{
			return
		}
		collection.Insert(rik)
	}
}
