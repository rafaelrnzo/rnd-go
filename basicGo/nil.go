package main

import "fmt"

func NewMap(name string) map[string]string {
	if name == "" {
		return nil
	} else {
		return map[string]string{
			name: name,
		}
	}
}

func main() {
	data := NewMap("")

	if data == nil {
		fmt.Println("data ksoong")
	} else {
		fmt.Println("data", data["name"])
	}
	//fmt.Println(data)
}
