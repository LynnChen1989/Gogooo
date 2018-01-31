package main

import (
	"fmt"
)

type Bitch interface {
	Chui()
	La()
	Tan()
	Chang()
}

type CheapBitch struct {
}

type ExpensiveBitch struct {
}

type GoodWoman struct {
}

func (this CheapBitch) Chui() {
	fmt.Println("马马虎虎")
}

func (this CheapBitch) La() {
	fmt.Println("马马虎虎")

}

func (this CheapBitch) Tan() {
	fmt.Println("马马虎虎")
}

func (this CheapBitch) Chang() {
	fmt.Println("马马虎虎")
}

func (this ExpensiveBitch) Chui() {
	fmt.Println("巴适")
}

func (this ExpensiveBitch) La() {
	fmt.Println("巴适")
}

func (this ExpensiveBitch) Tan() {
	fmt.Println("巴适")
}

func (this ExpensiveBitch) Chang() {
	fmt.Println("巴适")
}

func main() {
	a := CheapBitch{}
	b := ExpensiveBitch{}
	c := GoodWoman{}

	DoService(a)
	DoService(b)
	DoService(c)
}

func DoService(girl interface{}) {
	bitch, ok := girl.(Bitch)
	if ok {
		bitch.Chui()
		bitch.La()
		bitch.Tan()
		bitch.Chang()
	} else {
		fmt.Println("你找错人了")
	}

}
