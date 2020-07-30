package main

import (
	"fmt"
	"github.com/kniren/gota/dataframe"
	"github.com/sajari/regression"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"log"
	"os"
	"strconv"
)

func pacf(x []float64, lag int) float64 {


	var r regression.Regression
	r.SetObserved("x")


	//定义当前滞后和所有中间滞后。
	for i := 0; i < lag; i++ {
		r.SetVar(i, "x"+strconv.Itoa(i))
	}

	//移动系列。
	xAdj := x[lag:len(x)]


	//遍历系列，为回归创建数据集。
	for i, xVal := range xAdj {


		//循环中间的滞后以建立我们的自变量。
		laggedVariables := make([]float64, lag)
		for idx := 1; idx <= lag; idx++ {

			//获取滞后的系列变量。
			laggedVariables[idx-1] = x[lag+i-idx]
		}

		//将这些点添加到回归值中。
		r.Train(regression.DataPoint(xVal, laggedVariables))
	}


	//拟合回归。
	r.Run()

	return r.Coeff(lag)
}

func main() {

	passengersFile, err := os.Open("Partial_autocorrelation/AirPassengers.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer passengersFile.Close()

	passengersDF := dataframe.ReadCSV(passengersFile)

	passengers := passengersDF.Col("AirPassengers").Float()

	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}

	p.Title.Text = "Partial Autocorrelations for log(differenced passengers)"
	p.X.Label.Text = "Lag"
	p.Y.Label.Text = "PACF"
	p.Y.Min = 15
	p.Y.Max = -1

	w := vg.Points(3)

	// Create the points for plotting.
	numLags := 20
	pts := make(plotter.Values, numLags)


	// 遍历序列中的各种滞后值。
	fmt.Println("Partial Autocorrelation:")
	for i := 1; i < 11; i++ {


		//计算偏自相关。
		pac := pacf(passengers, i)
		fmt.Printf("Lag %d period: %0.2f\n", i, pac)
	}

	for i := 1; i <= numLags; i++ {

		// Calculate the partial autocorrelation.
		pts[i-1] = pacf(passengers, i)
	}

	// Add the points to the plot.
	bars, err := plotter.NewBarChart(pts, w)
	if err != nil {
		log.Fatal(err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = plotutil.Color(1)

	// Save the plot to a PNG file.
	p.Add(bars)
	if err := p.Save(8*vg.Inch, 4*vg.Inch, "pacf.png"); err != nil {
		log.Fatal(err)
	}
}
