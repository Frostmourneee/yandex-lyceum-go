package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	ex1()
}

func ex1() {
	chan_c := make(chan string)
	go func() {
		time.Sleep(1 * time.Second)
		chan_c <- "Инструкции выполнены успешно."
	}()

	timeout := time.AfterFunc(2*time.Second, func() {
		chan_c <- "Время выполнения истекло."
	})
	fmt.Println("Waiting for result...")

	result := <-chan_c
	timeout.Stop()
	fmt.Println(result)
	// timeout.Stop()
}

func GeneratePrimeNumbers(stop chan struct{}, prime_nums chan int, N int) {
	time.AfterFunc(time.Second/10, func() {
		close(stop)
	})

	for i := 2; i <= N; i++ {
		if !isPrime(i) {
			continue
		}
		select {
		case <-stop:
			close(prime_nums)
			return
		case prime_nums <- i:
		}
	}
	close(prime_nums)
}

func isPrime(n int) bool {
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func readJSON(ctx context.Context, path string, result chan<- []byte) {
	defer close(result)

	if err := ctx.Err(); err != nil {
		return
	}

	done := make(chan struct {
		data []byte
		err  error
	}, 1)
	go func() {
		data, err := os.ReadFile(path)
		done <- struct {
			data []byte
			err  error
		}{data, err}
	}()

	select {
	case <-ctx.Done():
		return
	case out := <-done:
		if out.err != nil {
			return
		}
		// select {
		// case <-ctx.Done():
		// 	return
		// case result <- out.data:
		// }
		result <- out.data
	}
}
