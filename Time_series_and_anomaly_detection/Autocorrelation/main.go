package main

import (
	"fmt"
	"github.com/kniren/gota/dataframe"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"log"
	"math"
	"os"
)
//自动回归模型：这是一种试图通过同一过程的一个或多个延迟或滞后版本来对时间序列过程进行建模的模型。
//例如，股价的自回归模型将尝试通过先前时间间隔内的股价值对股价进行建模。

//acf计算给定滞后序列的自相关。

//重要的是自相关

//自相关是度量信号与自身延迟版本之间的相关性的度量。
//例如，一个或多个先前对股票价格的观察可以与下一个对该股票价格的观察相关（或一起改变）。
//如果是这样的话，我们可以说股票价格受到一些滞后或延迟的影响。
//然后，我们可以通过延迟的自身模型以指示为高度相关的特定滞后对未来股价进行建模。
//要测量变量x t与自身的延迟版本（或带有滞后的版本）x s的自相关，我们可以利用自相关函数（ACF）

func acf(x []float64,lag int) float64 {
	xAdj := x[lag:len(x)]
	xLag := x[0 : len(x)-lag]


	//分子将保存我们的累积分子，分母将保存我们的累积分母。
	var numerator float64
	var denominator float64


	//计算x值的平均值，将在自相关的每个项中使用。
	xBar := stat.Mean(x, nil)

	//计算分子。
	for idx, xVal := range xAdj {
		numerator += (xVal - xBar) * (xLag[idx] - xBar)
	}

	//计算分母。
	for _, xVal := range x {
		denominator += math.Pow(xVal-xBar, 2)
	}

	return numerator / denominator
}

func main() {
	passengersFile, err := os.Open("Autocorrelation/AirPassengers.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer passengersFile.Close()

	passengersDF := dataframe.ReadCSV(passengersFile)

	passengers := passengersDF.Col("AirPassengers").Float()

	fmt.Println("Autocorrelation:")

	//遍历序列中的各种滞后值。
	for i := 0; i < 11; i++ {
		// 转移系列。
		adjusted := passengers[i:len(passengers)]
		lag := passengers[0 : len(passengers)-i]

		//计算自相关。
		ac := stat.Correlation(adjusted, lag, nil)
		fmt.Printf("Lag %d period: %0.2f\n", i, ac)

	}

	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}

	p.Title.Text = "Autocorrelations for AirPassengers"
	p.X.Label.Text = "Lag"
	p.Y.Label.Text = "ACF"
	p.Y.Min = 0
	p.Y.Max = 1

	w := vg.Points(3)

	//创建绘图点。
	numLags := 20
	pts := make(plotter.Values, numLags)

	//遍历序列中的各种滞后值。
	for i := 1; i <= numLags; i++ {

		// 计算自相关。
		pts[i-1] = acf(passengers, i)
	}

	//将点添加到图中。
	bars, err := plotter.NewBarChart(pts, w)
	if err != nil {
		log.Fatal(err)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = plotutil.Color(1)

	p.Add(bars)
	if err := p.Save(8*vg.Inch, 4*vg.Inch, "acf.png"); err != nil {
		log.Fatal(err)
	}
}