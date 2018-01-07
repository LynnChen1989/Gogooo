package main

import (
	"fmt"
)

type Animal struct { //定义animal
	Name string
	Mean bool
}

type Cat struct { //定义cat，继承自animal
	Animal
	MeowStrength int //猫叫强度
}

type Dog struct { //定义dog, 继承自animal
	Animal
	BarkStrength int //狗叫强度
}

func (dog *Dog) MakeNoise() { //实现狗叫方法
	barkStrength := dog.BarkStrength
	if dog.Mean == true {
		barkStrength *= 5
	}

	for bark := 0; bark < barkStrength; bark++ {
		fmt.Printf("BARK")
	}
	fmt.Println()
}

func (cat *Cat) MakeNoise() { //实现猫叫方法
	meowStrength := cat.MeowStrength

	if cat.Mean == true {
		meowStrength *= 5
	}

	for meow := 0; meow < meowStrength; meow++ {
		fmt.Printf("MEOW")
	}

	fmt.Println()
}

func main() {
	dog1 := &Dog{ //实例化一个对象
		Animal{
			"dog1",
			true,
		},
		2,
	}

	dog1.MakeNoise()
}
