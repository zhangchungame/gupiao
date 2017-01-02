package singleInstance

import (
	"fmt"
	"gopkg.in/olivere/elastic.v3"
	"sync"
)



var client *elastic.Client
var esonce sync.Once

func GetEsInstance() *elastic.Client {
	esonce.Do(func() {
		var err error
		client, err = elastic.NewClient(elastic.SetURL("http://localhost:9200"))
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		is_exist,err:=client.IndexExists("gupiao").Do()
		if !is_exist {
			client.CreateIndex("gupiao").Do()
		}
	})
	return client
}