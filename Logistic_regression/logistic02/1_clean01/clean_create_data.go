package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
scoreMax = 830.0
scoreMin = 640.0
)

func main() {

	f, err := os.Open("loan_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	//创建一个csv阅读器
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 2

	//开始阅读并记录
	rawCSVData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	f, err = os.Create("clean_loan_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)

	//顺序移动写出解析值的行。
	for idx, record := range rawCSVData {
		if idx==0 {
			 //写文件头
			 if err := w.Write(record); err!=nil {
			 	log.Fatal(err)
			 }
			 continue
		}
		//初始化一个切片以保存我们的解析值。
		outRecord := make([]string, 2)
		//解析并标准化FICO分数。
		score, err := strconv.ParseFloat(strings.Split(record[0],"-")[0],64)
		if err != nil {
			log.Fatal(err)
		}
		outRecord[0] = strconv.FormatFloat((score-scoreMin)/(scoreMax-scoreMin),'f',4,64)

		//解析利率类。把%去除
		rate, err := strconv.ParseFloat(strings.TrimSuffix(record[1],"%"),64)
		if err != nil {
			log.Fatal(err)
		}

		if rate<12.0 {
			outRecord[1] = "1.0"
			if err := w.Write(outRecord);err != nil {
				log.Fatal(err)
			}
			continue
		}
		outRecord[1] = "0.0"
		if err := w.Write(outRecord); err != nil {
			log.Fatal(err)
		}

		//将所有缓冲的数据写入基础写入器（标准输出）。
		w.Flush()
		if err:=w.Error();err!=nil {
			log.Fatal(err)
		}
	}

}
