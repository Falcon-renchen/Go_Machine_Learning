package main

import (
	"fmt"
	"github.com/kniren/gota/dataframe"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
	"os"
)

func main() {
	passengersFile, err := os.Open("Time_series/AirPassengers.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer passengersFile.Close()
	passengersDF := dataframe.ReadCSV(passengersFile)
	fmt.Println(passengersDF)
	//提取乘客人数列。
	yVals := passengersDF.Col("value").Float()

	pts := make(plotter.XYs,passengersDF.Nrow())
	//填充数据
	for i,floatVal := range passengersDF.Col("time").Float() {
		pts[i].X = floatVal
		pts[i].Y = yVals[i]
	}

	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.X.Label.Text = "time"
	p.Y.Label.Text = "passengers"
	p.Add(plotter.NewGrid())

	//添加时间序列的折线图点。
	l, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatal(err)
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Color = color.RGBA{B:255, A:255}

	p.Add(l)

	if err := p.Save(vg.Inch*10,vg.Inch*4,"passengers_ts.png");err!=nil {
		log.Fatal(err)
	}

}
