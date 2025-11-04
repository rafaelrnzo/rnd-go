package main

import "fmt"

type Person struct {
	Name string
}

func (person Person) GetName() string {
	return person.Name
}

type HasName interface {
	GetName() string
}

func SayHello(value HasName) {
	fmt.Println("hello", value.GetName())
}

func Ups() any {
	return 42
}

func main() {
	var kosong any = Ups()
	fmt.Println(kosong)

	person := Person{Name: "John"}
	SayHello(person)
}
