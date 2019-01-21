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

func Init() {

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

//func IndexHandler(w http.ResponseWriter, r *http.Request) {
//	HandleDbSlave()
//	fmt.Fprintln(w, "hello")
//}

func main() {
	//fmt.Println("HttpServer Running ON 0.0.0.0:80001 ...")
	//http.HandleFunc("/", IndexHandler)
	//http.ListenAndServe("0.0.0.0:8001", nil)
	Init()
	//StopSrv()
	//PauseTimer()
	//HandleDbSlave()
	//RedisClient()
	//httpPost("http://122.152.209.199:2046/api/v1/message/")
	//SendNotify("test:test", "test")
	HostGet()
}
