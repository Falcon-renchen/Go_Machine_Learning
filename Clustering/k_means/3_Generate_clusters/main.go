package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"github.com/mash/gokmeans"
	"strconv"
)

func main() {
	f, err := os.Open("fleet_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 3

	var data []gokmeans.Node

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if record[0] == "Driver_ID" {
			continue
		}
		//初始化point
		var point []float64

		//将数据填充到point中
		for i := 1; i < 3; i++ {
			val, err := strconv.ParseFloat(record[i], 64)
			if err != nil {
				log.Fatal(err)
			}

			//将此值附加到我们的观点。
			point = append(point, val)
		}

		//将我们的观点附加到数据上。
		data = append(data, gokmeans.Node{point[0], point[1]})
	}

	//用k均值生成聚类。
	success, centroids := gokmeans.Train(data, 2, 50)
	if !success {
		log.Fatal("Could not generate clusters")
	}
	fmt.Println("The centroids for our clusters are:")

	//为生成的集群提供以下质心
	for _, centroid := range centroids {
		fmt.Println(centroid)
	}
}
