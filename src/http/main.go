package main

import "fmt"

func main() {
	header := map[string]string{
		"Authorization": "Token 9c2dec0536aad82eb95f5b7dd3e640c218fae2b0",
	}
	url := "http://rms-backend.ops.dragonest.com/api/v1cdn/report/?month=201803"
	var data = httpGet(url, header)
	fmt.Println(data)
}
