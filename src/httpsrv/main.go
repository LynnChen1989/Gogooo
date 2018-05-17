package main


import (
	"fmt"
	"net/http"
	"strings"
	"log"
)

func GeneratePassHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") // write to client
}

func main() {
	http.HandleFunc("/", GeneratePassHandler) // route
	fmt.Println("server starting ...")
	err := http.ListenAndServe(":9090", nil) // listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
