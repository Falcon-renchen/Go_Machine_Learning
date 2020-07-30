package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const citiBikeURL = "https://gbfs.citibikenyc.com/gbfs/en/station_status.json"

type stationData struct {
	LastUpdated int `json:"last_updated"`
	TTL int `json:"ttl"`
	Data struct{
		Stations []station `json:"stations"`
	} `json:"data"`
}

type station struct {
	ID string `json:"station_id"`
	NumBikesAvailable int    `json:"num_bikes_available"`
	NumBikesDisabled  int    `json:"num_bike_disabled"`
	NumDocksAvailable int    `json:"num_docks_available"`
	NumDocksDisabled  int    `json:"num_docks_disabled"`
	IsInstalled       int    `json:"is_installed"`
	IsRenting         int    `json:"is_renting"`
	IsReturning       int    `json:"is_returning"`
	LastReported      int    `json:"last_reported"`
	HasAvailableKeys  bool   `json:"eightd_has_available_keys"`
}
func main() {
	//从URL获取JSON响应。
	response, err := http.Get(citiBikeURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	//将响应的正文读取到[] byte中。
	body,err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var sd stationData

	//将JSON数据解组到变量中。
	if err := json.Unmarshal(body, &sd); err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("%+v\n\n", sd.Data.Stations[0])

/*	//整理数据。
	outputData, err := json.Marshal(sd)
	if err != nil {
		log.Fatal(err)
	}

	//将封送的数据保存到文件中。
	if err := ioutil.WriteFile("citibike.json", outputData, 0644); err != nil {
		log.Fatal(err)
	}*/
}
