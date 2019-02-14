package main

import (
	"fmt"
	"regexp"
)

func main() {
	b := []byte("root:fhl3mjsdwj@tcp(172.16.1.18:50001)/db_cms01")
	pat := `\(.*\)`
	reg1 := regexp.MustCompile(pat)      // 第一匹配
	reg2 := regexp.MustCompilePOSIX(pat) // 最长匹配
	fmt.Printf("%s\n", reg1.Find(b))     // abc1
	fmt.Printf("%s\n", reg2.Find(b))     // abc1def1

}
