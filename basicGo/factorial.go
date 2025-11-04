package main

import "fmt"

func factorialFunction(value int) int {
	if value == 1 {
		return 1
	} else {
		return value * factorialFunction(value-1)
	}
}

func main() {
	fmt.Println(factorialFunction(20))
}
