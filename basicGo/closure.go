package main

import "fmt"

func main() {
	counter := 0
	increment := func() func(int) int {
		return func(step int) int {
			counter += step
			return counter
		}
	}

	counts := increment()

	fmt.Println(counts(2))
	fmt.Println(counts(3))
	fmt.Println(counts(5))
}
