package main

import (
	"encoding/json"
	"strconv"
	"log"
	"os"
)

type Status struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

type Record struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Enable string `json:"enable"`
	Weight string `json:"weight"`
}

type HandleReturnData struct {
	Status Status `json:"status"`
	Record Record `json:"record"`
}

type DnsRecord interface {
	RecordCreate(subDomain string, recordType string, recordLine string, value string, ttl string) string
	RecordRemove(recordId string) string
	RecordModify(recordId string, subDomain string, recordType string, recordLine string, value string, ttl string) string
	RecordList(subDomain string) string
}

type DnsPod struct {
	Id      string
	Token   string
	Domain  string
	BaseUrl string
	C       Http
}

var (
	CreatePostfix     = "Record.Create"
	RemovePostfix     = "Record.Remove"
	ModifyPostfix     = "Record.Modify"
	ListPostfix       = "Record.List"
	DomainInfoPostfix = "Domain.Info"
)

func (dp *DnsPod) SetBaseUrl(method string) {
	dp.BaseUrl = "https://dnsapi.cn/" + method
}

func (dp *DnsPod) dpPost(url string, reqData map[string][]string) {
	dp.C.Url = url
	dp.C.PostData = reqData
	dp.C.Post()
}

func (dp *DnsPod) RecordCreate(subDomain string, recordType string, recordLine string, value string, ttl string) string {
	if !(recordLine == "0" || recordLine == "3=0") {
		// 0 表示默认线路， 3=0表示海外线路
		log.Fatalln("recordLine error, just support: '0' or '3=0', get:", recordLine)
		os.Exit(1)
	}
	domainId := strconv.Itoa(dp.GetDomainId())
	dp.SetBaseUrl(CreatePostfix)
	reqData := map[string][]string{
		"login_token":    {dp.Id + "," + dp.Token},
		"format":         {"json"},
		"domain":         {dp.Domain},
		"sub_domain":     {subDomain},
		"record_type":    {recordType},
		"record_line_id": {recordLine},
		"value":          {value},
		"ttl":            {ttl},
		"domain_id":      {domainId},
	}
	dp.dpPost(dp.BaseUrl, reqData)
	var rs HandleReturnData
	json.Unmarshal([]byte(dp.C.ResponseData.(string)), &rs)
	if rs.Status.Code == "1" {
		return rs.Record.Id
	} else {
		return rs.Status.Message
	}
}

func (dp *DnsPod) RecordRemove(recordId string) string {
	domainId := strconv.Itoa(dp.GetDomainId())
	dp.SetBaseUrl(RemovePostfix)
	reqData := map[string][]string{
		"login_token": {dp.Id + "," + dp.Token},
		"format":      {"json"},
		"domain":      {dp.Domain},
		"domain_id":   {domainId},
		"record_id":   {recordId},
	}
	dp.dpPost(dp.BaseUrl, reqData)
	var rs HandleReturnData
	json.Unmarshal([]byte(dp.C.ResponseData.(string)), &rs)
	if rs.Status.Code == "1" {
		return rs.Record.Id
	} else {
		return rs.Status.Message
	}
}

func (dp *DnsPod) RecordModify(recordId string, subDomain string, recordType string,
	recordLine string, value string, ttl string) string {

	domainId := strconv.Itoa(dp.GetDomainId())
	dp.SetBaseUrl(ModifyPostfix)
	reqData := map[string][]string{
		"login_token":    {dp.Id + "," + dp.Token},
		"format":         {"json"},
		"domain_id":      {domainId},
		"record_id":      {recordId},
		"sub_domain":     {subDomain},
		"record_type":    {recordType},
		"record_line_id": {recordLine},
		"value":          {value},
		"ttl":            {ttl},
	}
	dp.dpPost(dp.BaseUrl, reqData)
	var rs HandleReturnData
	json.Unmarshal([]byte(dp.C.ResponseData.(string)), &rs)
	if rs.Status.Code == "1" {
		return rs.Record.Id
	} else {
		return rs.Status.Message
	}
}

func (dp *DnsPod) RecordList(subDomain string) string {
	// 根据子域名查询域名的详情
	dp.SetBaseUrl(ListPostfix)
	domainId := strconv.Itoa(dp.GetDomainId())
	reqData := map[string][]string{
		"login_token": {dp.Id + "," + dp.Token},
		"format":      {"json"},
		"domain_id":   {domainId},
		"sub_domain":  {subDomain},
	}
	dp.dpPost(dp.BaseUrl, reqData)
	var rs HandleReturnData
	json.Unmarshal([]byte(dp.C.ResponseData.(string)), &rs)
	if rs.Status.Code == "1" {
		return rs.Record.Id
	} else {
		return rs.Status.Message
	}
}
