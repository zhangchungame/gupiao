package allcode

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strings"
	"golang.org/x/text/transform"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gupiao/common"
	"time"
	"strconv"
	"github.com/tealeg/xlsx"
	"gupiao/singleInstance"
)
type Baseinfo struct {
	id int
	code string
	name string
	jiaoyisuo string
	a_or_b string
	market_time int
	zong_gu_ben float64
	liutong_gu_ben float64
}

var db *sql.DB

func init()  {
	db=singleInstance.GetMysqlInstance()
}

func Getallcodes() []interface{}{
	stmt,err:=db.Prepare("select code,jiaoyisuo from baseinfo where a_or_b=?")
	common.Checkerr(err)
	rows,err:=stmt.Query("A")
	common.Checkerr(err)
	result:=make([]interface{},0)
	//columns,_:=rows.Columns()
	code:=""
	jiaoyisuo:=""
	for rows.Next()  {
		err = rows.Scan(&code,&jiaoyisuo)
		common.Checkerr(err)
		tmp:=make(map[string]interface{})
		tmp["code"]=code
		tmp["jiaoyisuo"]=jiaoyisuo
		result=append(result,tmp)
	}
	return result
}
func decode(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	//O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	O := transform.NewReader(I, simplifiedchinese.GB18030.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}
func LoadCodeSz()  {
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

			if code!=""&&name!="" {
				tmp, err := cells[8].String()
				common.Checkerr(err)
				zong_gu_ben := moneytofloat(tmp)
				tmp, err = cells[9].String()
				common.Checkerr(err)
				liutong_gu_ben := moneytofloat(tmp)
				tmp, err = cells[7].String()
				common.Checkerr(err)
				mtime, err := time.ParseInLocation("2006-01-02", tmp, time.Local)
				common.Checkerr(err)
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
			if code!=""&&name!=""{
				tmp,err:=cells[13].String()
				common.Checkerr(err)
				zong_gu_ben:=moneytofloat(tmp)
				tmp,err=cells[14].String()
				common.Checkerr(err)
				liutong_gu_ben:=moneytofloat(tmp)
				tmp,err=cells[12].String()
				common.Checkerr(err)
				mtime,err:=time.ParseInLocation("2006-01-02",tmp,time.Local)
				common.Checkerr(err)
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
func DownloadCodeSh()  {
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
	b, err := decode(robots)
	str:=string(b)
	codes:=strings.Split(str,"\n")
	for i,d := range codes{
		if(i==0){
			continue
		}
		arr:=strings.Split(d,"	  ")
		if(len(arr)==8){
			mtime,err:=time.ParseInLocation("2006-01-02",arr[4],time.Local)
			common.Checkerr(err)
			market_time:=int(mtime.Unix())
			zong_gu_ben,err:=strconv.ParseFloat(arr[5],64)
			common.Checkerr(err)
			liutong_gu_ben,err:=strconv.ParseFloat(arr[6],64)
			common.Checkerr(err)
			baseinfo:=Baseinfo{0,arr[2],arr[3],"sh","A",market_time,zong_gu_ben,liutong_gu_ben}
			updatedb(baseinfo)
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
	b, err = decode(robots)
	str=string(b)
	codes=strings.Split(str,"\n")
	for i,d := range codes{
		if(i==0){
			continue
		}
		arr:=strings.Split(d,"	  ")
		if(len(arr)==8){
			mtime,err:=time.ParseInLocation("2006-01-02",arr[4],time.Local)
			common.Checkerr(err)
			market_time:=int(mtime.Unix())
			zong_gu_ben,err:=strconv.ParseFloat(arr[5],64)
			common.Checkerr(err)
			liutong_gu_ben,err:=strconv.ParseFloat(arr[6],64)
			common.Checkerr(err)
			baseinfo:=Baseinfo{0,arr[2],arr[3],"sh","B",market_time,zong_gu_ben,liutong_gu_ben}
			updatedb(baseinfo)
		}
	}
}

func updatedb(baseinfo Baseinfo)  {
	stmt,err:=db.Prepare("select id from baseinfo where code=?")
	common.Checkerr(err)
	rows,err:=stmt.Query(baseinfo.code)
	common.Checkerr(err)
	isexist:=false

	for rows.Next(){
		isexist=true
	}
	if isexist{
		stmt,err:=db.Prepare("update baseinfo set code=?,name=?,jiaoyisuo=?,a_or_b=?,market_time=?,zong_gu_ben=?,liutong_gu_ben=? where id=?")
		common.Checkerr(err)
		_,err=stmt.Exec(baseinfo.code,baseinfo.name,baseinfo.jiaoyisuo,baseinfo.a_or_b,baseinfo.market_time,baseinfo.zong_gu_ben,baseinfo.liutong_gu_ben,baseinfo.id)
		common.Checkerr(err)
		fmt.Println("update")
	}else{
		stmt,err:=db.Prepare("insert into baseinfo (code,name,jiaoyisuo,a_or_b,market_time,zong_gu_ben,liutong_gu_ben) values(?,?,?,?,?,?,?)")
		common.Checkerr(err)
		_,err=stmt.Exec(baseinfo.code,baseinfo.name,baseinfo.jiaoyisuo,baseinfo.a_or_b,baseinfo.market_time,baseinfo.zong_gu_ben,baseinfo.liutong_gu_ben)
		common.Checkerr(err)
		fmt.Println("insert")
	}
}