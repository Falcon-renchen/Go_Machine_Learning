package main

import (
	"encoding/csv"
	"fmt"
	"gonum.org/v1/gonum/stat"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

//MSE或均方差（MSD）是所有误差的平方的平均值
//MAE是所有误差的绝对值的平均值
//MSE对异常值更为敏感。
func main() {
	f, err := os.Open("evaluate_success/continuous_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)

	//观察和预测将保留从连续数据文件中解析出的观察和预测值。
	var observed []float64
	var predicted []float64

	line := 1

	//读入记录以查找列中的意外类型。
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		//跳过头一行
		if line == 1 {
			line++
			continue
		}
		//读取观测值和预测值
		observedVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
			continue
		}
		predictedVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
			continue
		}

		//如果记录具有预期的类型，则将其追加到切片中。
		observed = append(observed, observedVal)
		predicted = append(predicted, predictedVal)
		line++
	}
		//计算平均绝对误差和均方误差。
		var mAE float64
		var mSE float64
		for idx, oVal := range observed {
			mAE += math.Abs(oVal-predicted[idx]) / float64(len(observed))
			mSE += math.Pow(oVal-predicted[idx], 2) / float64(len(observed))
		}
		//为了判断这些值是否合适，我们需要将它们与观察数据中的值进行比较
		//MAE为2.55，我们观测值的平均值为14.0，因此我们的MAE约为平均值的20％。不是很好，
		fmt.Printf("\nMAE = %0.2f\n", mAE)
		fmt.Printf("\nMSE = %0.2f\n\n", mSE)


	// R^2 value.用作连续变量模型的评估指标
	//R平方测量我们在预测值中捕获的观测值中方差的比例
	//R平方是百分比，较高的百分比更好。在这里，我们捕获了我们试图预测的变量中约37％的方差。
	//不是很好。
	rSquared := stat.RSquaredFrom(observed, predicted, nil)


	fmt.Printf("\nR^2 = %0.2f\n\n", rSquared)
}
