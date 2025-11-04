package main

import "fmt"

func endApp() {
	fmt.Println("ruskakkk")
	message := recover()
	fmt.Println(message)
}

func runApp(error bool) {
	defer endApp()

	if error {
		panic("ups error")
	}
}

func main() {
	runApp(true)
	fmt.Println("aplah")
}
