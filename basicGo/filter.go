package main

import (
	"fmt"
	"strings"
)

type Filter func(string) string

func sayHelloWithFilter(name string, filter Filter) {
	fmt.Println(filter(name))
}

func spamFilter(name string) string {
	banned := []string{"anjing", "bangsat", "tolol", "Jokowo", "ngtd", "anj"}
	for _, w := range banned {
		if strings.EqualFold(name, w) {
			return "***"
		}
	}
	return name
}

func main() {
	sayHelloWithFilter("Rafael", spamFilter)
	sayHelloWithFilter("bangsat", spamFilter)
}
