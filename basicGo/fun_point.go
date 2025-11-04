package main

import "fmt"

type Employee struct {
	FirstName, LastName string
}

func ChangeNameToBudi(employee *Employee) {
	employee.LastName = "Lowe"
}

func main() {
	employee := Employee{"Rafael", "Lorenzo"}
	ChangeNameToBudi(&employee)

	fmt.Println(employee)
}
