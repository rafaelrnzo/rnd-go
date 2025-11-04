package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func ReadLine(prompt string) (string, error) {
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func ReadInt(prompt string) (int, error) {
	for {
		str, err := ReadLine(prompt)
		if err != nil {
			return 0, err
		}
		if str == "" {
			return 0, nil
		}
		v, err := strconv.Atoi(str)
		if err != nil {
			fmt.Println("⚠️  Masukkan angka yang valid.")
			continue
		}
		return v, nil
	}
}
