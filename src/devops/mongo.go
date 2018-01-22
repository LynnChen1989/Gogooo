package main

import (
	_ "fmt"
	"gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"
)

func getDB() *mgo.Database {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("test")
	return db
}
