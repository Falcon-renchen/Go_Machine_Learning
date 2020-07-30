package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gonum/matrix/mat64"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

//features：指向gonum mat64.Dense矩阵的指针。
//该矩阵包括一列我们正在使用的任何自变量（在我们的情况下为FICO评分）以及代表截距的1.0s列。
//labels：包含对应我们的所有类标签的浮动片段features。
//numSteps：优化的最大迭代次数。
//learningRate：可调整的参数，有助于优化的收敛。

func logistic(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}
//为了优化系数/权重，我们将使用一种称为随机梯度下降的技术

//该函数输出逻辑回归模型的优化权重：
// logisticRegression适合给定数据的逻辑回归模型。
func logisticRegression(features *mat64.Dense, labels []float64, numSteps int, learningRate float64) []float64 {

	//初始化随机权重
	_, numWeights := features.Dims()	// Dims返回矩阵中的行数和列数。
	weights := make([]float64, numWeights)

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for idx := range weights {
		weights[idx] = r.Float64()
	}

	// 迭代优化权重
	for i := 0; i < numSteps; i++ {

		//初始化变量以积累此迭代的错误。
		var sumError float64

		//为每个标签进行预测并累积错误。
		for idx, label := range labels {

			//获取与此标签相对应的特征。
			featureRow := mat64.Row(nil, idx, features)

			//计算该迭代权重的误差。
			pred := logistic(featureRow[0]*weights[0] + featureRow[1]*weights[1])
			predError := label - pred
			sumError += math.Pow(predError, 2)

			// 更新特征权重
			for j := 0; j < len(featureRow); j++ {
				weights[j] += learningRate * predError * pred * (1 - pred) * featureRow[j]
			}
		}
	}

	return weights
}

func main() {
	// Open the training dataset file.
	f, err := os.Open("train.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 2

	rawCSVData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// FeatureData和标签将保存所有最终将在我们的培训中使用的float值。
	featureData := make([]float64, 2*(len(rawCSVData)-1))
	labels := make([]float64, len(rawCSVData)-1)


	// featureIndex将跟踪要素矩阵值的当前索引。
	var featureIndex int

	// 依次将行移到浮点数切片中。
	for idx, record := range rawCSVData {

		if idx == 0 {
			continue
		}

		// Add the FICO score feature.
		featureVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		featureData[featureIndex] = featureVal

		// 添加拦截
		featureData[featureIndex+1] = 1.0

		// 增加功能行
		featureIndex += 2

		// 添加类标签
		labelVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		labels[idx-1] = labelVal
	}

	// 从特征形成矩阵
	features := mat64.NewDense(len(rawCSVData)-1, 2, featureData)

	// Train the logistic regression model.
	weights := logisticRegression(features, labels, 100, 0.3)

	// Output the Logistic Regression model formula to stdout.
	formula := "p = 1 / ( 1 + exp(- m1 * FICO.score - m2) )"
	fmt.Printf("\n%s\n\nm1 = %0.2f\nm2 = %0.2f\n\n", formula, weights[0], weights[1])
}
