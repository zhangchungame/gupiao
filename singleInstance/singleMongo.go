package singleInstance

import (
	"gopkg.in/mgo.v2"
	"sync"
	"fmt"
)

var mdb *mgo.Database
var mdb_once sync.Once

func GetMongoInstance() *mgo.Database {
	mdb_once.Do(func() {
		var err error
		session, err := mgo.Dial("localhost")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		session.SetMode(mgo.Monotonic, true)

		mdb = session.DB("gupiao")
	})
	return mdb
}
