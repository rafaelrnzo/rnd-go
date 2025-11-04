package main

import "fmt"

func main() {
	person := map[string]string{
		"name":    "Jokowi Boti",
		"address": "Ngawi",
		"test":    "true",
	}

	fmt.Println(person)
	fmt.Println(person["name"])
	fmt.Println(person["address"])

	delete(person, "address")

	fmt.Println(person)
}
