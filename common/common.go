package common

import "fmt"

func Checkerr(err error) {
	if err!=nil {
		fmt.Println(err)
	}
}