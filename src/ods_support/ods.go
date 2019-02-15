package main

import (
	"fmt"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	//HandleDbSlave()
	fmt.Fprintln(w, "hello")
}

func StopSrvHandler(w http.ResponseWriter, r *http.Request) {
	status := PauseSrv()
	if status == "success" {
		SendNotify("general", "ODS抽数:关闭服务入口SUCCESS", "关闭服务入口成功")
	} else if status == "failed" {
		SendNotify("fatal", "ODS抽数:关闭服务入口FAILURE", "关闭服务入口失败")
	}
	fmt.Fprintln(w, status)
}

func StartSrvHandler(w http.ResponseWriter, r *http.Request) {
	status := RestoreSrv()
	if status == "success" {
		SendNotify("general", "ODS抽数:开启服务入口SUCCESS", "开启服务入口成功")
	} else if status == "failed" {
		SendNotify("fatal", "ODS抽数:开启服务入口FAILURE", "开启服务入口失败")
	}
	fmt.Fprintln(w, status)
}

//func MQConsumeHandler(w http.ResponseWriter, r *http.Request) {
//	receiveMsg()
//}

func MQProductHandler(w http.ResponseWriter, r *http.Request) {
	pushMsg()
}

func CutDateHandler(w http.ResponseWriter, r *http.Request) {
	cd := &CutDate{}
	cd.CheckCutDateStatus()
	cd.CheckCutEndStatus()
}

func SlaveHandler(w http.ResponseWriter, r *http.Request) {
	BatchHandleDbSlave("stop")
}

func main() {
	logInit()
	envVariablesInitCheck()
	go receiveMsg()
	Info.Println("HttpServer Running ON 0.0.0.0:8001 ...")
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/stopsrv", StopSrvHandler)
	http.HandleFunc("/startsrv", StartSrvHandler)
	//http.HandleFunc("/consume-mq", MQConsumeHandler)
	http.HandleFunc("/push-mq", MQProductHandler)
	http.HandleFunc("/cutdate", CutDateHandler)
	http.HandleFunc("/ss", SlaveHandler)
	http.ListenAndServe("0.0.0.0:8001", nil)

	//Init()
	//StopSrv()
	//PauseTimer()
	//HandleDbSlave()
	//RedisClient()
	//httpPost("http://122.152.209.199:2046/api/v1/message/")
	//SendNotify("test:test", "test")
	//HostGet()
	//receiveMsg()
}
