package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

)
//处理意外数据,并且输出具体第几行第几列

type CSVRecord struct {
	SepalLength  float64
	SepalWidth   float64
	PetalLength  float64
	PetalWidth   float64
	Species      string
	ParseError   error
}

func main() {
	f, err := os.Open("dataform/iris_mixed_types.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	//创建一个切片值，该值将保存CSV中所有成功解析的记录。
	var csvData []CSVRecord

	// line将帮助我们跟踪记录的行号。
	line := 1

	//读取记录以查找意外类型。
	for {
		//连续阅读。 检查我们是否在文件末尾。
		record, err := reader.Read()
		if err==io.EOF {
			break
		}
		//为该行创建一个CSVRecord值。
		var csvRecord CSVRecord

		//根据期望的类型解析记录中的每个值。
		for idx, val := range record {
			//将记录中的值解析为字符串列的字符串。
			if idx == 4 {
				//验证该值不是一个空字符串。 如果该值为空字符串，则中断解析循环。
				if val == "" {
					log.Printf("Parsing line %d failed, unexpected type in column %d\n", line, idx)
					csvRecord.ParseError = fmt.Errorf("Empty string value")
					break
				}
				//将字符串值添加到CSVRecord。
				csvRecord.Species = val
				continue
			}
			//否则，将记录中的值解析为float64。
			//floatValue将保存数字列的记录的解析浮点值。
			var floatValue float64

			//如果无法将值解析为浮点型，请记录并中断解析循环。
			if floatValue, err = strconv.ParseFloat(val,64); err != nil {
				log.Printf("Parsing line %d failed, unexpected type in column %d\n", line, idx)
				csvRecord.ParseError = fmt.Errorf("Could not parse float")
				break
			}
			//将float值添加到CSVRecord中的相应字段。
			switch idx {
				case 0:
					csvRecord.SepalLength = floatValue
				case 1:
					csvRecord.SepalWidth = floatValue
				case 2:
					csvRecord.PetalLength = floatValue
				case 3:
					csvRecord.PetalWidth = floatValue
			}

		}
		//将成功解析的记录追加到上面定义的片上
		if csvRecord.ParseError == nil {
			csvData = append(csvData, csvRecord)
		}
		//递增行计数器。
		line++
	}

	fmt.Printf("successfully parsed %d lines\n", len(csvData))

}
