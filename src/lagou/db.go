package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "reflect"
	//"log"
)

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
	strFormat := `%s:%s@tcp(%s:%d)/%s?charset=%s`
	return fmt.Sprintf(strFormat, di.dbUser, di.dbPass, di.dbHost, di.dbPort, di.dbName, di.character)
}

func Conn() (db *sql.DB, err error) {
	DI := &DbInfo{
		driver:    "mysql",
		dbUser:    "root",
		dbPass:    "root",
		dbHost:    "127.0.0.1",
		dbPort:    13306,
		dbName:    "db_dnsdog",
		character: "utf8",
	}
	db, err = sql.Open(DI.driver, DI.ConnectInfo())
	return

}

//func main() {
//	conn, err := Conn()
//	if err != nil {
//		log.Fatal("database connect error:", err)
//	}
//	rows, err := conn.Query("SELECT * FROM applydns_cloud")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for rows.Next() {
//		var id int
//		var name string
//		var code string
//		if err := rows.Scan(&id, &name, &code); err != nil {
//			log.Fatal(err)
//		}
//		fmt.Printf("%s is %d\n", name)
//	}
//}
