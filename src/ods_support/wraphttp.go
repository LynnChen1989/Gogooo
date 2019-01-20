package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func httpPost(url string) {
	resp, err := http.Post(url,
		"application/json",
		strings.NewReader("name=cjb"))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}
