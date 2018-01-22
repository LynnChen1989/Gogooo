package main

import (
	"fmt"
)

// 7-51行代码是介绍iota, const, slice

//// iota
//const (
//    c0 = iota * 42    // iota = 0
//    c1 float64 = iota * 42 //iota =1
//    c2 = iota * 42 //iota=2
//)
//
//const (
//    Sunday = iota
//    Monday
//    Tuesday
//    Wednesday
//    Thursday  //=4
//    Friday
//    Saturday
//    numberOfDays // 这个常量没有导出
//)
//
//
//
//func main()  {
//    fmt.Println(c0, c1, c2)
//    fmt.Println(Thursday)
//
//    var arr [10]int = [10]int {1,2,3,4,5,6,7,8,9,10}
//    var myslice []int = arr[:5]
//
//    fmt.Println("arr is :")
//    for _, v := range arr{
//	fmt.Print(v, "")
//    }
//    fmt.Println("slice is: ")
//    for _, v := range myslice{
//	fmt.Print(v, "")
//    }
//    // use make
//    myslice1 := make([]int, 5, 8)
//    //myslice2 := make([]int, 5, 10)
//    myslice3 := []int{1,2,3,4,5,6}
//
//    myslice1 = append(myslice1, myslice3...)
//    fmt.Println(myslice1)
//}

//// map用法
//type PersonInfo struct {
//    ID string
//    Name string
//    Address string
//}
//
//func main()  {
//    var personDB map[string] PersonInfo
//    personDB = make(map[string] PersonInfo)
//
//    personDB["12345"] = PersonInfo{"12345", "Tom", "China"}
//    personDB["110"] = PersonInfo{"110", "Jack", "USA"}
//
//    person, ok := personDB["1234"]
//
//    if ok{
//	fmt.Println("Founf person", person.Name, "with ID 1234.")
//    }else {
//	fmt.Println("Not Found.")
//    }
//}

// 使用interface传递任意类型的参数
func PrintContent(args ...interface{}) {
	for _, arg := range args {
		switch arg.(type) {
		case int:
			fmt.Printf("%d is int.\n", arg)
		case string:
			fmt.Printf("%s is string.\n", arg)
		default:
			fmt.Printf("I dont kown what fuck type.\n")
		}
	}
}

func main() {
	var i int = 1
	var j string = "chenlin"
	PrintContent(i)
	PrintContent(j)
}
