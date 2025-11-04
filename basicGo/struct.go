package main

import "fmt"

type Customer struct {
	Name, Address string
	Age           int
}

func (customer Customer) sayHello(name string) {
	fmt.Println("Hello", name, "my name is", customer.Name)
}

func main() {
	var eko Customer
	eko.Name = "Eko"
	eko.Age = 18
	eko.Address = "Indonesia"

	jowoki := Customer{
		Name:    "Jowoki",
		Age:     18,
		Address: "Indonesia",
	}

	wokwok := Customer{"Wokwok", "Manokwari", 18}

	fmt.Println(eko)
	fmt.Println(jowoki)
	fmt.Println(wokwok)

	wokwok.sayHello("joko")
}
