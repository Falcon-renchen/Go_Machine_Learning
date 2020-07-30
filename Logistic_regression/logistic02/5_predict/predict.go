package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

//预测基于我们训练有素的逻辑回归模型进行预测。
func predict(score float64) float64 {

	//计算预测概率。
	p := 1 / (1 + math.Exp(-13.6*score+4.89))

	// 输出相应的类
	if p >= 0.5 {
		return 1.0
	}

	return 0.0
}

func main() {
	f, err := os.Open("test.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	//观察到的和预测的将保存标记的数据文件中已解析的观察到的和预测的值。
	var observed []float64
	var predicted []float64

	// 行将跟踪行号以进行记录.
	line := 1

	// 读入记录以查找列中的意外类型
	for {

		// 连续阅读。检查我们是否在文件末尾
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		// Skip the header.
		if line == 1 {
			line++
			continue
		}

		// Read in the observed value.
		observedVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Printf("Parsing line %d failed, unexpected type\n", line)
			continue
		}

		//做出相应的预测.
		score, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Printf("Parsing line %d failed, unexpected type\n", line)
			continue
		}

		predictedVal := predict(score)

		// 如果记录具有预期的类型，则将其追加到我们的切片中。
		observed = append(observed, observedVal)
		predicted = append(predicted, predictedVal)
		line++
	}

	//此变量将保存我们的真实正值和真实负值。
	var truePosNeg int

	// 累计真实的正数/负数。

	//通过predict判断
	for idx, oVal := range observed {
		if oVal == predicted[idx] {
			truePosNeg++
		}
	}

	// 计算精度（子集精度）
	accuracy := float64(truePosNeg) / float64(len(observed))

	// Output the Accuracy value to standard out.
	fmt.Printf("\nAccuracy = %0.2f\n\n", accuracy)


}
