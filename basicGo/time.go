package main

import (
	"fmt"
	"time"
)

func main() {
	currentTime := time.Now()

	fmt.Println("Current time:", currentTime)

	fmt.Println("Formatted time (RFC3339):", currentTime.Format(time.RFC3339))
	fmt.Println("Formatted time (Custom):", currentTime.Format("2006-01-02 15:04:05"))
}
