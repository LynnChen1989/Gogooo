package main

import (
	_ "../github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	_ "reflect"
)

//定义数据库连接信息的一个结构体
type DbInfo struct {
	driver    string
	dbUser    string
	dbPass    string
	dbHost    string
	dbPort    int
	dbName    string
	character string
}

func (di *DbInfo) ConnectInfo() string {
	//ConnectInfo方法，主要是构造一个连接信息的字符串
	strFormat := `%s:%s@tcp(%s:%d)/%s?charset=%s`
	return fmt.Sprintf(strFormat, di.dbUser, di.dbPass, di.dbHost, di.dbPort, di.dbName, di.character)
}

func Conn() (*sql.DB, error) {
	DI := DbInfo{
		driver:    "mysql",
		dbUser:    "roo",
		dbPass:    "",
		dbHost:    "127.0.0.1",
		dbPort:    3306,
		dbName:    "db_rms",
		character: "utf8",
	}
	db, err := sql.Open(DI.driver, DI.ConnectInfo())
	return db, err
}

func main() {
	// 初始化连接信息
	db, _ := Conn()
	row := db.QueryRow("SELECT * FROM rms_cloud")
	var col1 string
	row.Scan(&col1)
	fmt.Println(col1)
}
