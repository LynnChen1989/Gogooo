package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type ReplStatus struct {
	Status string `db:"repl_status"`
}

type OptSLave struct{}

type DbWorker struct {
	Dsn string
}

func HandleDbSlave() {
	dbw := DbWorker{
		Dsn: "root:fhl3mjsdwj@tcp(172.16.1.18:50001)/db_cms01",
	}
	db, err := sql.Open("mysql", dbw.Dsn)
	if err != nil {
		panic(err)
		return
	}
	a := &OptSLave{}
	a.GetSLaveStatus(db)
}

func (os *OptSLave) GetSLaveStatus(DB *sql.DB) (replStatus string) {
	sqlLan := fmt.Sprintf("SELECT variable_value AS repl_status FROM information_schema.global_status  WHERE variable_name='SLAVE_RUNNING'")
	rows := DB.QueryRow(sqlLan)
	rs := new(ReplStatus)
	if err := rows.Scan(&rs.Status); err != nil {
		Error.Println("scan data error:", err)
		return
	}
	DB.Close()
	Info.Println("current replication status is: ", rs.Status)
	return
}

func (os *OptSLave) StopSlave(DB *sql.DB) {
	sqlLan := fmt.Sprintf("STOP SLAVE")
	_, err := DB.Query(sqlLan)
	if err != nil {
		Error.Println("stop salve status error:", err)
		return
	}
	DB.Close()
	Info.Println("stop slave success")
}

func (os *OptSLave) StartSlave(DB *sql.DB) {
	sqlLan := fmt.Sprintf("START SLAVE")
	_, err := DB.Query(sqlLan)
	if err != nil {
		Error.Println("start salve status error:", err)
		return
	}
	DB.Close()
	Info.Println("start slave success")
}
