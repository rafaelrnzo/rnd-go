package main

import "fmt"

func main() {
	names := [...]string{"ahmad", "raya", "blongo", "ikhsan", "bram"}

	slice1 := names[3:5]
	slice2 := names[:3]
	slice3 := names[4:]
	slice4 := names[:]

	fmt.Println(slice1)
	fmt.Println(slice2)
	fmt.Println(slice3)
	fmt.Println(slice4)
}
