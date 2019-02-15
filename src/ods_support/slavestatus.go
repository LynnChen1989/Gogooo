package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"regexp"
	"strings"
)

type ReplStatus struct {
	Status string `db:"repl_status"`
}

type OptSLave struct{}

type DbWorker struct {
	Dsn string
}

func getMysqlDsnIpPort(Dns string) (ipPort string) {
	pat := `(?s)\((.*)\)`
	reg := regexp.MustCompile(pat)
	ipPort = string(reg.Find([]byte(Dns)))
	return
}

func BatchHandleDbSlave(status string) {
	/*
		批量处理主从的断开和恢复，status可选参数 start, stop
	*/
	slaveList := os.Getenv("SLAVE_LIST")
	var allSlaveStatus string
	for _, mysqlDns := range strings.Split(slaveList, "##") {
		dbw := DbWorker{
			Dsn: mysqlDns,
		}
		db, err := sql.Open("mysql", dbw.Dsn)
		if err != nil {
			panic(err)
			return
		}
		a := &OptSLave{}
		if status == "stop" {
			Info.Printf("now stop slave on [%s]", mysqlDns)
			//a.StopSlave(db)
			slaveStatus := a.GetSLaveStatus(db)
			allSlaveStatus += getMysqlDsnIpPort(mysqlDns) + "SLAVE:" + slaveStatus + ","
		} else if status == "start" {
			Info.Printf("now start slave on [%s]", mysqlDns)
			//a.StartSlave(db)
			slaveStatus := a.GetSLaveStatus(db)
			allSlaveStatus += getMysqlDsnIpPort(mysqlDns) + "SLAVE:" + slaveStatus + ","
		}
	}
	SendNotify("general", "ODS抽数:主从状态", allSlaveStatus)
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
	replStatus = rs.Status
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
