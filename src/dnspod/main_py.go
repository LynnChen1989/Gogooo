package main

import (
	"os"
	"C"
	"fmt"
)


//export WrapAddRecord
func WrapAddRecord(domain, subDomain, recordType, recordLine, value, ttl *C.char) *C.char {
	UserId, Token := os.Getenv("UserId"), os.Getenv("Token")
	dns := DnsRecord(&DnsPod{Id: UserId, Token: Token, Domain: C.GoString(domain)})
	id := dns.RecordCreate(C.GoString(subDomain), C.GoString(recordType), C.GoString(recordLine), C.GoString(value), C.GoString(ttl))
	fmt.Println()
	return C.CString(id)
}

//export WrapDelRecord
func WrapDelRecord(domain, recordId *C.char) string {
	UserId, Token := os.Getenv("UserId"), os.Getenv("Token")
	dns := DnsRecord(&DnsPod{Id: UserId, Token: Token, Domain: C.GoString(domain)})
	id := dns.RecordRemove(C.GoString(recordId))
	return id
}

//export WrapModifyRecord
func WrapModifyRecord(domain, recordId, subDomain, recordType, recordLine, value, ttl *C.char) string {
	UserId, Token := os.Getenv("UserId"), os.Getenv("Token")
	dns := DnsRecord(&DnsPod{Id: UserId, Token: Token, Domain: C.GoString(domain)})
	id := dns.RecordModify(C.GoString(recordId), C.GoString(subDomain), C.GoString(recordType),
		C.GoString(recordLine), C.GoString(value), C.GoString(ttl))
	return id
}

//export WrapGetRecord
func WrapGetRecord(domain, subDomain *C.char) string {
	UserId, Token := os.Getenv("UserId"), os.Getenv("Token")
	dns := DnsRecord(&DnsPod{Id: UserId, Token: Token, Domain: C.GoString(domain)})
	id := dns.RecordList(C.GoString(subDomain))
	return id
}

func main() {}
