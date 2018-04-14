package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func httpGet() {
	res, err := http.Get("http://www.011111y.com/11/accept.php?id=1")
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("fuck")
	fmt.Println(body)
}
