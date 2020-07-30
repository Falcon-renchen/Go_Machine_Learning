package main

import (
	"encoding/csv"
	"fmt"
	"github.com/sajari/regression"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("4_test_train/train.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	reader.FieldsPerRecord = 4
	trainingData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	var r regression.Regression
	//设置观察值   ==   y
	r.SetObserved("Sales")
	//设置自变量   ==   x
	r.SetVar(0,"TV")

	for i,record := range trainingData {
		if i == 0 {
			continue
		}

		//解析销售回归指标
		yVals, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}
		//解析电视值
		tvVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		r.Train(regression.DataPoint(yVals, []float64{tvVal}))

	}
	r.Run()

	fmt.Printf("\nRegression Formula:\n%v\n\n", r.Formula)
	//Predicted = 7.9819 + TV*0.0612


	f, err = os.Open("4_test_train/test.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	reader = csv.NewReader(f)

	reader.FieldsPerRecord = 4
	testData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	var mAe float64
	for i,record := range testData {
		if i == 0 {
			continue
		}
		// y   == Sales
		yObserved, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		// TV
		tvVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}
		yPredict, err := r.Predict([]float64{tvVal})
		if err != nil {
			log.Fatal(err)
		}
		mAe += math.Abs(yObserved-yPredict) / float64(len(testData))
	}
	//MAE小于我们的销售价值的标准差，约为平均值的20％,所以有预测功能
	fmt.Printf("MAE = %0.2f\n\n", mAe)
}
