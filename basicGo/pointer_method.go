package main

import "fmt"

type People struct {
	Name   string
	Gender string
}

func (m *People) Married() {
	switch m.Gender {
	case "Male":
		m.Name = "Mr. " + m.Name
	case "Female":
		m.Name = "Ms. " + m.Name

	}
}

func main() {
	eko := People{"Eko", "Male"}
	eka := People{"Eka", "Female"}
	eko.Married()
	eka.Married()

	fmt.Println(eko.Name)
	fmt.Println(eka.Name)
}
