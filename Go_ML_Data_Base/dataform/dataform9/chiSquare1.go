package main

import (
	"fmt"
	"gonum.org/v1/gonum/stat"
)
//卡方统计量
func main() {
	observed := []float64{48, 52}
	expected := []float64{50, 50}

	chiSquare := stat.ChiSquare(observed, expected)
	fmt.Println(chiSquare)	//偶然概率
}
