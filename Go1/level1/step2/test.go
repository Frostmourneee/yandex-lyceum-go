package main

import (
	"fmt"
)

type Animal struct {
	age int
}

type Cat struct {
	Animal
}

type Dog struct {
	Animal
}

type Naming interface {
	Name() string
}

func (d Dog) Name() string {
	if d.age < 12 {
		return "Щенок"
	}

	return "Собака"
}

func (c Cat) Name() string {
	if c.age < 12 {
		return "Котёнок"
	}

	return "Кошка"
}

func main() {
	var petType string
	var age int

	fmt.Scan(&petType, &age)

	var pet Naming
	switch petType {
	case "cat":
		pet = Cat{age: age}
	case "dog":
		pet = Dog{age: age}
	}

	fmt.Println(pet.Name())
}
