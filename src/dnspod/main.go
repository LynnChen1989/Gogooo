package main

import (
	"os"
	"C"
)

func WrapAddRecord(domain string, subDomain string, recordType string, recordLine string, value string, ttl string) string {
	UserId, Token := os.Getenv("UserId"), os.Getenv("Token")
	dns := DnsRecord(&DnsPod{Id: UserId, Token: Token, Domain: domain})
	id := dns.RecordCreate(subDomain, recordType, recordLine, value, ttl)
	return id
}

func WrapDelRecord(domain string, recordId string) string {
	UserId, Token := os.Getenv("UserId"), os.Getenv("Token")
	dns := DnsRecord(&DnsPod{Id: UserId, Token: Token, Domain: domain})
	id := dns.RecordRemove(recordId)
	return id
}

func WrapModifyRecord(domain string, recordId string, subDomain string, recordType string,
	recordLine string, value string, ttl string) string {

	UserId, Token := os.Getenv("UserId"), os.Getenv("Token")
	dns := DnsRecord(&DnsPod{Id: UserId, Token: Token, Domain: domain})
	id := dns.RecordModify(recordId, subDomain, recordType, recordLine, value, ttl)
	return id
}

func WrapGetRecord(domain string, subDomain string) string {
	UserId, Token := os.Getenv("UserId"), os.Getenv("Token")
	dns := DnsRecord(&DnsPod{Id: UserId, Token: Token, Domain: domain})
	id := dns.RecordList(subDomain)
	return id
}

//func main() {
//	UserId, Token := os.Getenv("UserId"), os.Getenv("Token")
//	dns := DnsRecord(&DnsPod{Id: UserId, Token: Token, Domain: "fffpl.cn"})
//	id := dns.RecordCreate("fuck.test", "A", "0", "10.10.10.18", "600")
//	fmt.Println(id)
//}
