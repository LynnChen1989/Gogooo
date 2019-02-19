package main

import (
	"database/sql"
	"fmt"
	"os"
)

// root:fhl3mjsdwj@tcp(172.16.1.18:50001)/db_dap
// root:fhl3mjsdwj@tcp(172.16.1.18:50001)/cas
// root:fhl3mjsdwj@tcp(172.16.1.18:50001)/db_act
// 日切 日终判断

type CutDate struct {
}

type CutEnd struct {
}

type DbQueryCnt struct {
	cnt int
}

func CheckBatchStatus(DbInfo string, smt string) (rows *sql.Row) {
	Info.Printf("DBInfo:[%s]", DbInfo)
	Info.Printf("SQL:[%s] will been execute", smt)
	dbw := DbWorker{Dsn: DbInfo}
	db, err := sql.Open("mysql", dbw.Dsn)
	err = db.Ping()
	if err != nil {
		Error.Printf("Failed to ping mysql: %s", err)
	}
	if err != nil {
		panic(err)
		return
	}
	rows = db.QueryRow(smt)
	db.Close()
	return rows
}

func (cd *CutDate) CheckCutDateStatus() (cutDateStatus bool) {
	// 日切状态检查
	var cnt int
	dapDbInfo := os.Getenv("DAP_DB_INFO")
	casDbInfo := os.Getenv("CAS_DB_INFO")
	sqlLan01 := fmt.Sprintf("select count(1) as cnt from t_batch_group_execution where group_id='GCasDayCutGroup' and status='2' and DATE_FORMAT(end_time,'yyyy-MM-dd')=DATE_FORMAT(CURRENT_DATE,'yyyy-MM-dd')")
	//sqlLan02 := fmt.Sprintf("select count(1) as cnt from t_sys_conf where ACT_DATE = CURRENT_DATE and ACT_STAT = '2' and TXN_DATE = CURRENT_DATE and TXN_STAT = '2'")
	sqlLan02 := fmt.Sprintf("select count(1) as cnt from t_sys_conf where ACT_DATE = '2024-10-07' and ACT_STAT = '2' and TXN_DATE = '2024-10-07' and TXN_STAT = '2'")
	// 序列服务
	dapRows := CheckBatchStatus(dapDbInfo, sqlLan01)
	dapRows.Scan(&cnt)
	if cnt > 0 {
		cutDateStatus = true
		Info.Printf("日切检查 step[1] check, success condition [cnt>0], cnt: %d,success", cnt)
	} else {
		cutDateStatus = false
		Error.Printf("日切检查 step[1] check, success condition [cnt>0], cnt: %d, failure", cnt)
	}
	// 信贷核心
	casRows := CheckBatchStatus(casDbInfo, sqlLan02)
	casRows.Scan(&cnt)
	if cnt > 0 {
		cutDateStatus = true
		Info.Printf("日切检查 step[2] check, success condition [cnt>0], cnt: %d,success", cnt)
	} else {
		cutDateStatus = false
		Error.Printf("日切检查 step[2] check, success condition [cnt>0], cnt: %d, failure", cnt)
	}
	//测试伪造状态
	//cutDateStatus = true
	Info.Printf("日切检查最终状态: [%t]", cutDateStatus)
	return
}

func (cd *CutDate) CheckCutEndStatus() (cutEndStatus bool) {
	// 检查日终状态
	var cnt int
	dapDbInfo := os.Getenv("DAP_DB_INFO")
	if dapDbInfo == "" {
		Error.Println("environment variable DAP_DB_INFO is needed")
		return
	}

	casDbInfo := os.Getenv("CAS_DB_INFO")
	if casDbInfo == "" {
		Error.Println("environment variable CAS_DB_INFO is needed")
		return
	}

	sqlLan01 := fmt.Sprintf("select count(1) from t_batch_group_execution where status <> '2'")
	sqlLan02 := fmt.Sprintf("select count(1) from t_act_system_para where ACCT_DATE = '2024-10-07' and stat = '1'")
	sqlLan03 := fmt.Sprintf("select count(1) from t_sys_conf where ACT_DATE = '2024-10-07' and ACT_STAT = '0' and TXN_DATE = '2024-10-07' and TXN_STAT = '0'")

	dapRows := CheckBatchStatus(dapDbInfo, sqlLan01)
	dapRows.Scan(&cnt)
	if cnt == 0 {
		cutEndStatus = true
		Info.Printf("日终检查 step[1] check, success condition [cnt=0], cnt: %d,success", cnt)
	} else {
		cutEndStatus = false
		Error.Printf("cut end step[1] check, success condition [cnt=0], cnt: %d, failure", cnt)
	}
	// 会计核算
	actRows := CheckBatchStatus(casDbInfo, sqlLan02)
	actRows.Scan(&cnt)
	if cnt == 1 {
		cutEndStatus = true
		Info.Printf("日终检查 step[2] check, success condition [cnt=1], cnt: %d,success", cnt)
	} else {
		cutEndStatus = false
		Error.Printf("日终检查 step[2] check, success condition [cnt=1], cnt: %d, failure", cnt)
	}

	// 信贷核心
	casRow := CheckBatchStatus(casDbInfo, sqlLan03)
	casRow.Scan(&cnt)
	if cnt == 1 {
		cutEndStatus = true
		Info.Printf("日终检查 step[3] check, success condition [cnt=1], cnt: %d,success", cnt)
	} else {
		cutEndStatus = false
		Error.Printf("日终检查 step[3] check, success condition [cnt=1], cnt: %d, failure", cnt)
	}

	//测试伪造状态
	//cutEndStatus = true
	Info.Printf("日终检查最终状态: [%t]", cutEndStatus)

	return
}
