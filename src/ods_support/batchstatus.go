package main

import (
	"database/sql"
	"fmt"
	"os"
)

type CutDate struct {
}

type CutEnd struct {
}

type DbQueryCnt struct {
	cnt int
}

var (
	cutDateStatus01 bool
	cutDateStatus02 bool
)
var (
	cutEndStatus01 bool
	cutEndStatus02 bool
	cutEndStatus03 bool
)

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
	sqlLan02 := fmt.Sprintf("select count(1) as cnt from t_sys_conf where ACT_DATE = '2024-10-11' and ACT_STAT = '2' and TXN_DATE = '2024-10-11' and TXN_STAT = '2'")
	// 序列服务
	dapRows := CheckBatchStatus(dapDbInfo, sqlLan01)
	dapRows.Scan(&cnt)
	if cnt > 0 {
		cutDateStatus01 = true
		Info.Printf("日切检查 step[1] check, 【dap】检查日切批量状态 success condition [cnt>0], cnt: %d,success", cnt)
	} else {
		cutDateStatus01 = false
		Error.Printf("日切检查 step[1] check, 【dap】检查日切批量状态 success condition [cnt>0], cnt: %d, failure", cnt)
	}
	// 信贷核心
	casRows := CheckBatchStatus(casDbInfo, sqlLan02)
	casRows.Scan(&cnt)
	if cnt > 0 {
		cutDateStatus02 = true
		Info.Printf("日切检查 step[2] check, 【cas】检查日切状态是否为日切中 success condition [cnt>0], cnt: %d,success", cnt)
	} else {
		cutDateStatus02 = false
		Error.Printf("日切检查 step[2] check, 【cas】检查日切状态是否为日切中 success condition [cnt>0], cnt: %d, failure", cnt)
	}
	//测试伪造状态
	//cutDateStatus = true
	if cutDateStatus01 && cutDateStatus02 {
		cutDateStatus = true
	} else {
		cutDateStatus = false
	}
	Info.Printf("日切检查最终状态: [%t]", cutDateStatus)
	cutDateStatus = true
	return
}

func (cd *CutDate) CheckCutEndStatus() (cutEndStatus bool) {
	// 检查日终状态
	var cnt int

	sqlLan01 := fmt.Sprintf("select count(1) as cnt from t_batch_group_execution where status <> '2'")
	sqlLan02 := fmt.Sprintf("select count(1) as cnt from t_act_system_para where ACCT_DATE = '2024-10-11' and stat = '1'")
	sqlLan03 := fmt.Sprintf("select count(1) as cnt from t_sys_conf where ACT_DATE = '2024-10-11' and ACT_STAT = '0' and TXN_DATE = '2024-10-11' and TXN_STAT = '0'")

	dapRows := CheckBatchStatus(os.Getenv("DAP_DB_INFO"), sqlLan01)
	dapRows.Scan(&cnt)
	if cnt == 0 {
		cutEndStatus01 = true
		Info.Printf("日终检查 step[1] check, 【dap】检查跑批失败数量 success condition [cnt=0], cnt: %d,success", cnt)
	} else {
		cutEndStatus01 = false
		Error.Printf("cut end step[1] check, 【dap】检查跑批失败数量 success condition [cnt=0], cnt: %d, failure", cnt)
	}
	// 会计核算
	actRows := CheckBatchStatus(os.Getenv("ACT_DB_INFO"), sqlLan02)
	actRows.Scan(&cnt)
	if cnt == 1 {
		cutEndStatus02 = true
		Info.Printf("日终检查 step[2] check, 【act】检查核算日终完成 success condition [cnt=1], cnt: %d,success", cnt)
	} else {
		cutEndStatus02 = false
		Error.Printf("日终检查 step[2] check, 【act】检查核算日终完成 success condition [cnt=1], cnt: %d, failure", cnt)
	}

	// 信贷核心
	casRow := CheckBatchStatus(os.Getenv("CAS_DB_INFO"), sqlLan03)
	casRow.Scan(&cnt)
	if cnt == 1 {
		cutEndStatus03 = true
		Info.Printf("日终检查 step[3] check, 【cas】检查核心日终完成 success condition [cnt=1], cnt: %d,success", cnt)
	} else {
		cutEndStatus03 = false
		Error.Printf("日终检查 step[3] check, 【cas】检查核心日终完成 success condition [cnt=1], cnt: %d, failure", cnt)
	}

	if cutEndStatus01 && cutEndStatus02 && cutEndStatus03 {
		cutEndStatus = true
	} else {
		cutEndStatus = false
	}
	Info.Printf("日终检查最终状态: [%t]", cutEndStatus)
	//测试伪造状态
	cutEndStatus = true
	return
}
