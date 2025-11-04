package main

import "fmt"

func getCompleteName() (firstName, lastName string) {
	firstName, lastName = "rafael", "lorenzo"
	return firstName, lastName
}
func main() {
	a, b := getCompleteName()
	fmt.Println(a, b)
}
