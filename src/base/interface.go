package main

//
///* 定义接口 */
//type interface_name interface {
//   method_name1 [return_type]
//   method_name2 [return_type]
//   method_name3 [return_type]
//   ...
//   method_namen [return_type]
//}
//
///* 定义结构体 */
//type struct_name struct {
//   /* variables */
//}
//
///* 实现接口方法 */
//func (struct_name_variable struct_name) method_name1() [return_type] {
//   /* 方法实现 */
//}
//...
//func (struct_name_variable struct_name) method_namen() [return_type] {
//   /* 方法实现*/
//}

import (
	"fmt"
)

//定义一个叫Phone的接口
type Phone interface {
	call()
}

//定义结构体
type VivoPhone struct{}
type IPhone struct{}

//方法实现
func (vivophone *VivoPhone) call() {
	fmt.Println("I am vivo")
}

func (iPhone *IPhone) call() {
	fmt.Println("I am iphone")
}

func main() {
	var phone Phone
	phone = new(VivoPhone)
	phone.call()

	phone = new(IPhone)
	phone.call()
}
