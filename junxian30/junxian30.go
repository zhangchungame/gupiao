package junxian30

import (
	"gupiao/allcode"
	"gopkg.in/mgo.v2"
	"gupiao/singleInstance"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"gupiao/rikxian"
)

var collection *mgo.Collection

func init()  {
	collection=singleInstance.GetMongoInstance().C("rikxian")
}


func CalculateAll()  {
	codes := allcode.Getallcodes()
	for _, val := range codes {
		junxian30(val.Code)
		//go junxian30(val.Code)
	}
}

func junxian30(code string)  {
	var data  []rikxian.Rikxian
	collection.Find(bson.M{"code":code}).Sort("-data_int")..All(&data)
	//collection.Find(bson.M{"code":code}).All(&data)
	fmt.Println(data)
}
