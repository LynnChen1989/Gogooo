package main

import (
	"fmt"
)

type Integer int

func (a Integer) Less(b Integer) bool {
	return a < b
}

func (a *Integer) Add(b Integer) {
	*a += b
}

type LessAdder interface {
	Less() bool
	Add()
}

func main() {
	var a Integer = 1
	fmt.Println("a is less than b:", a.Less(2))
	fmt.Println(a.Add(8))
}
