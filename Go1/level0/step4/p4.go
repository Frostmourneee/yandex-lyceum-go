package main

import (
	"fmt"
	"time"
)

func main() {
	var x string
	fmt.Scan(&x)

	y, _ := time.Parse("2006-01-02/15:04:05", x)

	fmt.Printf("Текущее время %d часов, %d минут. Ты точно не забыл про важные дела на сегодня?", y.Hour(), y.Minute())
}
