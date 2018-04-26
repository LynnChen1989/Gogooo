package main

import (
	"encoding/json"
	"strconv"
)

type DomainInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Domain struct {
	Domain DomainInfo `json:"domain"`
}

func (dp *DnsPod) GetDomainId() int {
	dp.SetBaseUrl(DomainInfoPostfix)
	reqData := map[string][]string{
		"login_token": {dp.Id + "," + dp.Token},
		"format":      {"json"},
		"domain":      {dp.Domain},
	}
	dp.C.Url = dp.BaseUrl
	dp.C.PostData = reqData
	dp.C.Post()
	var domainInfo Domain
	switch dp.C.ResponseData.(type) {
	case string:
		json.Unmarshal([]byte(dp.C.ResponseData.(string)), &domainInfo)
		domainId, _ := strconv.Atoi(domainInfo.Domain.Id)
		return domainId
	}
	return 999999
}
