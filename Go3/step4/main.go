package main

import (
	"fmt"
	"net/http"
)

func main() {
	ex3()
}

type apiHandler struct{}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func ex1() {
	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler{})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	http.ListenAndServe(":8080", mux)
}

func ex2() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("name")
		msg := specifyClient(name)
		fmt.Fprint(w, msg)
	})

	http.ListenAndServe(":8080", mux)
}

func specifyClient(name string) string {
	if name == "" {
		return "hello stranger"
	}

	if checkOnlyEnglishLetters(name) {
		return "hello " + name
	}

	return "hello dirty hacker"
}

func checkOnlyEnglishLetters(name string) bool {
	for _, c := range name {
		if c < 'a' || ('z' < c && c < 'A') || c > 'Z' {
			return false
		}
	}

	return true
}

func ex3() {
	counter := counter()
	fib := fibCounter()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%d", fib())
	})

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%d", counter())
	})

	http.ListenAndServe(":8080", mux)
}

func fibCounter() func() int {
	a, b := 0, 1
	return func() int {
		res := a
		a, b = b, a+b
		return res
	}
}

func counter() func() int {
	c := 0

	return func() int {
		c++
		return c
	}
}
