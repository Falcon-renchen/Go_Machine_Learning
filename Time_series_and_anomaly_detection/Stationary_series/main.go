package main

import (
	"encoding/csv"
	"github.com/go-gota/gota/dataframe"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
	"os"
	"strconv"
)

//我们可以应用一个常见的技巧来使我们的序列平稳，称为差分。
//可视化数据

func main() {
	passengersFile, err := os.Open("Stationary_series/AirPassengers.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer passengersFile.Close()

	passengersDF := dataframe.ReadCSV(passengersFile)

	passengerVals := passengersDF.Col("AirPassengers").Float()
	timeVals := passengersDF.Col("time").Float()

	pts := make(plotter.XYs, passengersDF.Nrow()-1)

	//差异将保存我们的差异值，这些差异值将输出到新的CSV文件中。
	//差异值是 后一个人数 减去  前一个人数
	var differenced [][]string
	differenced = append(differenced, []string{"time", "differenced_passengers"})

	for i := 1; i < len(passengerVals); i++ {
		pts[i-1].X = timeVals[i]
		pts[i-1].Y = passengerVals[i] - passengerVals[i-1]
		differenced = append(differenced, []string{
			strconv.FormatFloat(timeVals[i], 'f', -1, 64),
			strconv.FormatFloat(passengerVals[i]-passengerVals[i-1], 'f', -1, 64),
		})
	}

	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.X.Label.Text = "time"
	p.Y.Label.Text = "differenced passengers"
	p.Add(plotter.NewGrid())

	l, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatal(err)
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	p.Add(l)
	if err := p.Save(10*vg.Inch, 4*vg.Inch, "diff_passengers_ts.png"); err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("diff_series.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.WriteAll(differenced)

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

}
