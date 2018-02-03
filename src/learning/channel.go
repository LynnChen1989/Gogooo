package main

import "fmt"

//FIFO
func sum(a []int, c chan int) {
	total := 0
	for _, v := range a {
		total += v
	}

	c <- total
}

func main() {
	a := []int{5, 6, 3, -8, 3}
	c := make(chan int)
	go sum(a[:len(a)/2], c) //前面两个数求和， 然后塞进c，11
	go sum(a[len(a)/2:], c) //后面三个数求和， 然后塞进c，-2

	x, y := <-c, <-c //channel是FIFO的，所以11先进去后出来

	fmt.Printf("x=%d, y=%d, x+y=%d", x, y, x+y)

	//缓冲
	c2 := make(chan int, 1) //塞两个数进去， 长度不够， 报错
	c2 <- 1
	c2 <- 2
	fmt.Println(<-c2)
	fmt.Println(<-c2)
}
