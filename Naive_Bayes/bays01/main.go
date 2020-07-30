package main

import (
	"fmt"
	_ "github.com/sjwhitworth/golearn"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/filters"
	"github.com/sjwhitworth/golearn/naive"
	"log"
)

// convertToBinary利用内置的golearn功能将我们的标签转换为二进制标签格式。
func convertToBinary(src base.FixedDataGrid) base.FixedDataGrid {
	// NewBinaryConvertFilter创建一个空白的BinaryConvertFilter
	b := filters.NewBinaryConvertFilter()

	// NonClassAttrs返回所有未指定为类属性的属性。
	attrs := base.NonClassAttributes(src)
	for _, a := range attrs {
		b.AddAttribute(a)   // AddAttribute向此过滤器添加新的属性
	}
	b.Train()

	//将给定的Filter应用于其包含的属性后，
	//NewLazilyFitleredInstances返回一个新的FixedDataGrid。未经过滤的属性无需修改即可传递。
	ret := base.NewLazilyFilteredInstances(src,b)
	return ret
}

func main() {

	// ParseCSVToInstances读取文件路径给定的CSV文件，并返回读取的实例。
	trainingData, err := base.ParseCSVToInstances("training.csv",true)
	if err != nil {
		log.Fatal(err)
	}

	//初始化新的朴素贝叶斯分类器。
	nb := naive.NewBernoulliNBClassifier()

	//适配朴素贝叶斯分类器。
	nb.Fit(convertToBinary(trainingData))

	//将贷款测试数据集读入golearn“ instances”。这次，我们将使用前一组实例的模板来验证测试集的格式。
	testData, err := base.ParseCSVToInstances("test.csv",true)
	if err != nil {
		log.Fatal(err)
	}

	//做出我们的预测。
	prediction, err := nb.Predict(convertToBinary(testData))
	if err != nil {
		log.Fatal(err)
	}

	//生成混淆矩阵
	cm, err := evaluation.GetConfusionMatrix(testData,prediction)
	if err != nil {
		log.Fatal(err)
	}

	//检索精度。
	accuracy := evaluation.GetAccuracy(cm)

	fmt.Printf("\nAccuracy: %0.2f\n\n", accuracy)

}
