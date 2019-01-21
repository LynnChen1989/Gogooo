package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync/atomic"
)

type Params map[string]interface{}

var (
	_host string
	_api  *API
)

type Request struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Auth    string      `json:"auth,omitempty"` // omitempty表示如果为空置则忽略字段
	Id      int32       `json:"id"`
}

type Response struct {
	Jsonrpc string      `json:"jsonrpc"`
	Error   *ZbxError   `json:"error"`
	Result  interface{} `json:"result"`
	Id      int32       `json:"id"`
}

type ZbxError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// 实现打印错误信息方法
func (e *ZbxError) Error() string {
	return fmt.Sprintf("%d (%s): %s", e.Code, e.Message, e.Data)
}

type ExpectOneResult int

func (e *ExpectOneResult) Error() string {
	return fmt.Sprintf("Expectd exactly one result, got %d.", *e)
}

type ExpectMore struct {
	Expected int
	Got      int
}

func (e *ExpectMore) Error() string {
	return fmt.Sprintf("Expectd %d, got %d", e.Expected, e.Got)
}

type API struct {
	AuthToken string
	Logger    *log.Logger
	url       string
	c         http.Client
	id        int32
}

func NewAPI(url string) (api *API) {
	return &API{url: url, c: http.Client{}}
}

func (api *API) SetClient(c *http.Client) {
	api.c = *c
}

// 实现日志打印方法
func (api *API) Printf(format string, v ...interface{}) {
	if api.Logger != nil {
		api.Logger.Printf(format, v...)
	}
}

func (api *API) callBytes(method string, params interface{}) (b []byte, err error) {
	id := atomic.AddInt32(&api.id, 1)
	js := Request{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		Auth:    api.AuthToken,
		Id:      id,
	}
	b, err = json.Marshal(js)
	if err != nil {
		return
	}
	api.Printf("Request(%s): %s", "POST", b)
	req, err := http.NewRequest("POST", api.url, bytes.NewReader(b))
	if err != nil {
		return
	}

	// 增加请求的header
	req.ContentLength = int64(len(b))
	req.Header.Add("Content-Type", "application/json")

	// 发起请求
	res, err := api.c.Do(req)
	if err != nil {
		api.Printf("Error: %s", err)
		return
	}

	defer res.Body.Close()

	b, err = ioutil.ReadAll(res.Body)
	api.Printf("Response(%d): %s", res.StatusCode, b)
	return
}

// 封装调用方法
func (api *API) Call(method string, params interface{}) (response Response, err error) {
	b, err := api.callBytes(method, params)
	if err == nil {
		err = json.Unmarshal(b, &response)
	}
	return
}

func (api *API) CallWithError(method string, params interface{}) (response Response, err error) {
	response, err = api.Call(method, params)
	if err == nil && response.Error != nil {
		err = response.Error
	}
	return
}

func (api *API) Login(user, password string) (auth string, err error) {
	params := map[string]string{"user": user, "password": password}
	response, err := api.CallWithError("user.login", params)
	if err != nil {
		return
	}
	auth = response.Result.(string)
	api.AuthToken = auth
	return
}

func (api *API) Version() (v string, err error) {
	auth := api.AuthToken
	api.AuthToken = ""
	response, err := api.CallWithError("APIInfo.version", Params{})
	api.AuthToken = auth
	if e, ok := err.(*ZbxError); ok && e.Code == -32602 {
		response, err = api.CallWithError("APIInfo.version", Params{})
	}
	if err != nil {
		return
	}
	v = response.Result.(string)
	return
}

func getAPI() *API {
	if _api != nil {
		return _api
	}
	url, user, password := os.Getenv("ZABBIX_URL"), os.Getenv("ZABBIX_USER"), os.Getenv("ZABBIX_PASSWORD")
	_api = NewAPI(url)
	_api.SetClient(http.DefaultClient)
	v := os.Getenv("ZABBIX_VERBOSE")
	if v != "" && v != "0" {
		_api.Logger = log.New(os.Stderr, "[zabbix] ", 0)
	}

	if user != "" {
		auth, err := _api.Login(user, password)
		if err != nil {
			log.Fatal("Login Error")
		}
		if auth == "" {
			log.Fatal("Login Failed")
		}
	}

	return _api
}

func (api *API) HostGet(params Params) (result map[int]string, err error) {
	// 判断key为output是否存在
	output := [...]string{"hostid", "name"}
	if _, present := params["output"]; !present {
		params["output"] = output
	}
	response, err := api.CallWithError("host.get", params)

	if err != nil {
		return
	}

	tempResult := response.Result.([]interface{})
	result = make(map[int]string)
	for _, v := range tempResult {
		v := reflect.ValueOf(v)
		name := v.MapIndex(reflect.ValueOf("name")).Interface().(string)
		hostId := v.MapIndex(reflect.ValueOf("hostid")).Interface().(string)
		n, _ := strconv.Atoi(hostId)
		result[n] = name
	}
	return
}

func HostGet() {
	api := getAPI()
	_, err := api.Version()
	if err != nil {
		return
	}
}
