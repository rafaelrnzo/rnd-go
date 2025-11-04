package main

import "fmt"

func getGoogBye(name string) string {
	return "Good bye " + name
}

func main() {
	contoh := getGoogBye
	fmt.Println(contoh("Goog"))
}
