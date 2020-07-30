package main

import (
	"fmt"
	"github.com/gonum/floats"
)
func main()  {
	distance := floats.Distance([]float64{1,2},[]float64{3,4},2)  //	L是平方
	fmt.Printf("distance = %0.2f",distance)
}
