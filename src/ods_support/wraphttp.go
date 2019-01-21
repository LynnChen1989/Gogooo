package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func httpPost(url string, pb string) {
	var jsonStr = []byte(pb)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	Info.Println("response Status:", resp.Status)
	Info.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	Info.Println("response Body:", string(body))
}

func httpGet(url string) {
	//
}
