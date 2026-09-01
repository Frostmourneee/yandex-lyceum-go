package main

import "fmt"

func main() {
	var (
		name   string
		flat   int
		pwd    int
		amount float64
	)
	fmt.Scan(&name, &flat, &pwd, &amount)
	fmt.Printf("Привет, %s! Приглашаю тебя на соревнование по программированию, которое пройдёт, как всегда, в квартире %d. Оно будет длиться примерно %.1f часа. Не забудь секретный пароль для входа: %d.",
		name,
		flat,
		amount,
		pwd)
}
