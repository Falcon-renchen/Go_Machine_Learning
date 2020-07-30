package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)
//只检查第几行数据有错误。
func main() {
	f, err := os.Open("dataform/iris_unexpected_fields.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	reader.FieldsPerRecord = 5

	var rawCSVData [][]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			continue
		}
		rawCSVData = append(rawCSVData, record)
	}
	fmt.Printf("parsed %d lines successfully\n",len(rawCSVData))
}
