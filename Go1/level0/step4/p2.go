package main

import (
	"fmt"
	"math"
)

func main() {
	var x, y, z, m, n float64
	fmt.Scan(&x, &y, &z, &m, &n)
	fmt.Println((5 * x * math.Sin(2*y)) / (z + math.Pow(n, math.Log(m))))
}
