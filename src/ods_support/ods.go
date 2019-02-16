package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIResponse struct {
	Data string `json:"data"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "application/json")
	var rs APIResponse
	rs.Data = "Hello, ODS"
	ret, _ := json.Marshal(rs)
	fmt.Fprint(w, string(ret))
}

func StopSrvHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "application/json")
	status := PauseSrv()
	if status == "success" {
		SendNotify("general", "ODS抽数:关闭服务入口SUCCESS", "[API]关闭服务入口成功")
	} else if status == "failed" {
		SendNotify("fatal", "ODS抽数:关闭服务入口FAILURE", "[API]关闭服务入口失败")
	}
	var rs APIResponse
	rs.Data = "停服" + status
	ret, _ := json.Marshal(rs)
	fmt.Fprint(w, string(ret))
}

func StartSrvHandler(w http.ResponseWriter, r *http.Request) {
	status := RestoreSrv()
	if status == "success" {
		SendNotify("general", "ODS抽数:开启服务入口SUCCESS", "[API]开启服务入口成功")
	} else if status == "failed" {
		SendNotify("fatal", "ODS抽数:开启服务入口FAILURE", "[API]开启服务入口失败")
	}
	var rs APIResponse
	rs.Data = "启服" + status
	ret, _ := json.Marshal(rs)
	fmt.Fprint(w, string(ret))
}

func MQProductHandler(w http.ResponseWriter, r *http.Request) {
	PushCutBatchMsg()
}

//
//func CutDateHandler(w http.ResponseWriter, r *http.Request) {
//	cd := &CutDate{}
//	cd.CheckCutDateStatus()
//	cd.CheckCutEndStatus()
//}

//func SlaveHandler(w http.ResponseWriter, r *http.Request) {
//	BatchHandleDbSlave("cascms", "stop")
//}

func StopAlertHandler(w http.ResponseWriter, r *http.Request) {
	api := getAPI()
	Info.Printf("Mid: [%d]", api.PauseAlert())
}

func main() {
	LogInit()
	EnvVariablesInitCheck()
	go OdsTimer()
	go ReceiveMsg()
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/stopsrv", StopSrvHandler)
	http.HandleFunc("/startsrv", StartSrvHandler)
	http.HandleFunc("/push", MQProductHandler)
	//http.HandleFunc("/cutdate", CutDateHandler)
	//http.HandleFunc("/ss", SlaveHandler)
	http.HandleFunc("/stopalert", StopAlertHandler)
	http.ListenAndServe("0.0.0.0:8001", nil)
	Info.Println("HttpServer Running ON 0.0.0.0:8001 ...")
}
