package main

import (
	"bufio"
	"github.com/kniren/gota/dataframe"
	"log"
	"os"
)

func main() {
	f, err := os.Open("clean_loan_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	loanDF := dataframe.ReadCSV(f)

	trainingNum := (4 * loanDF.Nrow()) / 5
	testingNum := loanDF.Nrow() / 5

	if trainingNum+testingNum<loanDF.Nrow() {
		trainingNum++
	}

	traningIdx := make([]int, trainingNum)
	testIdx := make([]int, testingNum)

	//列举训练指标
	for i:=0; i<trainingNum; i++ {
		traningIdx[i] = i
	}

	//列举测试指标
	for i:=0; i<testingNum; i++ {
		testIdx[i] = i
	}

	//创建子集数据框。
	traningDF := loanDF.Subset(traningIdx)
	testDF := loanDF.Subset(testIdx)

	//创建一个映射，该映射将用于将数据写入文件
	setMap := map[int]dataframe.DataFrame{
		0:traningDF,
		1:testDF,
	}

	for idx, setName := range []string{"train.csv","test.csv"} {
		f, err := os.Create(setName)
		if err != nil {
			log.Fatal(err)
		}
		w := bufio.NewWriter(f)
		if err := setMap[idx].WriteCSV(w); err!=nil {
			log.Fatal(err)
		}

	}

}
