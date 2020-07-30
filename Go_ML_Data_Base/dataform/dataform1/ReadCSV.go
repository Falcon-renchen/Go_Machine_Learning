package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main()  {
	// Open the CSV.
	f, err := os.Open("dataform/myfile.csv")
	if err != nil {
		log.Fatal(err)
	}

	// Read in the CSV records.
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//在整数列中获取最大值。
	var intMax int
	for _, record := range records {

		//解析整数值。
		intVal, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatal(err)
		}

		//如果合适，请替换最大值。
		if intVal > intMax {
			intMax = intVal
		}
	}

	// Print the maximum value.
	fmt.Println(intMax)
}
