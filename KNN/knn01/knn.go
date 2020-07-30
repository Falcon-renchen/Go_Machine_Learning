package main

import (
	"fmt"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	knn2 "github.com/sjwhitworth/golearn/knn"
	"log"
	"math"
)

func main() {
	//使用golearn时候，必须先将数据转化为instance
	irisData, err := base.ParseCSVToInstances("iris.csv",true)
	if err != nil {
		log.Fatal(err)
	}

	//初始化新的KNN分类器。我们将使用一个简单的欧几里德距离度量，并且k = 2。
	knn := knn2.NewKnnClassifier("euclidean","linear",2)

	//使用交叉折叠验证以对数据集的5个折叠进行连续训练和评估模型。
	cv, err := evaluation.GenerateCrossFoldValidationConfusionMatrices(irisData,knn,5)
	if err != nil {
		log.Fatal(err)
	}

	//获取交叉验证准确性的均值，方差和标准偏差。
	mean, variance := evaluation.GetCrossValidatedMetric(cv,evaluation.GetAccuracy)
	stdev := math.Sqrt(variance)

	//将交叉指标输出到标准输出。
	fmt.Printf("\nAccuracy\n%.2f (+/- %.2f)\n\n", mean, stdev*2)

}
