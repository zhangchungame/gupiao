package rimingxi

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Rimingxiget() string {
	req, err := http.Get("http://quotes.money.163.com/service/chddata.html?code=0600321&start=20160524&end=20160525&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;TURNOVER;VOTURNOVER;VATURNOVER;TCAP;MCAP")
	if err != nil {
		fmt.Println(err)
	}
	robots, err := ioutil.ReadAll(req.Body)
	str := string(robots)
	req.Body.Close()
	return str
}
