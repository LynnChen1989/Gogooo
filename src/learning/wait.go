package main

import (
	"fmt"
	"sync"
)

var wait sync.WaitGroup

func Afunction(num int) {
	fmt.Println(num)
	wait.Done()
	// 这个函数每执行一次就把wait -1
}

func main() {
	for i := 0; i <= 10; i++ {
		wait.Add(1)
		go Afunction(i)
	}
	wait.Wait()
}
