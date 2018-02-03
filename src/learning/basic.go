package main

import (
	"fmt"
	//"math"
)

//func main() {
//	var x float64 //静态类型声明
//	x = 20.0
//	y := 42 //动态类型声明
//	const LENGTH int  = 10
//	const WIDTH int  = 5
//	var area int
//	area = LENGTH * WIDTH
//	fmt.Println(area)
//	fmt.Println(x,y)
//	fmt.Printf("x is type of %T \n", x)
//	fmt.Printf("y is type of %T \n", y)
//}

//func main ()  {
//	var a int = 100;
//	if(a < 20){
//		fmt.Println("a is less than 20.")
//	}else{
//		fmt.Println("a is not less than 20.")
//	}
//}

//func main(){
//	var b int = 5
//	var a int
//	numbers := [4]int{1,2,3,5}
//
//	for a := 0; a < 5; a++ {
//		fmt.Printf("TEST A ==> value of a :%d\n", a)
//	}
//
//	for a < b {
//		a++
//		fmt.Printf("TEST B ==> value of a : %d\n", a)
//	}
//
//	for i, x:=range numbers{
//		fmt.Printf("TEST C ==> value of x = %d at %d\n", x, i)
//	}
//}

//func main(){
//    var i, j int
//    for i=0; i<100; i++ {
//	for j=2; j<=(i/j); j++{
//	    if(i%j == 0){
//		break;
//	    }
//	}
//	if(j > i/j){
//	    fmt.Printf("%d ids prime\n",i)
//	}
//    }
//}

//func main(){
//    //getSquareRoot := func(x float64)  float64 { //函数作为值
//	//return math.Sqrt(x)
//    //}
//    //fmt.Println(getSquareRoot(9))
//}

//func getSequence() func() int{
//    i:=0
//    return func() int {
//	i += 1
//	return  i
//    }
//}
//
//func main() {
//    nextNumber := getSequence()
//    fmt.Println(nextNumber())
//    fmt.Println(nextNumber())
//    fmt.Println(nextNumber())
//}
//
//type Circle struct {
//    x,y,radius float64
//}
//
//func (circle Circle) area() float64 {
//    return math.Pi * circle.radius * circle.radius
//}
//
//func main()  {
//    circle := Circle{x:0, y:0, radius:5}
//    fmt.Printf("Circle area: %f", circle.area())
//}
//

//func getAverage(arr []int, size int) float32{
//    var i int
//    var sum int
//    var avg float32
//
//    for i=0; i<size; i++{
//	sum += arr[i]
//    }
//    avg = float32(sum / size)
//    return avg
//}
//
//func main()  {
//    var balance = []int{2,3,4,5}
//    var avg float32
//    avg = getAverage(balance, 4)
//    fmt.Printf("Average value is %f", avg)
//}

//func main()  {
//    var a int = 10
//    var ip *int
//    ip = &a
//    fmt.Printf("Address of `a` variable :%x\n", &a)
//    fmt.Printf("Address of `ip` variable :%x\n", ip)
//    fmt.Printf("Value of *ip variable: %d\n", *ip)
//}

//const MAX int = 3
//
//func main()  {
//    a := []int{10,200,100}
//    var i int
//    var ptr [MAX]*int
//
//    for i=0; i< MAX; i++{
//    	ptr[i] = &a[i]
//    }
//
//    for i=0; i< MAX; i++{
//    	fmt.Printf("a[%d] = %d\n", i, *ptr[i])
//    }
//}

//func main()  {
//    var a int
//    var ptr *int
//    var pptr **int
//
//    a = 3000
//    ptr = &a
//    pptr = &ptr
//
//    fmt.Printf("变量a = %d\n", a)
//    fmt.Printf("指针变量*ptr = %d\n", *ptr)
//    fmt.Printf("指向指针的指针变量*pptr = %d\n", **pptr)
//
//}

//func main() {
//    /* 定义局部变量 */
//    var a int = 100
//    var b int = 200
//    var c int = 1
//    var d int = 2
//
//    fmt.Printf("交换前 a 的值 : %d\n", a )ma
//    fmt.Printf("交换前 b 的值 : %d\n", b )
//
//    /* 调用函数用于交换值
//    * &a 指向 a 变量的地址
//    * &b 指向 b 变量的地址
//    */
//    swap(&a, &b);
//
//    fmt.Printf("交换后 a 的值 : %d\n", a )
//    fmt.Printf("交换后 b 的值 : %d\n", b )
//
//    c ,d = d, c
//    fmt.Printf("交换后 c 的值 : %d\n", c )
//    fmt.Printf("交换后 d 的值 : %d\n", d )
//}
//
//func swap(x *int, y *int) {
//    var temp int
//    temp = *x    /* 保存 x 地址的值 */
//    *x = *y      /* 将 y 赋值给 x */
//    *y = temp    /* 将 temp 赋值给 y */
//}

//// 结构体
//type Person struct  {
//    name string
//    age int
//    sex string
//}
//
//func printPerson( p *Person ){
//    fmt.Printf("Person 1 name is : %s\n", p.name)
//    fmt.Printf("Person 1 age is : %d\n", p.age)
//    fmt.Printf("Person 1 sex is : %s\n", p.sex)
//}
//
//func main()  {
//    var p1 Person
//    p1.name = "chenlin"
//    p1.age = 27
//    p1.sex = "male"
//
//    printPerson(&p1)
//}
//

//切片

//func main()  {
//    var age = []int{23,45}
//    s := []int{1,2,3}
//    var num = make([]int, 2, 4)
//
//    fmt.Println(s[0:1])
//    fmt.Println(age[0:1])
//    fmt.Printf("len=%d cap=%d slice=%v\n",len(num),cap(num),num) // len() 方法获取长度,计算容量的方法 cap() 可以测量切片最长可以达到多少
//}

//range

func main() {
	nums := []int{1, 2, 3}
	sum := 0

	for _, num := range nums {
		sum += num
	}

	fmt.Println("sum:", sum)

	for i, num := range nums {
		if num == 3 {
			fmt.Println("index:", i)
		}
	}

	// range 可以再map的键值对上
	kvs := map[string]string{"a": "apple", "b": "banana"}
	for k, v := range kvs {
		fmt.Printf("%s -- > %s\n", k, v)
	}

	for i, c := range "go" {
		fmt.Println(i, c)
	}
}
