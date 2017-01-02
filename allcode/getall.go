package allcode

import (
	"bytes"
	"fmt"
	"github.com/tealeg/xlsx"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gupiao/common"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"gupiao/singleInstance"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Baseinfo_data struct {
	Code           string
	Name           string
	Jiaoyisuo      string
	A_or_b         string
	Market_time    int
	Zong_gu_ben    float64
	Liutong_gu_ben float64
}
type Baseinfo struct {
	Id 		bson.ObjectId    `bson:"_id"`
	Code           string
	Name           string
	Jiaoyisuo      string
	A_or_b         string
	Market_time    int
	Zong_gu_ben    float64
	Liutong_gu_ben float64
}

func Getallcodes() []Baseinfo {
	esclient:=singleInstance.GetEsInstance()
	fmt.Println(esclient)
	var result []Baseinfo
	mdb:=singleInstance.GetMongoInstance()
	collection:=mdb.C("baseinfo")
	collection.Find(bson.M{"a_or_b":"A"}).All(&result)
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

func InsertAllCode()  {
	LoadCodeSz()
	DownloadCodeSh()
}
func LoadCodeSz() {
	mdb:=singleInstance.GetMongoInstance()
	collection:=mdb.C("baseinfo")
	excelFileName := "aaa.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	common.Checkerr(err)
	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i == 0 {
				continue
			}
			cells := row.Cells
			code, err := cells[5].String()
			common.Checkerr(err)
			name, err := cells[6].String()
			common.Checkerr(err)

			if code != "" && name != "" {
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
				baseinfo := Baseinfo_data{ code, name, "sz", "A", market_time, zong_gu_ben, liutong_gu_ben}
				updatedb(baseinfo, collection)
			}

			//b股
			code, err = cells[10].String()
			common.Checkerr(err)
			name, err = cells[11].String()
			common.Checkerr(err)
			if code != "" && name != "" {
				tmp, err := cells[13].String()
				common.Checkerr(err)
				zong_gu_ben := moneytofloat(tmp)
				tmp, err = cells[14].String()
				common.Checkerr(err)
				liutong_gu_ben := moneytofloat(tmp)
				tmp, err = cells[12].String()
				common.Checkerr(err)
				mtime, err := time.ParseInLocation("2006-01-02", tmp, time.Local)
				common.Checkerr(err)
				market_time := int(mtime.Unix())
				common.Checkerr(err)
				baseinfo := Baseinfo_data{code, name, "sz", "B", market_time, zong_gu_ben, liutong_gu_ben}
				updatedb(baseinfo, collection)
			}
		}
	}
}
func moneytofloat(money string) float64 {
	arr := strings.Split(money, ",")
	tmpmoney := ""
	for _, d := range arr {
		tmpmoney += d
	}
	result, err := strconv.ParseFloat(tmpmoney, 64)
	common.Checkerr(err)
	return result

}
func DownloadCodeSh() {
	mdb:=singleInstance.GetMongoInstance()
	collection:=mdb.C("baseinfo")
	client := new(http.Client)
	url := "http://query.sse.com.cn/security/stock/downloadStockListFile.do?csrcCode=&stockCode=&areaName=&stockType=1"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header["Referer"] = append(req.Header["Referer"], "http://www.sse.com.cn/assortment/stock/list/share/")
	resp, err := client.Do(req)
	if resp.Status != "200 OK" {
		return
	}
	robots, err := ioutil.ReadAll(resp.Body)
	b, err := decode(robots)
	str := string(b)
	codes := strings.Split(str, "\n")
	for i, d := range codes {
		if i == 0 {
			continue
		}
		arr := strings.Split(d, "	  ")
		if len(arr) == 8 {
			mtime, err := time.ParseInLocation("2006-01-02", arr[4], time.Local)
			common.Checkerr(err)
			market_time := int(mtime.Unix())
			zong_gu_ben, err := strconv.ParseFloat(arr[5], 64)
			common.Checkerr(err)
			liutong_gu_ben, err := strconv.ParseFloat(arr[6], 64)
			common.Checkerr(err)
			baseinfo := Baseinfo_data{ arr[2], arr[3], "sh", "A", market_time, zong_gu_ben, liutong_gu_ben}
			updatedb(baseinfo, collection)
		}
	}
	url = "http://query.sse.com.cn/security/stock/downloadStockListFile.do?csrcCode=&stockCode=&areaName=&stockType=2"

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header["Referer"] = append(req.Header["Referer"], "http://www.sse.com.cn/assortment/stock/list/share/")
	resp, err = client.Do(req)
	robots, err = ioutil.ReadAll(resp.Body)
	b, err = decode(robots)
	str = string(b)
	codes = strings.Split(str, "\n")
	for i, d := range codes {
		if i == 0 {
			continue
		}
		arr := strings.Split(d, "	  ")
		if len(arr) == 8 {
			mtime, err := time.ParseInLocation("2006-01-02", arr[4], time.Local)
			common.Checkerr(err)
			market_time := int(mtime.Unix())
			zong_gu_ben, err := strconv.ParseFloat(arr[5], 64)
			common.Checkerr(err)
			liutong_gu_ben, err := strconv.ParseFloat(arr[6], 64)
			common.Checkerr(err)
			baseinfo := Baseinfo_data{ arr[2], arr[3], "sh", "B", market_time, zong_gu_ben, liutong_gu_ben}
			updatedb(baseinfo, collection)
		}
	}
}

func updatedb(baseinfo Baseinfo_data, collection *mgo.Collection) {
	var exist_baseinfo Baseinfo
	err:=collection.Find(bson.M{"code":baseinfo.Code}).One(&exist_baseinfo)
	if(exist_baseinfo.Id==""){
		err=collection.Insert(baseinfo)
	}else{
		err=collection.Update(bson.M{"_id":exist_baseinfo.Id},baseinfo)
	}
	common.Checkerr(err)
}
