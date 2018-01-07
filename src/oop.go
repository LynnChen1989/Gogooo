package main

import "fmt"

type Integer int

func (a Integer) Less(b Integer) bool { //面向对象
	return a < b
}

func Integer_Less(a Integer, b Integer) bool { //面向过程
	return a < b
}

func main() {
	var (
		i Integer = 1
		j Integer = 1
	)

	fmt.Println(i.Less(2))

	if Integer_Less(j, 2) {
		fmt.Println("j less 2.")
	}
}
