package main

import "fmt"

type Address struct {
	City, Province, Country string
}

func main() {
	address := Address{"Subang", "Jawa Barat", "Indonesia"}
	address2 := &address

	address2.City = "New York"
	fmt.Println(address2)
	fmt.Println(address)

}
