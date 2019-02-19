package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
		SendNotify("general", "ODS抽数:关闭服务入口SUCCESS"+NowFormatDate("20060102"), "[API]关闭服务入口成功")
	} else if status == "failed" {
		SendNotify("fatal", "ODS抽数:关闭服务入口FAILURE"+NowFormatDate("20060102"), "[API]关闭服务入口失败")
	}
	var rs APIResponse
	rs.Data = "停服" + status
	ret, _ := json.Marshal(rs)
	fmt.Fprint(w, string(ret))
}

func StartSrvHandler(w http.ResponseWriter, r *http.Request) {
	status := RestoreSrv()
	if status == "success" {
		SendNotify("general", "ODS抽数:开启服务入口SUCCESS"+NowFormatDate("20060102"), "[API]开启服务入口成功")
	} else if status == "failed" {
		SendNotify("fatal", "ODS抽数:开启服务入口FAILURE"+NowFormatDate("20060102"), "[API]开启服务入口失败")
	}
	var rs APIResponse
	rs.Data = "启服" + status
	ret, _ := json.Marshal(rs)
	fmt.Fprint(w, string(ret))
}

func MQProductHandler(w http.ResponseWriter, r *http.Request) {
	PushCutBatchMsg()
	var rs APIResponse
	rs.Data = "通知下游服务消息已发送"
	ret, _ := json.Marshal(rs)
	fmt.Fprint(w, string(ret))
}

func CutDateHandler(w http.ResponseWriter, r *http.Request) {
	cd := &CutDate{}
	cutDate := cd.CheckCutDateStatus()
	cutEnd := cd.CheckCutEndStatus()
	var rs APIResponse
	rs.Data = "日切:" + strconv.FormatBool(cutDate) + "," + "日终:" + strconv.FormatBool(cutEnd)
	ret, _ := json.Marshal(rs)
	fmt.Fprint(w, string(ret))
}

func SlaveHandler(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	var rs APIResponse
	system := vars.Get("system")
	opt := vars.Get("opt")
	if system == "" && opt == "" {
		rs.Data = "system或opt参数缺失"
	} else {
		BatchHandleDbSlave(system, opt)
	}
	ret, _ := json.Marshal(rs)
	fmt.Fprint(w, string(ret))
}

func AlertHandler(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	opt := vars.Get("opt")
	api := getAPI()
	var rs APIResponse
	if opt == "" {
		rs.Data = "opt参数缺失"
	} else {
		if opt == "start" {
			api.RestoreAlert()
			rs.Data = "已解除告警屏蔽"
		} else if opt == "stop" {
			api.PauseAlert()
			rs.Data = "已屏蔽告警"
		}
	}
	ret, _ := json.Marshal(rs)
	fmt.Fprint(w, string(ret))
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
	http.HandleFunc("/cut", CutDateHandler)
	http.HandleFunc("/slave", SlaveHandler)
	http.HandleFunc("/alert", AlertHandler)
	http.ListenAndServe("0.0.0.0:8001", nil)
	Info.Println("HttpServer Running ON 0.0.0.0:8001 ...")
}
