package rimingxi

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"gupiao/common"
	"strings"
	"time"
	"github.com/astaxie/beego/orm"
	"gupiao/allcode"
	"strconv"
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

func RimingxiTime(baseinfo allcode.Baseinfo,ch chan int)  {
	sql:="CREATE TABLE IF NOT EXISTS `rimingxi_"+baseinfo.Code+"` (`id` int(11) NOT NULL AUTO_INCREMENT,`date` varchar(14) NOT NULL DEFAULT '',`date_int` int(11) NOT NULL DEFAULT '0',`chengjiaojia` double(8,4) NOT NULL DEFAULT '0.0000',`zhangdiee` double(8,4) NOT NULL DEFAULT '0.0000',`chengjiaoshou` int(11) NOT NULL DEFAULT '0',`chengjiaoe` double(20,4) NOT NULL DEFAULT '0.0000',`buy_sall` varchar(2) NOT NULL DEFAULT '',PRIMARY KEY (`id`),KEY `date` (`date`),KEY `date_int` (`date_int`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;";
	db:=orm.NewOrm()
	_,err:=db.Raw(sql).Exec()
	common.Checkerr(err)
	sql="select * from rimingxi_"+baseinfo.Code+" order by id desc limit 0,1"
	tmprimingxi:=Rimingxi{}
	db.Raw(sql).QueryRow(&tmprimingxi)
	startStr:="2016-01-01"
	if tmprimingxi.Date_int>0{
		startStr=tmprimingxi.Date
	}
	fmt.Println("startDate======================="+startStr)
	startDate,err:=time.ParseInLocation("2006-01-02",startStr,time.Local)
	common.Checkerr(err)
	//startInt:=startDate.Unix()
	nowDate:=time.Now()
	for startDate.Unix()<=nowDate.Unix(){
		if rimingxiget(startDate.Format("2006-01-02"),baseinfo)==0{
			rimingxiget(startDate.Format("2006-01-02"),baseinfo)
		}
		startDate=startDate.AddDate(0,0,1)
	}
	fmt.Println(baseinfo.Code+"++++++++++++++++++++++++++++++=finish")
	ch<-1
}

func rimingxiget(date string,baseinfo allcode.Baseinfo) int64{
	defer func() {
		if err := recover(); err != nil {
			fmt.Print("error___________________________")
			fmt.Println(err)
		}
	}()
	db:=orm.NewOrm()
	var isexist Rimingxi
	sql:="select * from rimingxi_"+baseinfo.Code+" where date='"+date+"'"
	db.Raw(sql).QueryRow(&isexist)
	if isexist.Id!=0{
		fmt.Println(baseinfo.Code+"----"+date+"-----isexist")
		return 2
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

	sql="insert into rimingxi_"+baseinfo.Code+" (date,date_int,chengjiaojia,zhangdiee,chengjiaoshou,chengjiaoe,buy_sall) values";
	sql_arr:=make([]string,0);
	count:=0;
	for index,val:=range strs{
		if index==0{
			continue
		}
		arr:=strings.Split(val,"\t")
		if len(arr)!=6{
			continue
		}
		tmptime,err:=time.ParseInLocation("2006-01-02 15:04:05",date+" "+arr[0],time.Local)
		common.Checkerr(err)
		zhangdiee:=arr[2];
		if strings.EqualFold(arr[2],"--"){
			zhangdiee="0"
		}
		buy_sall:="";
		if strings.EqualFold(arr[5],"买盘"){
			buy_sall="b"
		}else if strings.EqualFold(arr[5],"卖盘"){
			buy_sall="s"
		}else{
			buy_sall="e"
		}
		sql_arr=append(sql_arr,"('"+date+"',"+strconv.FormatInt(tmptime.Unix(),10)+","+arr[1]+","+zhangdiee+","+arr[3]+","+arr[4]+",'"+buy_sall+"')")
		count++
		if count>=200{
			_,err:=db.Raw(sql+strings.Join(sql_arr,",")).Exec()
			common.Checkerr(err)
			count=0
			sql_arr=sql_arr[(len(sql_arr)-1):]
		}
	}
	if len(sql_arr)>0{
		_,err:=db.Raw(sql+strings.Join(sql_arr,",")).Exec()
		common.Checkerr(err)
	}

	fmt.Println(date+"#####"+baseinfo.Code)
	return 1
}


