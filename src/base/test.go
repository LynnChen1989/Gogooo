package main

import (
	"fmt"
)

type Items struct {
	Hosts []string `json:"hosts"`
}

func main() {
	s := "abc"
	fmt.Println("s:", []byte(s))
}
