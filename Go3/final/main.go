package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mux := http.NewServeMux()
	mux.Handle("GET /users/{id}", getUserOuter(getUserHandler))
	mux.HandleFunc("POST /users", createUserHandler)

	handler := loggingMiddleware(logger, mux)
	http.ListenAndServe(":8080", handler)
}

var store = NewStore()

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	if user, ok := store.GetUser(id); !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
	} else {
		jsonUser, _ := json.Marshal(user)
		w.Write(jsonUser)
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		Age  int
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	} else {
		user := store.CreateUser(input.Name, input.Age)
		jsonUser, _ := json.Marshal(user)
		w.Write(jsonUser)
	}
}

func getUserOuter(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := strconv.Atoi(r.PathValue("id")); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Store struct {
	users map[int]User
	mu    sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		users: make(map[int]User),
		mu:    sync.RWMutex{},
	}
}

func (s *Store) CreateUser(name string, age int) User {
	s.mu.Lock()
	defer s.mu.Unlock()

	newID := len(s.users) + 1
	newUser := User{
		ID:   newID,
		Name: name,
		Age:  age,
	}
	s.users[newID] = newUser

	return newUser
}

func (s *Store) GetUser(id int) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if user, ok := s.users[id]; !ok {
		return user, false
	} else {
		return user, true
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		w,
		200,
	}
}

func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		logger.Info("http request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.statusCode),
			slog.Duration("duration", time.Since(start)),
		)
	})
}
