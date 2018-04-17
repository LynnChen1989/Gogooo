package main

import (
	"encoding/json"
	"log"
)

type AppDetail struct {
	ProjectName string `json:"display_name"`
	ProjectCode string `json:"code"`
	Stage       string `json:"stage"`
}

// 接收接口返回的json
type Hosts struct {
	PublicIp   []string  `json:"public_ip"`
	Hostname   string    `json:"hostname"`
	ServerType string    `json:"type"`
	App        AppDetail `json:"app"`
}

type SavedHosts struct {
	PublicIp    string    `json:"public_ip"`
	Hostname    string    `json:"hostname"`
	Type        string    `json:"type"`
	ProjectInfo AppDetail `json:"project_info"`
}

type CheckExist struct {
	Hostname string `json:"hostname"`
}

func getHost() (hosts []Hosts) {
	header := map[string]string{
		"Authorization": "Token 01B7CF014E8FF0164EFF6BFDC69CA12B0283ABBB",
	}
	url := "https://cmdb-api.ops.dragonest.com/api/hosts/"
	var data = httpGet(url, header)
	json.Unmarshal([]byte(data), &hosts)
	return
}

func saveHost(hosts []Hosts) {
	for _, v := range hosts {
		var ip string
		if len(v.PublicIp) > 0 {
			ip = v.PublicIp[0]
		} else {
			ip = "0.0.0.0"
		}
		data := &SavedHosts{
			PublicIp:    ip,
			Hostname:    v.Hostname,
			Type:        v.ServerType,
			ProjectInfo: v.App,
		}
		condition := map[string]interface{}{
			"publicip": ip,
		}
		check := &CheckExist{}
		getOneByCondition("db_go", "tb_cmdb_hosts", condition, check)
		if len(check.Hostname) > 0 {
			log.Printf("%s already existed, ignore.", v.Hostname)
		} else {
			log.Printf("%s does not exist, save it.", v.Hostname)
			Insert("db_go", "tb_cmdb_hosts", data)
		}
	}
}
