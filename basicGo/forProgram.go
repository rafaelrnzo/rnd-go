package main

import "fmt"

func main() {
	counter := 1

	for counter < 100 {
		fmt.Println("angka ", counter)
		counter++
	}

	for angka := 1; angka <= 10; angka++ {
		fmt.Println("wok ", angka)
	}

	names := []string{"john", "james", "jowo", "jiko", "jums", "konts"}
	for i := 0; i < len(names); i++ {
		fmt.Println(names[i])
	}

	for index, name := range names {
		fmt.Println(index, name, "=", name)
	}

	for _, name := range names {
		fmt.Println(name)
	}

}
