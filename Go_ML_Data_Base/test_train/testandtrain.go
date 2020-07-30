package main

import (
	"bufio"
	"github.com/kniren/gota/dataframe"
	"log"
	"os"
)

func main() {
	f, err := os.Open("test_train/diabetes.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	diabeteDF := dataframe.ReadCSV(f)

	/*********************************************************/

	//计算每个集合中的元素数。
	// Nrow返回DataFrame上的行数。
	trainNum := (4 * diabeteDF.Nrow()) / 5
	testNum := diabeteDF.Nrow() / 5
	if trainNum+testNum < diabeteDF.Nrow() {
		trainNum++
	}

	//创建子集索引
	trainingIdx := make([]int, trainNum)
	testIdx := make([]int, testNum)

	//列举训练指标
	for i:=0; i<trainNum; i++ {
		trainingIdx[i] = i
	}
	//列举测试指标
	for i:=0; i<testNum; i++ {
		testIdx[i] = trainNum + i
	}
	//创建子集数据框
	trainingDF := diabeteDF.Subset(trainingIdx)
	testDF := diabeteDF.Subset(testIdx)

	//创建用于写入数据的map到文件。
	setMap := map[int]dataframe.DataFrame{
		0:trainingDF,
		1:testDF,
	}
	//创建相应的文件
	for idx, setName := range []string{"traning.csv", "test.csv"} {
		//保存过滤的数据集文件。
		f, err := os.Create(setName)
		if err != nil {
			log.Fatal(err)
		}
		//创建一个缓冲的
		w := bufio.NewWriter(f)
		//将数据框写为CSV。
		if err := setMap[idx].WriteCSV(w); err != nil {
			log.Fatal(err)
		}
	}
}
