package main

import (
	"fmt"
)

func main() {
	var v1 int
	fmt.Scan(&v1)

	var msg string
	if v1 == 0 {
		msg = "Число 0"
	} else if -10 < v1 && v1 < 10 {
		msg = "Число однозначное"
	} else if v1 % 2 == 0 {
		msg = "Число чётное"
	} else if v1 > 0 {
		msg = "Число положительное"
	} else {
		msg = "Число красивое"
	}

	fmt.Println(msg)
}
