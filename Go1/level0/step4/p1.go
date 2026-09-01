package main

import (
	"fmt"
	"math"
)

func main() {
	var v1, v2, v3 float64
	fmt.Scan(&v1, &v2, &v3)
	fmt.Println(math.Min(math.Min(v1, v2), v3))
}
