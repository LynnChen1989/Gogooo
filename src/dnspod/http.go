package main

import (
	"net/http"
	"log"
	"io/ioutil"
)

type Http struct {
	Url            string
	Header         map[string]string
	Username       string
	Password       string
	PostData       map[string][]string
	ResponseStatus string
	ResponseHeader http.Header
	ResponseData   interface{}
}

func (h *Http) Post() {
	// just post form
	log.Println("post request url: ", h.Url)
	req, err := http.PostForm(h.Url, h.PostData)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err.Error())
	}
	h.ResponseData = string(body)
	h.ResponseHeader = req.Header
	h.ResponseStatus = req.Status
}
