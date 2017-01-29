package rikxian
import (
	"net/http"
	"fmt"
	"strings"
	"time"
	"io/ioutil"
	"gupiao/allcode"
	"gupiao/common"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type Rikxian struct {
	Id int64
	Code string
	Date string
	Date_int int64
	Kaipan float64
	Shoupan float64
	Zuigao float64
	Zuidi float64
	Zhangdiee float64
	Zhangdiefu float64
	Chengjiaoliang int64
	Chengjiaoe float64
	Huanshoulv float64
	Zongshizhi float64
	Liutongshizhi float64
	Rijun30 float64
	Rijun30_cha float64
	Rijun10 float64
	Rijun10_cha float64
	Rijun5 float64
	Rijun5_cha float64
}

func Getallkxian()  {
	codes:=allcode.Getallcodes()
	process_count:=0
	total_count:=0
	chprocess:=make(chan int)

	for _,val:=range codes{
		go Rikxianget(val,chprocess)
		//Rikxianget(val,chprocess)
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

func Rikxianget(baseinfo allcode.Baseinfo,chprocess chan int) {
	var code string
	if baseinfo.Jiaoyisuo=="sh"{
		code="0"+baseinfo.Code
	}else{
		code="1"+baseinfo.Code
	}
	fmt.Println(code,"---start")
	now_date:=time.Now().Format("20060102")
	begin_date:="20140101"
	url:="http://quotes.money.163.com/service/chddata.html?code="+code+"&start="+begin_date+"&end="+now_date+"&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;VOTURNOVER;VATURNOVER;TURNOVER;TCAP;MCAP"
	resp,err:=http.Get(url);

	common.Checkerr(err)
	robots, err := ioutil.ReadAll(resp.Body)
	common.Checkerr(err)
	str:=string(robots)
	resp.Body.Close()
	strs:=strings.Split(str, "\r\n")
	db:=orm.NewOrm()
	var tmptime time.Time
	for i,str:=range strs {
		if i==0{
			continue
		}
		res:=strings.Split(str,",")
		if len(res)==15{
			rikxian:=Rikxian{}
			rikxian.Code=baseinfo.Code
			rikxian.Date=res[0]
			if len(res[0])==10{
				tmptime,err=time.ParseInLocation("2006-01-02",res[0],time.Local)
			}else{
				tmptime,err=time.ParseInLocation("2006-01-02","1990-01-01",time.Local)
			}
			common.Checkerr(err)
			rikxian.Date_int=tmptime.Unix()

			rikxian.Kaipan,err=strconv.ParseFloat(res[6],64)
			common.Checkerr(err)
			rikxian.Shoupan,err=strconv.ParseFloat(res[3],64)
			common.Checkerr(err)
			rikxian.Zuigao,err=strconv.ParseFloat(res[4],64)
			common.Checkerr(err)
			rikxian.Zuidi,err=strconv.ParseFloat(res[5],64)
			common.Checkerr(err)
			rikxian.Zhangdiee,err=strconv.ParseFloat(res[8],64)
			if err!=nil{
				if err.Error()!="strconv.ParseFloat: parsing \"None\": invalid syntax"{
					fmt.Println(err)
				}
				rikxian.Zhangdiee=0
			}
			rikxian.Zhangdiefu,err=strconv.ParseFloat(res[9],64)
			if err!=nil{
				if err.Error()!="strconv.ParseFloat: parsing \"None\": invalid syntax"{
					fmt.Println(err)
				}
				rikxian.Zhangdiefu=0
			}
			rikxian.Chengjiaoliang,err=strconv.ParseInt(res[10],0,64)
			common.Checkerr(err)
			rikxian.Chengjiaoe,err=strconv.ParseFloat(res[11],64)
			common.Checkerr(err)
			rikxian.Huanshoulv,err=strconv.ParseFloat(res[12],64)
			common.Checkerr(err)
			rikxian.Zongshizhi,err=strconv.ParseFloat(res[13],64)
			common.Checkerr(err)
			rikxian.Liutongshizhi,err=strconv.ParseFloat(res[14],64)
			common.Checkerr(err)

			exist := db.QueryTable("rikxian").Filter("code", rikxian.Code).Filter("date_int",rikxian.Date_int).Exist()
			if !exist{
				db.Insert(&rikxian)
			}
		}
	}
	fmt.Println(code,"---end")
	chprocess<-1
}
//
//func rikxianget(code string,stamp int64)  {
//	date:=time.Unix(stamp,0).Format("20060102")
//	url:="http://quotes.money.163.com/service/chddata.html?code="+code+"&start="+date+"&end="+date+"&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;VOTURNOVER;VATURNOVER;TURNOVER;TCAP;MCAP"
//	resp,err:=http.Get(url);
//
//	common.Checkerr(err)
//	robots, err := ioutil.ReadAll(resp.Body)
//	common.Checkerr(err)
//	str:=string(robots)
//	resp.Body.Close()
//	strs:=strings.Split(str, "\r\n")
//	if strs[1] !=""{
//		stmt,err:=db.Prepare("select id from rikxian where code=? and date_int=?")
//		common.Checkerr(err)
//		rows,err:=stmt.Query()
//		common.Checkerr(err)
//		isexist:=false
//
//		for rows.Next(){
//			isexist=true
//		}
//		res:=strings.Split(strs[1],",")
//		stmt, err := db.Prepare(`INSERT rikxian (code,date,date_int,kaipan,shoupan,zuigao,zuidi,zhangdiee,zhangdiefu,chengjiaoliang,chengjiaoe,huanshoulv,zongshizhi,liutongshizhi) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
//		common.Checkerr(err)
//		_, err = stmt.Exec(code,res[0],stamp,res[6],res[3],res[4],res[5],res[8],res[9],res[10],res[11],res[12],res[13],res[14])
//		common.Checkerr(err)
//	}
//}
