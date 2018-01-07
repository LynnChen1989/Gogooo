package main

import (
	"errors"
	"fmt"
)

type item struct {
	Name string
}

func (i item) String() string {
	return fmt.Sprintf("item name: %v", i.Name)
}

type person struct {
	Name string
	Sex  string
}

func (p person) String() string {
	return fmt.Sprintf("person name: %v sex: %v", p.Name, p.Sex)
}

func Parse(i interface{}) interface{} {
	switch i.(type) {
	case string:
		return &item{
			Name: i.(string),
		}
	case []string:
		data := i.([]string)
		length := len(data)
		if length == 2 {
			return &person{
				Name: data[0],
				Sex:  data[1],
			}
		} else {
			return nil
		}
	default:
		panic(errors.New("type match miss"))
	}
	return nil
}

func main() {
	p1 := Parse("Apple").(*item)
	fmt.Println(p1)
	p2 := Parse([]string{"zhangsan", "man"}).(*person)
	fmt.Println(p2)
}
