package main

import (
	"encoding/json"
	"fmt"
)

type Student struct {
	Name    string
	Age     int
	Classes []string
	Price   float32
}

func (s *Student) ShowStu() {
	fmt.Println("show Student :")
	fmt.Println("\tName\t:", s.Name)
	fmt.Println("\tAge\t:", s.Age)
	fmt.Println("\tPrice\t:", s.Price)
	fmt.Printf("\tClasses\t: ")

	for _, a := range s.Classes {
		fmt.Printf("%s ", a)
	}
	fmt.Println()
}

func main() {
	st := &Student{
		"Xiao Ming",
		27,
		[]string{"Math", "English", "Chinese"},
		9.99,
	}

	fmt.Println("before JSON encoding :")
	st.ShowStu()

	b, err := json.Marshal(st)
	if err != nil {
		fmt.Println("encoding failed.")
	} else {
		fmt.Println("encoded data: ")
		fmt.Println(b)
		fmt.Println(string(b))
	}
}
