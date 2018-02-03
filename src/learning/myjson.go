package main

import (
	"fmt"
)

func Greeting(prefix string, who ...string) {
	fmt.Println(prefix)
	for _, name := range who {
		fmt.Println(name)
	}
}

func main() {
	Greeting("hello", "chenlin", "hehhe")
}
