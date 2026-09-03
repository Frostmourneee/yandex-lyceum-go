package main

import (
	"fmt"
	"math"
)

func main() {
	var v1 float64
	fmt.Scan(&v1)

	if v1 < 0 {
		fmt.Println(-1)
	} else {
		fmt.Printf("%.3f", math.Sqrt(v1))
	}
}
