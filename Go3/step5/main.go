package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sync"
)

func main() {
	t3()
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func t1() {
	mux := http.NewServeMux()
	mux.Handle("/hello", Logger(http.HandlerFunc(helloHandler)))

	http.ListenAndServe(":8080", mux)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := slog.String("path", r.URL.Path)
		method := slog.String("method", r.Method)
		slog.Info("incoming request", method, path)

		next.ServeHTTP(w, r)
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, middleware!")
}

func t2() {
	mux := http.NewServeMux()
	mux.Handle("/hello", SetDefaultName(Sanitize(HelloHandler)))

	http.ListenAndServe(":8080", mux)
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello "+r.URL.Query().Get("name"))
}

func SetDefaultName(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("name") == "" {
			fmt.Fprint(w, "hello stranger")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Sanitize(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if anyNonEnglish(r.URL.Query().Get("name")) {
			fmt.Fprint(w, "hello dirty hacker")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func anyNonEnglish(s string) bool {
	for _, c := range s {
		if c < 'A' || ('Z' < c && c < 'a') || c > 'z' {
			return true
		}
	}

	return false
}

type SafeCounter struct {
	c  int
	mu sync.RWMutex
}

var fib = fibClosure()

var sc SafeCounter = SafeCounter{}

func (sc *SafeCounter) Inc() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.c++
}

func (sc *SafeCounter) Get() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.c
}

func t3() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Metrics(FiboHandler))

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", sc.Get())
	})

	http.ListenAndServe(":8080", mux)
}

func FiboHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d", fib())
}

func fibClosure() func() int {
	a, b := 0, 1

	return func() int {
		res := a
		a, b = b, a+b
		return res
	}
}

func Metrics(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sc.Inc()
		next.ServeHTTP(w, r)
	})
}

func t4() {
	mux := http.NewServeMux()
	mux.HandleFunc("/answer/", Authorization(answerHandler))

	http.ListenAndServe(":8080", mux)
}

func answerHandler(w http.ResponseWriter, r *http.Request) {
	username, _, _ := r.BasicAuth()
	fmt.Fprintf(w, "Welcome, %s!", username)
}

func Authorization(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || username == "" || password == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
