package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	TestSequentionalHttpVsGoroutines() 
}

func Contains(ctx context.Context, r io.Reader, seq []byte) (bool, error) {
	window := make([]byte, 2*len(seq))
	initN, initErr := r.Read(window)
	if initN > 0 && bytes.Contains(window[:initN], seq) {
		return true, nil
	}
	if initErr == io.EOF {
		return false, nil
	}
	if initErr != nil {
		return false, initErr
	}

	tmp := make([]byte, len(seq))
	for {
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		n, err := r.Read(tmp)
		if n > 0 {
			window = append(window[len(window)/2:], tmp[:n]...)
			if bytes.Contains(window, seq) {
				return true, nil
			}
		}
		if err == io.EOF {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
}

func TestContext() {
	// ctx := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(2)
	var flag bool
	go func() {
		defer wg.Done()
		if err := processSourceData(&flag); err != nil {
			fmt.Printf("processSourceData(ctx): %s", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := readSource(); err != nil {
			flag = true
			fmt.Printf("readSource(ctx): %s", err)
		}
	}()

	wg.Wait()
}

func readSource() error {
	time.Sleep(3 * time.Second)
	return fmt.Errorf("some error in readSource")
}

func processSourceData(flag *bool) error {
	for {
		<-time.After(time.Second)
		fmt.Println("process data bit by bit...")

		if *flag {
			fmt.Println("processSourceData was canceled")
			return nil
		}
	}
}

func f1() {
	req, err := http.NewRequest(http.MethodGet, "https://ya.ru", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

func Server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	http.ListenAndServe(":8080", nil)
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8081/provideData", nil)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(body)
}

func StartServer(maxTimeout time.Duration) {
	http.Handle("/readSource", http.TimeoutHandler(http.HandlerFunc(myHandler), maxTimeout, "Timeout!"))

	http.ListenAndServe(":8080", nil)
}	

type APIResponseSmall struct {
	Data       string
	StatusCode int
}

var client = &http.Client{}

func fetchAPI(ctx context.Context, url string, timeout time.Duration) (*APIResponseSmall, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &APIResponseSmall{
		Data:       string(body),
		StatusCode: resp.StatusCode,
	}, nil
}

type APIResponse struct {
	URL        string
	Data       string
	StatusCode int
	Err        error
}

func FetchAPI(ctx context.Context, urls []string, timeout time.Duration) []*APIResponse {
	n := len(urls)
	responses := make([]*APIResponse, n)
	wg := sync.WaitGroup{}
	wg.Add(n)

	for idx, url := range urls {
		go func(idx int, url string) {
			defer wg.Done()

			r, err := fetchAPI(ctx, url, timeout)
			if err != nil {
				resp := &APIResponse{URL: url, Err: err}
				responses[idx] = resp
				return
			}

			responses[idx] = &APIResponse{
				URL:        url,
				Data:       r.Data,
				StatusCode: r.StatusCode,
				Err:        err,
			}
		}(idx, url)
	}

	wg.Wait()
	return responses
}

func TestSequentionalHttpVsGoroutines() {
	var urls [100]string
	for i := range urls {
		urls[i] = "https://httpbin.org/get"
	}

	// start := time.Now()
	// for _, url := range urls {
	// 	fetchAPI(context.Background(), url, 1*time.Second)
	// }
	// elapsed := time.Since(start)
	// fmt.Printf("Sequential HTTP: %v\n", elapsed) 
	// 1m36.7 on https://httpbin.org/get url

	start := time.Now()
	FetchAPI(context.Background(), urls[:], 1*time.Second)
	elapsed := time.Since(start)
	fmt.Printf("Concurrent HTTP: %v\n", elapsed) 
	// 1.002 sec on https://httpbin.org/get url
}
