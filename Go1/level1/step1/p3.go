package main

import (
	"fmt"
)

func main() {
	var v1, v2, price int
	fmt.Scan(&v1, &v2, &price)

	if float64(v1) + float64(v2) / float64(100) >= float64(price) {
		fmt.Println("Сегодня будет вкусный кофе!")
	} else {
		fmt.Println("Стоит подкопить")
	}
}
