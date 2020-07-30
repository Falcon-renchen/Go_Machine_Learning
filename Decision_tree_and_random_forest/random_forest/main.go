package main

import (
	"fmt"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/ensemble"
	"github.com/sjwhitworth/golearn/evaluation"
	"log"
	"math"
	"math/rand"
)

func main() {
	irisData, err := base.ParseCSVToInstances("iris.csv",true)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(44111342)

	//组装一个随机森林，该森林包含10棵树和每棵树2个特征，这是明智的默认设置（通常设置每棵树的特征数到sqrt（功能数））。
	rf := ensemble.NewRandomForest(10,2)

	//使用交叉折叠验证以对数据集的5个折叠进行连续训练和评估模型。
	cv, err := evaluation.GenerateCrossFoldValidationConfusionMatrices(irisData,rf,4)
	if err != nil {
		log.Fatal(err)
	}
	//获取交叉验证准确性的均值，方差和标准偏差。
	mean, variance := evaluation.GetCrossValidatedMetric(cv,evaluation.GetAccuracy)
	stdev := math.Sqrt(variance)

	fmt.Printf("\nAccuracy\n%.2f (+/- %.2f)\n\n", mean, stdev*2)
}
