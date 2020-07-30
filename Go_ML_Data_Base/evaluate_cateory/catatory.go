package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

//类别1可能对应于欺诈性交易，
//类别2可能对应于非欺诈性交易，
//而类别3可能对应于无效交易
func main() {
	f, err := os.Open("evaluate_cateory/labeled.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	var observed []int
	var predicted []int

	////行将跟踪行号以进行记录。
	line := 1
	//读入记录以查找列中的意外类型。
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if line == 1 {
			line++
			continue
		}

		observedVal, err := strconv.Atoi(record[0])
		if err != nil {
			panic(err)
		}
		predictedVal, err := strconv.Atoi(record[1])
		if err != nil {
			panic(err)
		}

		//如果记录具有预期的类型，则将其追加到切片中。
		observed = append(observed, observedVal)
		predicted = append(predicted, predictedVal)
		line++
	}
		//此变量将保存我们的真实正值和真实负值。
		var truePosNeg int

		//累积真实的正数/负数。
		//判断观测值和预测值是否相等，如果相等，truePosNeg++
		for idx, oVal := range observed {
			if oVal==predicted[idx] {
				truePosNeg++
			}
		}
		//计算精度（子集精度）。
		accuracy := float64(truePosNeg)/float64(len(observed))
		//将Accuracy值输出到标准输出。
		fmt.Printf("\nAccuracy = %0.2f\n\n", accuracy)
}

