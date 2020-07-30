package main

import (
	"fmt"
	"gonum.org/v1/gonum/integrate"
	"gonum.org/v1/gonum/stat"
)

func main() {
	score := []float64{0.1,0.35,0.4,0.8}
	classes := []bool{true,false,true,false}

	//计算召回率，误报率
	tpr, fpr, _ := stat.ROC(nil,score,classes,nil)
	//计算曲线下的面积
	auc := integrate.Trapezoidal(fpr,tpr)
	fmt.Printf("true positive rate: %v\n", tpr)
	fmt.Printf("false positive rate: %v\n", fpr)
	fmt.Printf("auc: %v\n", auc)

}
