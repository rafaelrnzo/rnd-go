package main

import "fmt"

type CopAddress struct {
	City, Province, Country string
}

func main() {
	copAddress := CopAddress{"Subang", "Jawa Barat", "Indonesia"}
	address2 := &copAddress

	address2.City = "New York"

	fmt.Println(address2)
	fmt.Println(copAddress)

	*address2 = CopAddress{"Bnadung", "Jawa Timur", "Uruguay"}
	fmt.Println(address2)
}
