package main

import (
	"fmt"
	"time"
)

func main() {
	var s1, s2 string
	fmt.Scan(&s1, &s2)

	fut, _ := time.Parse("2006-01-02", s1)
	now, _ := time.Parse("2006-01-02", s2)

	fmt.Printf("%d year ago", fut.Year() - now.Year())
}
