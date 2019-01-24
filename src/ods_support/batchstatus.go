package main

import (
	"database/sql"
	"fmt"
	"os"
)

// root:fhl3mjsdwj@tcp(172.16.1.18:50001)/db_dap
// root:fhl3mjsdwj@tcp(172.16.1.18:50001)/cas
// root:fhl3mjsdwj@tcp(172.16.1.18:50001)/db_act
func CheckBatchStatus(DbInfo string, smt string) (rows *sql.Row) {
	dbw := DbWorker{Dsn: DbInfo}
	db, err := sql.Open("mysql", dbw.Dsn)
	if err != nil {
		panic(err)
		return
	}
	rows = db.QueryRow(smt)
	db.Close()
	return rows
}

func CheckGlobalStatus() {
	sqlLan := fmt.Sprintf("select * from t_batch_group_execution where status in ('3') and start_time > date_sub(now(), interval 30 minute);")
	DbInfo := os.Getenv("DAP_DB_INFO")
	if DbInfo == "" {
		Error.Println("environment variable DAP_DB_INFO is needed")
		return
	}
	CheckBatchStatus(DbInfo, sqlLan)
}

func CheckCASStat() {
	sqlLan := fmt.Sprintf("select count(1) from t_sys_conf where ACT_DATE = CURRENT_DATE and ACT_STAT = '0' and TXN_DATE = CURRENT_DATE and TXN_STAT = '0' and ACT_MODIFY_TIME > date_sub(now(), interval 30 minute) and TXN_MODIFY_TIME > date_sub(now(), interval 30 minute)")
	DbInfo := os.Getenv("CAS_DB_INFO")
	if DbInfo == "" {
		Error.Println("environment variable CAS_DB_INFO is needed")
		return
	}
	CheckBatchStatus(DbInfo, sqlLan)
}

func CheckACTStat() {
	sqlLan := fmt.Sprintf("select count(1) from t_act_system_para where ACCT_DATE = CURRENT_DATE and stat = '1' and MODIFY_TIME > date_sub(now(), interval 30 minute);")
	DbInfo := os.Getenv("ACT_DB_INFO")
	if DbInfo == "" {
		Error.Println("environment variable ACT_DB_INFO is needed")
		return
	}
	CheckBatchStatus(DbInfo, sqlLan)
}

func CheckDAPStat() {
	sqlLan := fmt.Sprintf("select count(1) from t_batch_group_execution where status <> '2' and start_time > date_sub(now(), interval 30 minute);")
	DbInfo := os.Getenv("DAP_DB_INFO")
	if DbInfo == "" {
		Error.Println("environment variable DAP_DB_INFO is needed")
		return
	}
	CheckBatchStatus(DbInfo, sqlLan)
}
