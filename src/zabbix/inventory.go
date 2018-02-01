package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func (api *API) GetGroup(params Params) (result map[int]string, err error) {
	// 返回数据类型：{group_id: groupName}
	output := []string{"groupid", "name"}
	if _, present := params["output"]; !present {
		params["out"] = output
	}

	response, err := api.CallWithError("hostgroup.get", params)

	if err != nil {
		return
	}

	tempResult := response.Result.([]interface{})
	result = make(map[int]string)
	for _, v := range tempResult {
		v := reflect.ValueOf(v)
		name := v.MapIndex(reflect.ValueOf("name")).Interface().(string)
		hostId := v.MapIndex(reflect.ValueOf("groupid")).Interface().(string)
		n, _ := strconv.Atoi(hostId)
		result[n] = name
	}
	return
}

// 根据主机组名获取主机
func (api *API) GetHostByGroup(group string) (results map[int]string) {
	groups, err := api.GetGroup(Params{})
	if err != nil {
		return
	}
	for gid, name := range groups {
		if strings.Trim(name, "") != strings.Trim(group, "") {
			continue
		}
		hosts, err := api.HostGet(Params{"groupids": gid})
		if err != nil {
			return
		}
		results = hosts
	}
	return
}

func (api *API) PrintResult(hosts map[int]string, group string) {
	result := make(map[string]interface{})

	for hid, name := range hosts {
		api.printf("%d:%s", hid, name)
		ip, err := api.GetInterfaceById(hid)
		if err != nil {
			return
		}
		/*
			下面这段主要解决慢的问题
			http://docs.ansible.com/ansible/latest/dev_guide/developing_inventory.html#tuning-the-external-inventory-script
		*/
		hostVar := make(map[string]interface{})
		hostVar["hostvars"] = make(map[string]interface{})
		result["_meta"] = hostVar

		if item, ok := result[group].(*Items); ok {
			result[group].(*Items).Hosts = append(item.Hosts, ip)
		} else {
			it := &Items{Hosts: []string{ip}}
			result[group] = it
		}
	}

	data, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		return
	}
	fmt.Println(string(data))
}

func main() {
	api := getAPI()
	_, err := api.Version()
	if err != nil {
		return
	}

	groupName := flag.String("g", "", "Zabbix Group Name")
	flag.Parse()
	if *groupName == "" {
		api.printf("Sorry, You need use [-g] appoint group")
		return
	}
	hosts := api.GetHostByGroup(*groupName)
	api.PrintResult(hosts, *groupName)
}
