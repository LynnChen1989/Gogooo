package main

import (
	"net/http"
	"fmt"
	"strings"
	"log"
	"os"
	"encoding/json"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var ServerHost = getEnv("ServerHost", "0.0.0.0")
var ServerPort = getEnv("ServerPort", "9090")

func GenerateAvatarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()
	var Value string
	var keyCnt int
	for k, v := range r.Form {
		if k == "key" {
			keyCnt += 1
			Value = strings.Join(v, "")
		}
	}
	log.Println("generate avatar key word is:", Value)
	i := Generate(Value)
	Render(i, Value)

	var data Response
	if keyCnt < 1 {
		data = Response{Code: 400, Message: "bad request, you must provide args named `key`"}
	} else {
		data = Response{Code: 200, Message: fmt.Sprintf("http://%s:%s/avatar/%s.png", ServerHost, ServerPort, Value)}
	}
	jData, err := json.Marshal(data)
	if err != nil {
		log.Println("get json data error:", err)
	}
	w.Write(jData)
}

func GetAvatarHandler(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Path
	png := strings.Split(vars, "/")
	http.ServeFile(w, r, png[2])
}

func main() {
	http.HandleFunc("/", GenerateAvatarHandler)   // route
	http.HandleFunc("/avatar/", GetAvatarHandler) // route
	log.Printf("server starting on http://%s:%s", ServerHost, ServerPort)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", ServerHost, ServerPort), nil) // listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
