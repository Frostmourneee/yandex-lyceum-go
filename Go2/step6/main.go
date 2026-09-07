package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	var once sync.Once
	wg := &sync.WaitGroup{}
	initializeResources := func() {
		time.Sleep(time.Second)
		fmt.Println("Only once initialize something")
	}
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(initializeResources)
		}()
	}
	wg.Wait()
}

func listen(name string, data map[string]string, c *sync.Cond) {
	c.L.Lock()
	c.Wait()

	fmt.Printf("[%s] %s\n", name, data["key"])

	c.L.Unlock()
}

func getMark(name string) (string, error) {
	url := "http://localhost:8082/mark?name=" + url.QueryEscape(name)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: status %d", name, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func Compare(name1, name2 string) (string, error) {
	wg := &sync.WaitGroup{}

	var mark1, mark2 string
	var err1, err2 error
	wg.Add(2)
	go func() { mark1, err1 = getMark(name1); wg.Done() }()
	go func() { mark2, err2 = getMark(name2); wg.Done() }()
	wg.Wait()

	if err1 != nil {
		return "", err1
	}
	if err2 != nil {
		return "", err2
	}

	switch {
	case mark1 > mark2:
		return ">", nil
	case mark1 < mark2:
		return "<", nil
	default:
		return "=", nil
	}
}

func Average(names []string) (int, error) {
	type response struct {
		mark string
		err  error
	}

	n := len(names)
	wg := &sync.WaitGroup{}
	wg.Add(n)

	responses := make([]response, n)
	for i, name := range names {
		go func(idx int, s string) {
			mark, err := getMark(s)
			responses[idx] = response{mark: mark, err: err}
			wg.Done()
		}(i, name)
	}
	wg.Wait()

	total := 0
	for _, response := range responses {
		if response.err != nil {
			return 0, response.err
		}

		num, err := strconv.Atoi(response.mark)
		if err != nil {
			return 0, err
		}

		total += num
	}

	return total / n, nil
}

func BestStudents(names []string) (string, error) {
	type response struct {
		name string
		mark string
		err  error
	}

	n := len(names)
	wg := &sync.WaitGroup{}
	wg.Add(n)

	responses := make([]response, n)
	for i, name := range names {
		go func(idx int, s string) {
			mark, err := getMark(s)
			responses[idx] = response{name: s, mark: mark, err: err}
			wg.Done()
		}(i, name)
	}
	wg.Wait()

	total := 0
	for _, response := range responses {
		if response.err != nil {
			return "", response.err
		}

		num, err := strconv.Atoi(response.mark)
		if err != nil {
			return "", err
		}

		total += num
	}

	best := slices.DeleteFunc(responses, func(r response) bool {
		mark, _ := strconv.Atoi(r.mark)
		return mark*n <= total
	})

	slices.SortFunc(best, func(a, b response) int {
		return strings.Compare(a.name, b.name)
	})

	resnames := make([]string, len(best))
	for i, r := range best {
		resnames[i] = r.name
	}
	return strings.Join(resnames, ","), nil
}

func CompareList(names []string) (map[string]string, error) {
	type response struct {
		name string
		mark string
		err  error
	}

	n := len(names)
	wg := &sync.WaitGroup{}
	wg.Add(n)

	responses := make([]response, n)
	for i, name := range names {
		go func(idx int, s string) {
			mark, err := getMark(s)
			responses[idx] = response{name: s, mark: mark, err: err}
			wg.Done()
		}(i, name)
	}
	wg.Wait()

	total := 0
	for _, response := range responses {
		if response.err != nil {
			return nil, response.err
		}

		num, err := strconv.Atoi(response.mark)
		if err != nil {
			return nil, err
		}

		total += num
	}

	resnames := make(map[string]string, n)
	for _, r := range responses {
		mark, _ := strconv.Atoi(r.mark)
		switch {
		case mark * n > total:
			resnames[r.name] = ">"
		case mark * n < total:
			resnames[r.name] = "<"
		default:
			resnames[r.name] = "="
		}
	}
	return resnames, nil
}
