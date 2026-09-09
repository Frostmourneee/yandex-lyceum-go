package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	inputs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fn := func(i int) int {
		return i * 2
	}
	workers := 5
	results, _ := ParallelMapCtx(context.Background(), inputs, fn, workers)
	// if err != nil {
	// log.Fatalf("ParallelMapCtx: %v", err)
	// }
	fmt.Println(results)
}

func ParallelMapCtx1(ctx context.Context, inputs []int, fn func(int) int, workers int) ([]int, error) {
	n := len(inputs)
	wg := &sync.WaitGroup{}
	wg.Add(2 * n)

	results := make([]int, n)
	workersChannel := make(chan bool, workers)
	jobResult := make(chan struct {
		index int
		num   int
	}, n)

	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case result := <-jobResult:
				wg.Done()
				results[result.index] = result.num
			}
		}
	}(ctxWithCancel)

	for i := 0; i < workers; i++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		go func(i int) {
			defer wg.Done()

			jobResult <- struct {
				index int
				num   int
			}{i, fn(inputs[i])}
			workersChannel <- true
		}(i)
	}

	for i := workers; i < n; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-workersChannel:
			go func(i int) {
				defer wg.Done()

				jobResult <- struct {
					index int
					num   int
				}{i, fn(inputs[i])}
				workersChannel <- true
			}(i)
		}
	}

	wg.Wait()

	return results, nil
}

func ParallelMapCtx(ctx context.Context, inputs []int, fn func(int) int, workers int) ([]int, error) {
	n := len(inputs)
	results := make([]int, n)
	ch := make(chan int)

	wg := &sync.WaitGroup{}
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for idx := range ch {
				results[idx] = fn(inputs[idx])
			}
		}()
	}

	for i := 0; i < n; i++ {
		select {
		case <-ctx.Done():
			close(ch)
			wg.Wait()
			return nil, ctx.Err()
		case ch <- i:
		}
	}
	close(ch)
	wg.Wait()

	return results, nil
}
