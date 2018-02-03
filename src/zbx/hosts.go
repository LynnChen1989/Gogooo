package zbx

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
)

var (
	_host string
	_api  *API
)

type Items struct {
	Hosts []string `json:"hosts"`
}

type Group struct {
	GroupId string `json:"hostid"`
	Name    string `json:"name"`
}

func init() {
	rand.Seed(time.Now().UnixNano())

	var err error
	_host, err = os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	_host += "-testing"
	if os.Getenv("ZABBIX_URL") == "" {
		log.Fatal("Set environment variables ZABBIX_URL (and optionally ZABBIX_USER and ZABBIX_PASSWORD)")
	}
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

// 根据host id查找host group
func (api *API) GetGroupByID(hid int) (name string, err error) {
	params := Params{"hostids": hid}
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("hostgroup.get", params)
	if err != nil {
		return
	}
	tempResult := response.Result.([]interface{})
	if len(tempResult) > 1 {
		e := ExpectOneResult(len(tempResult))
		err = &e
		api.Printf("HostId %d find Multi Group", hid)
	} else {
		for _, v := range tempResult {
			v := reflect.ValueOf(v)
			name = v.MapIndex(reflect.ValueOf("name")).Interface().(string)
		}
	}
	return
}

// 根据host id查找 ip
func (api *API) GetInterfaceById(hid int) (ip string, err error) {
	response, err := api.CallWithError("hostinterface.get", Params{"hostids": hid})
	if err != nil {
		return
	}
	tempResult := response.Result.([]interface{})
	if len(tempResult) != 1 {
		e := ExpectOneResult(len(tempResult))
		err = &e
		api.Printf("HostId %d find Multi IpAddress", len(tempResult))
	} else {
		for _, v := range tempResult {
			v := reflect.ValueOf(v)
			ip = v.MapIndex(reflect.ValueOf("ip")).Interface().(string)
		}
	}
	return
}

//func main() {
//	api := getAPI()
//	_, err := api.Version()
//	if err != nil {
//		return
//	}
//
//	hosts, err := api.HostGet(Params{"groupids": 51})
//	if err != nil {
//		return
//	}
//
//	var result map[string]interface{}
//	result = make(map[string]interface{})
//	for hid, name := range hosts {
//		api.printf("%d:%s", hid, name)
//		group, _ := api.GetGroupByID(hid)
//		ip, _ := api.GetInterfaceById(hid)
//
//		if item, ok := result[group].(*Items); ok {
//			result[group].(*Items).Hosts = append(item.Hosts, ip)
//		} else {
//			it := &Items{Hosts: []string{ip}}
//			result[group] = it
//		}
//	}
//	data, err := json.MarshalIndent(result, "", "\t")
//	fmt.Println(string(data))
//}
