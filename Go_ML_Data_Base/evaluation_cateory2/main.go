package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

//我们可以将类别1视为正面， 将其他类别视为负面，
//将类别2视为正面，其他类别视为负面
func main() {
	f, err := os.Open("evaluate_cateory/labeled.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	var observed []int
	var predicted []int

	line := 1

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
			log.Printf("Parsing line %d failed, unexpected type\n", line)
			continue
		}

		predictedVal, err := strconv.Atoi(record[1])
		if err != nil {
			log.Printf("Parsing line %d failed, unexpected type\n", line)
			continue
		}

		observed = append(observed, observedVal)
		predicted = append(predicted, predictedVal)
		line++
	}

	//类在标记的数据中包含三个可能的类。
	classes := []int{0, 1, 2}

	for _, class := range classes {

		//这些变量将保存我们的真实肯定计数和错误肯定计数。
		var truePos int
		var falsePos int
		var falseNeg int

		//累加真实的正数和错误的正数。
		for idx, oVal := range observed {
			switch oVal {
			//如果观察到的值是相关的类别，则应检查是否预测了该类别。
			case class:
				if predicted[idx] == class {
					truePos++
					continue
				}

				falseNeg++

				//如果观察值是不同的类，则应检查是否预测到假阳性。
			default:
				if predicted[idx] == class {
					falsePos++
				}
			}
		}
		//计算精度。
		precision := float64(truePos) / float64(truePos+falsePos)

		//计算召回率。
		recall := float64(truePos) / float64(truePos+falseNeg)

		fmt.Printf("\nPrecision (class %d) = %0.2f", class, precision)
		fmt.Printf("\nRecall (class %d) = %0.2f\n\n", class, recall)

	}
}
