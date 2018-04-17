package main

import (
	_ "fmt"
	_ "net/http"
	_ "net/url"
	"fmt"
)

var channel chan int = make(chan int)

func saySomething(str string) {
	for i := 0; i <= 10; i ++ {
		fmt.Println(str)
	}
	channel <- 0
}

func main() {
	go saySomething("fuck")
}
