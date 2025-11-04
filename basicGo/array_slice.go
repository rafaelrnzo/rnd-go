package main

import "fmt"

func main() {
	days := [...]string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}
	fmt.Println(days)

	daysSlice := days[5:]
	fmt.Println(daysSlice)

	daysSlice[0] = "Manday"
	daysSlice[1] = "Mday"
	fmt.Println(days)

	daysSlice2 := append(daysSlice, days[:]...)
	fmt.Println(daysSlice2)

	newSlice := make([]string, 2, 5)
	newSlice[0] = "Monday"
	newSlice[1] = "Tuesday"

	fmt.Println(newSlice)
	fmt.Println(len(newSlice))
	fmt.Println(cap(newSlice))

	fromSlice := days[:]
	toSlice := make([]string, len(fromSlice), cap(fromSlice))

	copy(toSlice, fromSlice)
	fmt.Println(toSlice)

	iniArray := [...]int{1, 2, 3}
	iniSlice := []int{1, 2, 3}

	fmt.Println(iniArray)
	fmt.Println(iniSlice)
}
