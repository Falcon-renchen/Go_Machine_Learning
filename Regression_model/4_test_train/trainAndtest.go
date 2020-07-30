package main

import (
	"bufio"
	"github.com/kniren/gota/dataframe"
	"log"
	"os"
)
//我们在此处使用的数据没有以任何方式按数据排序或排序。但是，如果要处理按响应，
//日期或其他任何方式排序的数据，则将数据随机分为训练和测试集很重要。
//如果您不这样做，则您的训练和测试集可能仅包括响应的某些范围，可能受时间/日期等的人为影响。
func main() {
	f, err := os.Open("Advertising.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	advertDF := dataframe.ReadCSV(f)

	trainingNum := (4 * advertDF.Nrow()) / 5
	testNum := advertDF.Nrow() / 5
	if trainingNum+testNum < advertDF.Nrow() {
		trainingNum++
	}

	trainingIdx := make([]int,trainingNum)
	testIdx := make([]int, testNum)

	//列举训练指标
	for i:=0; i<trainingNum; i++ {
		trainingIdx[i] = i
	}
	for i:=0; i<testNum; i++ {
		testIdx[i] = trainingNum+i
	}

	trainingDF := advertDF.Subset(trainingIdx)
	testDF := advertDF.Subset(testIdx)

	setMap := map[int]dataframe.DataFrame{
		0:trainingDF,
		1:testDF,
	}

	for idx,setName := range []string{"train.csv","test.csv"} {
		f, err := os.Create(setName)
		if err != nil {
			log.Fatal(err)
		}
		w := bufio.NewWriter(f)
		if err:=setMap[idx].WriteCSV(w); err != nil {
			log.Fatal(err)
		}
	}


}
