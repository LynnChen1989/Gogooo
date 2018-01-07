package main

import (
	"fmt"
	"runtime"
	"time"
)

const COUNT int = 100
const SIZE int = 300

//把0-100放在一个容量为300的数组中，然后求和 [1,2,3,4,...,100]
func main() {
	var num [SIZE]int
	for i := 0; i <= COUNT; i++ {
		num[i] = i
	}

	fmt.Println()
	fmt.Printf("result = %d\n", calmul(num[0:]))
}

func calmul(num []int) int {
	t1 := time.Now()
	var MULTICORE int = runtime.NumCPU() //计算cpu的数
	runtime.GOMAXPROCS(MULTICORE)        //在多核上运行
	fmt.Printf("with %d core\n", MULTICORE)

	ch := make(chan int) //chan是一个FIFO的队列
	for i := 0; i < MULTICORE; i++ {
		go calsome(i*COUNT/MULTICORE, (i+1)*COUNT/MULTICORE, num[0:], ch)
	}

	result := 0
	for i := 0; i < MULTICORE; i++ {
		temp := <-ch
		fmt.Printf("multicore #%d result: %d\n", i, temp)
		result += temp
	} //从channel中读取部分结果，所有结果读取完毕，循环结束

	t2 := time.Now()

	fmt.Printf("multicore total time :%d\n", t2.Sub(t1))

	return result
}

//真正开始计算的函数
func calsome(from, to int, num []int, ch chan int) {
	someresult := 0
	for i := from; i < to; i++ {
		someresult += num[i]
	}

	ch <- someresult //把结果推送到channel中去
}
