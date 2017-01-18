package rimingxi

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"gupiao/common"
	"strings"
	"strconv"
	"time"
	"github.com/astaxie/beego/orm"
	"gupiao/allcode"
)

type Rimingxi struct {
	Id int64
	Code string
	Date string
	Date_int int64
	Chengjiaojia float64
	Zhangdiee float64
	Chengjiaoshou int64
	Chengjiaoe float64
	Buy_sall string
}

func RimingxigetAll()  {
	codes:=allcode.Getallcodes()
	process_count:=0
	total_count:=0
	chprocess:=make(chan int)

	for _,val:=range codes{
		go RimingxiTime(val,chprocess)
		//RimingxiTime(val,chprocess)
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

func RimingxiTime(baseinfo allcode.Baseinfo,ch chan int)  {
	startStr:="2015-01-01"
	startDate,err:=time.ParseInLocation("2006-01-02",startStr,time.Local)
	common.Checkerr(err)
	//startInt:=startDate.Unix()
	nowDate:=time.Now()
	for startDate.Unix()<=nowDate.Unix(){
		rimingxiget(startDate.Format("2006-01-02"),baseinfo)
		startDate=startDate.AddDate(0,0,1)
	}
	fmt.Println(baseinfo.Code+"++++++++++++++++++++++++++++++=finish")
	ch<-1
}

func rimingxiget(date string,baseinfo allcode.Baseinfo){
	db:=orm.NewOrm()
	exist := db.QueryTable("rimingxi").Filter("code", baseinfo.Code).Filter("date",date).Exist()
	if exist{
		return
	}
	url:="http://market.finance.sina.com.cn/downxls.php?date="+date+"&symbol="+baseinfo.Jiaoyisuo+baseinfo.Code;
	req,err:=http.Get(url);
	if(err !=nil){
		fmt.Println(err)
	}
	robots, err := ioutil.ReadAll(req.Body)

	b, err := common.Gb2utf_decode(robots)
	str:=string(b)
	req.Body.Close()
	strs:=strings.Split(str,"\n")
	result:=make([]Rimingxi,0)
	for index,val:=range strs{
		if index==0{
			continue
		}
		arr:=strings.Split(val,"\t")
		if len(arr)!=6{
			continue
		}
		tmp_rimingxi:=Rimingxi{}
		tmp_rimingxi.Code=baseinfo.Code
		tmp_rimingxi.Date=date
		//算時間
		tmptime,err:=time.ParseInLocation("2006-01-02 15:04:05",date+" "+arr[0],time.Local)
		common.Checkerr(err)
		tmp_rimingxi.Date_int=tmptime.Unix()
		tmp_rimingxi.Chengjiaojia,err=strconv.ParseFloat(arr[1],64)
		common.Checkerr(err)
		if strings.EqualFold(arr[2],"--"){
			tmp_rimingxi.Zhangdiee=0
		}else{
			tmp_rimingxi.Zhangdiee,err=strconv.ParseFloat(arr[2],64)
			common.Checkerr(err)
		}
		tmp_rimingxi.Chengjiaoshou,err=strconv.ParseInt(arr[3],10,64)
		common.Checkerr(err)
		tmp_rimingxi.Chengjiaoe,err=strconv.ParseFloat(arr[4],64)
		common.Checkerr(err)
		if strings.EqualFold(arr[5],"买盘"){
			tmp_rimingxi.Buy_sall="b"
		}else if strings.EqualFold(arr[5],"卖盘"){
			tmp_rimingxi.Buy_sall="s"
		}else{
			tmp_rimingxi.Buy_sall="e"
		}
		result=append(result,tmp_rimingxi)
	}
	db.InsertMulti(100,&result)
	fmt.Println(date+"#####"+baseinfo.Code)
}

