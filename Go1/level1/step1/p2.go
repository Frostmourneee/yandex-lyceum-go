package main

import (
	"fmt"
)

func main() {
	var v1, v2 int
	fmt.Scan(&v1, &v2)

	if v1 == v2 {
		fmt.Println("Числа равны")
		return
	}

	if v1 > v2 {
		fmt.Println("Первое число больше второго")
	} else {
		fmt.Println("Второе число больше первого")
	}
}
