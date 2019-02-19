package main

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Today() (today string) {
	currentTime := time.Now()
	today = currentTime.Format("20060102")
	return
}

func EnvVariablesInitCheck() {
	envAll := [...]string{
		"CAS_CMS_SLAVE_LIST",
		"ACT_SLAVE_LIST",
		"DAP_DB_INFO",
		"CAS_DB_INFO",
		"MESSAGE_URL",
		"YHB_REDIS_HOST",
		"YHB_REDIS_PASSWORD",
		"YHB_REDIS_DB",
		"ZABBIX_URL",
		"ZABBIX_USER",
		"ZABBIX_PASSWORD",
		"MQ_URI",
		"CACHE_REDIS_HOST",
		"CACHE_REDIS_PASSWORD",
		"CACHE_REDIS_DB",
		"CRON_STOP_SRV",
		"CRON_CUT_DATE",
		"CRON_CUT_END",
		"CRON_RESTORE_SRV"}

	for _, e := range envAll {
		envValue := os.Getenv(e)
		if envValue == "" {
			Error.Printf("environment variable [%s] is needed", e)
			panic("you must input all environment variables")
		} else {
			Info.Printf("get environment %s, value: [%s]", e, envValue)
		}
	}
}

func LogInit() {
	logFile, err := os.OpenFile("ods.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("open log file error: ", err)
	}

	Info = log.New(io.MultiWriter(os.Stderr, logFile),
		"[INFO] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(io.MultiWriter(os.Stderr, logFile),
		"[WARNING] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(io.MultiWriter(os.Stderr, logFile),
		"[ERROR] ",
		log.Ldate|log.Ltime|log.Lshortfile)

}
