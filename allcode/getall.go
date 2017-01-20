package allcode

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"fmt"
	"io/ioutil"
	"github.com/tealeg/xlsx"
	"time"
	"gupiao/common"
	"strings"
	"net/http"
)
type Baseinfo struct {
	Id int
	Code string
	Name string
	Jiaoyisuo string
	A_or_b string
	Market_time int
	Zong_gu_ben float64
	Liutong_gu_ben float64
}

func Getallcodes() []Baseinfo{
	db:=orm.NewOrm()
	var result []Baseinfo
	db.QueryTable("baseinfo").Filter("a_or_b","A").OrderBy("id").Limit(-1).All(&result)
	return result
}

func Setallcodes()  {
	loadCodeSz()
	loadCodeSh()
}

func loadCodeSz()  {
	excelFileName := "aaa.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	common.Checkerr(err)
	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i==0{
				continue
			}
			cells:=row.Cells
			code,err:=cells[5].String()
			common.Checkerr(err)
			name,err:=cells[6].String()
			common.Checkerr(err)

			if !strings.EqualFold(code,"")&&!strings.EqualFold(name,"") {
				tmp, err := cells[8].String()
				common.Checkerr(err)
				zong_gu_ben := moneytofloat(tmp)
				tmp, err = cells[9].String()
				common.Checkerr(err)
				liutong_gu_ben := moneytofloat(tmp)
				tmp, err = cells[7].String()
				common.Checkerr(err)
				mtime, err := time.ParseInLocation("2006-01-02", tmp, time.Local)
				if err!=nil {
					fmt.Println(err)
					mtime,_=time.ParseInLocation("2006-01-02", "1990-01-01", time.Local)
				}
				market_time := int(mtime.Unix())
				common.Checkerr(err)
				baseinfo := Baseinfo{0, code, name, "sz", "A", market_time, zong_gu_ben, liutong_gu_ben}
				updatedb(baseinfo)
			}

			//b股
			code,err=cells[10].String()
			common.Checkerr(err)
			name,err=cells[11].String()
			common.Checkerr(err)
			if !strings.EqualFold(code,"")&&!strings.EqualFold(name,"") {
				tmp,err:=cells[13].String()
				common.Checkerr(err)
				zong_gu_ben:=moneytofloat(tmp)
				tmp,err=cells[14].String()
				common.Checkerr(err)
				liutong_gu_ben:=moneytofloat(tmp)
				tmp,err=cells[12].String()
				common.Checkerr(err)
				mtime,err:=time.ParseInLocation("2006-01-02",tmp,time.Local)
				if err!=nil {
					fmt.Println(err)
					mtime,_=time.ParseInLocation("2006-01-02", "1990-01-01", time.Local)
				}
				market_time:=int(mtime.Unix())
				common.Checkerr(err)
				baseinfo:=Baseinfo{0,code,name,"sz","B",market_time,zong_gu_ben,liutong_gu_ben}
				updatedb(baseinfo)
			}
		}
	}
}
func moneytofloat(money string) float64 {
	arr:=strings.Split(money,",")
	tmpmoney:=""
	for _,d :=range arr{
		tmpmoney+=d
	}
	result,err:=strconv.ParseFloat(tmpmoney,64)
	common.Checkerr(err)
	return result

}
func loadCodeSh()  {
	client:=new(http.Client)
	url:="http://query.sse.com.cn/security/stock/downloadStockListFile.do?csrcCode=&stockCode=&areaName=&stockType=1"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header["Referer"]=append(req.Header["Referer"],"http://www.sse.com.cn/assortment/stock/list/share/")
	resp,err:= client.Do(req)
	if resp.Status!="200 OK"{
		return
	}
	robots, err := ioutil.ReadAll(resp.Body)
	b, err := common.Gb2utf_decode(robots)
	str:=string(b)
	codes:=strings.Split(str,"\n")
	baseinfo:=Baseinfo{}
	baseinfo.Jiaoyisuo="sh"
	baseinfo.A_or_b="b"
	var mtime time.Time
	for i,d := range codes{
		if(i==0){
			continue
		}
		arr:=strings.Split(d,"	  ")
		if(len(arr)==8){
			baseinfo.Code=arr[2]
			baseinfo.Name=arr[3]
			if len(arr[4])==10{
				mtime,_=time.ParseInLocation("2006-01-02",arr[4],time.Local)
			}else{
				mtime,_=time.ParseInLocation("2006-01-02", "1990-01-01", time.Local)
			}
			baseinfo.Market_time=int(mtime.Unix())
			baseinfo.Zong_gu_ben,err =strconv.ParseFloat(arr[5],64)
			common.Checkerr(err)
			baseinfo.Liutong_gu_ben,err=strconv.ParseFloat(arr[6],64)
			common.Checkerr(err)
			updatedb(baseinfo)
			fmt.Println(baseinfo)
		}
	}
	url="http://query.sse.com.cn/security/stock/downloadStockListFile.do?csrcCode=&stockCode=&areaName=&stockType=2"

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header["Referer"]=append(req.Header["Referer"],"http://www.sse.com.cn/assortment/stock/list/share/")
	resp,err= client.Do(req)
	robots, err = ioutil.ReadAll(resp.Body)
	b, err = common.Gb2utf_decode(robots)
	str=string(b)
	codes=strings.Split(str,"\n")
	baseinfo.Jiaoyisuo="sh"
	baseinfo.A_or_b="b"
	for i,d := range codes{
		if(i==0){
			continue
		}
		arr:=strings.Split(d,"	  ")
		if(len(arr)==8){
			baseinfo.Code=arr[2]
			baseinfo.Name=arr[3]
			if len(arr[4])==10{
				mtime,_=time.ParseInLocation("2006-01-02",arr[4],time.Local)
			}else{
				mtime,_=time.ParseInLocation("2006-01-02", "1990-01-01", time.Local)
			}
			baseinfo.Market_time=int(mtime.Unix())
			baseinfo.Zong_gu_ben,err =strconv.ParseFloat(arr[5],64)
			common.Checkerr(err)
			baseinfo.Liutong_gu_ben,err=strconv.ParseFloat(arr[6],64)
			common.Checkerr(err)
			updatedb(baseinfo)
		}
	}
}

