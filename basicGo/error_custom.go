package main

import "fmt"

type validationError struct {
	Message string
}

func (v *validationError) Error() string {
	return v.Message
}

func SaveData(id string, data any) error {
	if id == "" {
		return &validationError{Message: "id is empty"}
	}

	if id != "eko" {
		return &validationError{Message: "id is not eko"}
	}

	return nil
}

func main() {
	err := SaveData("eso", nil)
	if err != nil {
		if validationErr, ok := err.(*validationError); ok {
			fmt.Println(validationErr.Error())
		} else if notfoundErr, ok := err.(*validationError); ok {
			fmt.Println(notfoundErr.Error())
		} else {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("Save data success")
	}
}
