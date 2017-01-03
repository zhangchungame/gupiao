package rikxian
import (
	"net/http"
	"fmt"
	"strings"
	"time"
	"io/ioutil"
	"gupiao/allcode"
	"gupiao/common"
	"gupiao/singleInstance"
	"database/sql"
)
var count int

var db *sql.DB

func init()  {
	db=singleInstance.GetMysqlInstance()
}
func Getallkxian()  {
	codes:=allcode.Getallcodes()
	process_count:=0
	total_count:=0
	chprocess:=make(chan int)

	for _,val:=range codes{
		tmp:=val.(map[string]interface{})
		go Rikxianget(tmp,chprocess)
		//Rikxianget(tmp,chprocess)
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
func Rikxianget(tmp map[string]interface{},chprocess chan int) {
	//db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/gupiao?charset=utf8")
	var code string
	if tmp["jiaoyisuo"]=="sh"{
		code="0"+tmp["code"].(string)
	}else{
		code="1"+tmp["code"].(string)
	}
	fmt.Println(code,"---start")
	count=0
	nowstamp := time.Now().Unix()
	beginstr:="20160101000000"
	begintime,err:=time.ParseInLocation("20060102030405",beginstr,time.Local)
	common.Checkerr(err)
	beginstamp:=begintime.Unix()
	for startstamp:=beginstamp;startstamp<nowstamp ;startstamp+=3600*24  {
		date:=time.Unix(startstamp,0).Format("20060102")
		url:="http://quotes.money.163.com/service/chddata.html?code="+code+"&start="+date+"&end="+date+"&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;VOTURNOVER;VATURNOVER;TURNOVER;TCAP;MCAP"
		resp,err:=http.Get(url);

		common.Checkerr(err)
		robots, err := ioutil.ReadAll(resp.Body)
		common.Checkerr(err)
		str:=string(robots)
		resp.Body.Close()
		strs:=strings.Split(str, "\r\n")
		if strs[1] !=""{
			stmt,err:=db.Prepare("select id from rikxian where code=? and date_int=?")
			common.Checkerr(err)
			rows,err:=stmt.Query(tmp["code"].(string),startstamp)
			common.Checkerr(err)
			isexist:=false

			for rows.Next(){
				isexist=true
			}
			if(isexist){
				continue
			}
			res:=strings.Split(strs[1],",")
			stmt, err = db.Prepare(`INSERT rikxian (code,date,date_int,kaipan,shoupan,zuigao,zuidi,zhangdiee,zhangdiefu,chengjiaoliang,chengjiaoe,huanshoulv,zongshizhi,liutongshizhi) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
			common.Checkerr(err)
			_, err = stmt.Exec(tmp["code"].(string),res[0],startstamp,res[6],res[3],res[4],res[5],res[8],res[9],res[10],res[11],res[12],res[13],res[14])
			if err!=nil{
				fmt.Println(strs[1])
				fmt.Println(res)
				fmt.Println(err)
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
