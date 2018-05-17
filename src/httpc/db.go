package main

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"os"
	"gopkg.in/mgo.v2/bson"
)

const MONGO_URI = "mongodb://db_go:db_go@172.16.254.224:27017/admin"

func dbSession() interface{} {
	session, err := mgo.Dial(MONGO_URI)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(0)
	}
	return session
}

func getOneByCondition(db string, tb string, condition map[string]interface{}, result interface{}) (err error) {
	/*
	USE CASE:
		result := &Person{}
		getOneByCondition("db_go", "tb_go", condition, result)
		fmt.Println(result.Name)
	*/
	ss := dbSession()
	session, ok := ss.(*mgo.Session)
	if !ok {
		panic("db session error.")
	}
	return session.DB(db).C(tb).Find(bson.M(condition)).One(result)
}

func getManyByCondition(db string, tb string, condition map[string]interface{}, result interface{}) (err error) {
	/*
	USE CASE:
		var result []Person
		getManyByCondition("db_go", "tb_go", condition, &result)
		fmt.Println(result)
	*/
	ss := dbSession()
	session, ok := ss.(*mgo.Session)
	if !ok {
		panic("db session error")
	}
	return session.DB(db).C(tb).Find(bson.M(condition)).All(result)
}

func Insert(db string, tb string, doc interface{}) (err error) {
	/*
	USE CASE:
		p := &Person{
			Name: "chenboyu",
			Age:  1,
			Sex:  "male",
		}
	Insert("db_go", "tb_go", p)
	*/
	ss := dbSession()
	session, ok := ss.(*mgo.Session)
	if !ok {
		panic("db session error")
	}
	return session.DB(db).C(tb).Insert(doc)
}

func Update(db string, tb string, condition map[string]interface{}, data map[string]interface{}) (change interface{}, err error) {
	/*
	满足条件的批量更新
	USE CASE:

	condition := map[string]interface{}{
		"name": "zhangmeixi",
	}
	value := map[string]interface{}{
		"name": "chenlin2",
	}
	Update("db_go", "tb_go", condition, value)
	*/
	ss := dbSession()
	session, ok := ss.(*mgo.Session)
	if !ok {
		panic("db session error")
	}

	selector := bson.M(condition)
	value := bson.M{"$set": data}
	return session.DB(db).C(tb).UpdateAll(selector, value)
}

func DeleteByCondition(db string, tb string, condition map[string]interface{}) (change interface{}, err error) {
	ss := dbSession()
	session, ok := ss.(*mgo.Session)
	if !ok {
		panic("db session error")
	}
	return session.DB(db).C(tb).RemoveAll(condition)
}
