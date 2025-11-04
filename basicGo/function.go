package main

import "fmt"

func sayHello(firstname string, lastname string, age int) {
	fmt.Println("Hello", firstname, lastname, age)
}

func getHello(name string) string {
	hello := "Hello " + name
	return hello
}

func getFullname() (string, string) {
	return "Eku", "Kunedy"
}

func main() {
	sayHello("rafael", "lorenzo", 20)
	result := getHello("rafael")
	fmt.Println(result)

	firstname, _ := getFullname()
	fmt.Println(firstname)
}
