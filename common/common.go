package common

import (
	"fmt"
	"bytes"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
)

func Checkerr(err error) {
	if err!=nil {
		fmt.Println(err)
	}
}

func Gb2utf_decode(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	//O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	O := transform.NewReader(I, simplifiedchinese.GB18030.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

