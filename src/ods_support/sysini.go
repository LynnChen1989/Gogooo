package main

import (
	"io"
	"log"
	"os"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func envVariablesInitCheck() {
	envAll := [...]string{"SLAVE_LIST", "DAP_DB_INFO", "CAS_DB_INFO", "MESSAGE_URL", "REDIS_HOST",
		"REDIS_PASSWORD", "REDIS_DB", "ZABBIX_URL", "ZABBIX_USER", "ZABBIX_PASSWORD", "MQ_URI", "JUDGE_ODS_DB_INFO"}

	for _, e := range envAll {
		envValue := os.Getenv(e)
		if envValue == "" {
			Error.Printf("environment variable [%s] is needed", e)
			return
		} else {
			Info.Printf("get environment %s, value: [%s]", e, envValue)
		}
	}
}

func logInit() {
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
