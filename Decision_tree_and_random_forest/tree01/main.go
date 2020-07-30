package main

import (
	"fmt"
	_ "github.com/sjwhitworth/golearn"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/trees"
	"log"
	"math"
	"math/rand"
)

//最好坚持使用单一决策树。这个单一的决策树也更具解释性和效率。

func main() {
	irisData, err := base.ParseCSVToInstances("iris.csv",true)
	if err != nil {
		log.Fatal(err)
	}
	//这是播种与构建决策树有关的随机过程的种子。
	rand.Seed(44111342)

	//我们将使用ID3算法来构建决策树。还有，我们将以参数0.6开头，该参数控制火车修剪分裂。
	tree := trees.NewID3DecisionTree(0.6)

	//使用交叉折叠验证以对数据集的5个折叠进行连续训练和评估模型。
	//GenerateCrossFoldValidationConfusionMatrices将数据分为多个折叠，然后训练并评估每个折叠上的分类器，从而生成新的ConfusionMatrix。
	cv, err := evaluation.GenerateCrossFoldValidationConfusionMatrices(irisData,tree,2)
	if err != nil {
		log.Fatal(err)
	}
	//获取交叉验证准确性的均值，方差和标准偏差。
	mean, variance := evaluation.GetCrossValidatedMetric(cv,evaluation.GetAccuracy)
	stdev := math.Sqrt(variance)

	fmt.Printf("\nAccuracy\n%.2f (+/- %.2f)\n\n", mean, stdev*2)

}
