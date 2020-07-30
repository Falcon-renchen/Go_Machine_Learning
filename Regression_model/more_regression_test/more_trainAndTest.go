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
	//y轴  观测值 为 Sales
	r.SetObserved("Sales")
	r.SetVar(0,"TV")
	r.SetVar(1,"Radio")

	for i, record := range trainingData {
		if i == 0 {
			continue
		}
		yVal, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}
		tvVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		radioVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		r.Train(regression.DataPoint(yVal,[]float64{tvVal,radioVal}))
	}
	r.Run()
	fmt.Printf("\nRegression Formula:\n%v\n\n", r.Formula)


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

	var mAE float64
	for i,record := range testData {
		if i == 0 {
			continue
		}
		yObserved, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}
		tvVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}
		radioVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		yPredict,err := r.Predict([]float64{tvVal,radioVal})
		if err != nil {
			log.Fatal(err)
		}

		mAE += math.Abs(yObserved-yPredict) / float64(len(testData))
	}
	//如果MAE小于平均值的20%左右，可预测
	fmt.Printf("MAE = %0.2f\n\n", mAE)

}
