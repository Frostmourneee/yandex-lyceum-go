package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func main() {
	res, err := TimeoutFibonacci(10, 1*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TimeoutFibonacci(n int, timeout time.Duration) (int, error) {
	if n < 0 {
		return 0, errors.New("n must be non-negative")
	}
	
	ch := make(chan int, 1)
	go func() {
		ch <- fibonacci(n)
	}()

	select {
	case result := <-ch:
		return result, nil
	case <-time.After(timeout):
		return 0, errors.New("timeout")
	}
}

func fibonacci(n int) int {
	switch n {
	case 0:
		return 0
	case 1, 2:
		return 1
	default:
		a, b := 0, 1
		for i := 1; i < n; i++ {
			a, b = b, a+b
		}
		return b
	}
}

func QuizRunner(questions, answers []string, answerCh chan string) int {
	correct := 0
	for i := 0; i < len(questions); i++ {
		select {
		case answer := <-answerCh:
			if strings.EqualFold(strings.TrimSpace(answer), strings.TrimSpace(answers[i])) {
				correct++
			}
		case <-time.After(time.Second):
			continue
		}
	}

	return correct
}
