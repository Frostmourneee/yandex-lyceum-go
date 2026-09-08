package main

import (
	"fmt"
	"reflect"
)

func main() {
	myFunc("hello")
	myFunc(42)
}

func myFunc(a any) {
	t := reflect.TypeOf(a)
	fmt.Printf("Type of '%v' is %v\n", a, t)
}

type TestNumeric interface {
	int | float64
}

func Sum[T TestNumeric](arr []T) T {
	var sum T
	for _, v := range arr {
		sum += v
	}
	return sum
}

func Filter[T any](arr []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range arr {
		if (predicate(v)) {
			result = append(result, v)
		}
	}
	return result
}

