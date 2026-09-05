package main

import (
	"fmt"
	"time"
)

func main() {
	go doSomething()
	time.Sleep(time.Second)
}

func doSomething() {
	fmt.Println("hello world")
}

func Send(ch chan int, num int) {
	ch <- num
}

func Receive(ch chan int) int {
	return <- ch
}

func Send1(ch1, ch2 chan int) {
	go func() {
		ch1 <- 0
		ch1 <- 1
		ch1 <- 2
	}()
	go func() {
		ch2 <- 0
		ch2 <- 1
		ch2 <- 2
	}()

	for range 3 {
		fmt.Println(<-ch1)
		fmt.Println(<-ch2)
	}
}

func Process(nums []int) chan int {
	ch := make(chan int, 10)
	for _, val := range nums {
		ch <- val
	}

	return ch
}