func updatedb(baseinfo Baseinfo)  {
	db:=orm.NewOrm()
	exist := db.QueryTable("baseinfo").Filter("Code", baseinfo.Code).Exist()
	if exist{
		db.QueryTable("baseinfo").Filter("Code", baseinfo.Code).Update(orm.Params{
			"code": baseinfo.Code,
			"name":baseinfo.Name,
			"market_time":baseinfo.Market_time,
			"zong_gu_ben":baseinfo.Zong_gu_ben,
			"liutong_gu_ben":baseinfo.Liutong_gu_ben,
		})
		fmt.Println("update"+baseinfo.Code)
	}else{
		db.Insert(&baseinfo)
		fmt.Println("insert"+baseinfo.Code)
	}

	//stmt,err:=db.Prepare("select id from baseinfo where code=?")
	//common.Checkerr(err)
	//rows,err:=stmt.Query(baseinfo.code)
	//common.Checkerr(err)
	//isexist:=false
	//
	//for rows.Next(){
	//	isexist=true
	//}
	//if isexist{
	//	stmt,err:=db.Prepare("update baseinfo set code=?,name=?,jiaoyisuo=?,a_or_b=?,market_time=?,zong_gu_ben=?,liutong_gu_ben=? where id=?")
	//	common.Checkerr(err)
	//	_,err=stmt.Exec(baseinfo.code,baseinfo.name,baseinfo.jiaoyisuo,baseinfo.a_or_b,baseinfo.market_time,baseinfo.zong_gu_ben,baseinfo.liutong_gu_ben,baseinfo.id)
	//	common.Checkerr(err)
	//	fmt.Println("update")
	//}else{
	//	stmt,err:=db.Prepare("insert into baseinfo (code,name,jiaoyisuo,a_or_b,market_time,zong_gu_ben,liutong_gu_ben) values(?,?,?,?,?,?,?)")
	//	common.Checkerr(err)
	//	_,err=stmt.Exec(baseinfo.code,baseinfo.name,baseinfo.jiaoyisuo,baseinfo.a_or_b,baseinfo.market_time,baseinfo.zong_gu_ben,baseinfo.liutong_gu_ben)
	//	common.Checkerr(err)
	//	fmt.Println("insert")
	//}
}