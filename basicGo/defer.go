package main

import "fmt"

func logging() {
	fmt.Println("This is logging function")
}

func runApplication() {
	defer logging()
	fmt.Println("This is running application")
}

func main() {
	runApplication()
}
