package main

import (
	"fmt"
	"github.com/gonum/floats"
	"github.com/kniren/gota/dataframe"
	"github.com/montanaflynn/stats"
	"gonum.org/v1/gonum/stat"
	"log"
	"os"
)

func main() {
	irisFile, err := os.Open("dataform/iris.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer irisFile.Close()

	irisDF := dataframe.ReadCSV(irisFile)

	sepalLength := irisDF.Col("sepal_length").Float()
	meanVal := stat.Mean(sepalLength, nil)	//平均值
	fmt.Println("meanVal: ",meanVal)

	//众数
	modeVal, modeCount := stat.Mode(sepalLength,nil)		// 分布
	fmt.Printf("modeVal: %f\nmodeCount: %f\n",modeVal,modeCount)


	medianVal, err := stats.Median(sepalLength)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("medianVal: %f",medianVal)

	maxVal := floats.Max(sepalLength)
	minVal := floats.Min(sepalLength)

	rangeVal := maxVal - minVal
	varianceVal := stat.Variance(sepalLength, nil)	//计算方差

	stdDevVal := stat.StdDev(sepalLength,nil)	//计算标准差

	//对值进行排列
	indx := make([]int, len(sepalLength))
	floats.Argsort(sepalLength, indx)

	//获取分位数
	quant25 := stat.Quantile(0.25, stat.Empirical, sepalLength, nil)
	quant50 := stat.Quantile(0.50, stat.Empirical, sepalLength, nil)
	quant75 := stat.Quantile(0.75, stat.Empirical, sepalLength, nil)


	fmt.Printf("\nSepal Length Summary Statistics:\n")
	fmt.Printf("Max value: %0.2f\n", maxVal)
	fmt.Printf("Min value: %0.2f\n", minVal)
	fmt.Printf("Range value: %0.2f\n", rangeVal)			//整个值的范围
	fmt.Printf("Variance value: %0.2f\n", varianceVal)	//方差
	fmt.Printf("Std Dev value: %0.2f\n", stdDevVal)		//标准差
	fmt.Printf("25 Quantile: %0.2f\n", quant25)
	fmt.Printf("50 Quantile: %0.2f\n", quant50)
	fmt.Printf("75 Quantile: %0.2f\n\n", quant75)




}
