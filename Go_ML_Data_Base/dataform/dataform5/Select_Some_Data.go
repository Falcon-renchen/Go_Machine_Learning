package main

import (
	"fmt"
	"github.com/kniren/gota/dataframe"
	"log"
	"os"
)

func main() {
	irisFile, err := os.Open("dataform/iris_labeled.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer irisFile.Close()

	irisDF := dataframe.ReadCSV(irisFile)
	fmt.Println(irisDF)

	//为数据框创建过滤器。就是挑选出来的数据的最后一列
	filter := dataframe.F{
		Colname:    "species",
		Comparator: "==",
		Comparando: "Iris-versicolor",
	}

	//过滤数据框以仅查看虹膜种类为“ Iris-versicolor”的行。
	//查看Comparando设置的那一列的某种属性的
	versiColorDF := irisDF.Filter(filter)
	if versiColorDF.Err != nil {
		log.Fatal(versiColorDF.Err)
	}

	//将结果输出到标准输出。
	fmt.Println(versiColorDF)

	//查看这两列的属性
	versiColorDF = irisDF.Filter(filter).Select([]string{"sepal_width","species"})
	fmt.Println(versiColorDF)

	//查看这两列属性的前三个(第0，1，2行的数据)
	versiColorDF = irisDF.Filter(filter).Select([]string{"sepal_width","species"}).Subset([]int{0,1,2})
	fmt.Println(versiColorDF)
}
