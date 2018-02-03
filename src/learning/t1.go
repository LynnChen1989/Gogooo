package main

import (
	"fmt"
	_ "fmt"
	"net/http"
	_ "net/http"
	_ "net/url"
)

var client = &http.Client{}

func getMethod() {
	uri := "http://studygolang.com/"
	request, _ := http.NewRequest("GET", uri, nil)
	response, _ := client.Do(request)

	defer response.Body.Close()

	if response.StatusCode == 200 {
		fmt.Println("get request OK")
	}
}

func main() {
	getMethod()
}
