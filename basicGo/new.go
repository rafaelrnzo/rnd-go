package main

import "fmt"

type NewAddress struct {
	city, province, country string
}

func main() {
	var alamat1 *NewAddress = new(NewAddress)
	var alamat2 *NewAddress = alamat1

	alamat2.country = "Indonesia"

	fmt.Println(alamat2)
	fmt.Println(alamat1)

}
