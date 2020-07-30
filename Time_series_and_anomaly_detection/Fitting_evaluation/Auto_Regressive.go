package main

import (
	"encoding/csv"
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/sajari/regression"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"log"
	"math"
	"os"
	"strconv"
)


//自回归以给定顺序为系列计算AR模型。
func autoregressive(x []float64, lag int) ([]float64, float64) {

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

	r.Run()

	// coeff保留我们的滞后系数。
	var coeff []float64
	for i := 1; i <= lag; i++ {
		coeff = append(coeff, r.Coeff(i))
	}

	return coeff, r.Coeff(0)
}

func main() {
	passengersFile, err := os.Open("Fitting_evaluation/log_diff_series.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer passengersFile.Close()

	passengersDF := dataframe.ReadCSV(passengersFile)

	passengers := passengersDF.Col("log_differenced_passengers").Float()

	//计算滞后1和滞后2的系数以及我们的误差。
	coeffs, intercept := autoregressive(passengers, 2)

	//将AR（2）模型输出到stdout。
	fmt.Printf("\nlog(x(t)) - log(x(t-1)) = %0.6f + lag1*%0.6f + lag2*%0.6f\n\n", intercept, coeffs[0], coeffs[1])


	transFile, err := os.Open("Fitting_evaluation/log_diff_series.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer transFile.Close()

	transReader := csv.NewReader(transFile)

	transReader.FieldsPerRecord = 2
	transData, err := transReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//遍历预测转换后的观测值的数据。
	var transPredictions []float64
	for i, _ := range transData {

		//跳过标题和前两个观察值（因为我们需要两个滞后才能进行预测）。
		if i == 0 || i == 1 || i == 2 {
			continue
		}

		//解析第一个滞后。
		lagOne, err := strconv.ParseFloat(transData[i-1][1], 64)
		if err != nil {
			log.Fatal(err)
		}


		//解析第二个滞后。
		lagTwo, err := strconv.ParseFloat(transData[i-2][1], 64)
		if err != nil {
			log.Fatal(err)
		}


		//使用我们训练有素的AR模型预测转换后的变量。
		transPredictions = append(transPredictions, 0.008159+0.234953*lagOne-0.173682*lagTwo)
	}

	origFile, err := os.Open("Fitting_evaluation/AirPassengers.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer origFile.Close()

	origReader := csv.NewReader(origFile)

	origReader.FieldsPerRecord = 2
	origData, err := origReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}


	// pts *将保存绘图值。
	ptsObs := make(plotter.XYs, len(transPredictions))
	ptsPred := make(plotter.XYs, len(transPredictions))

	//逆变换并计算MAE。
	var mAE float64
	var cumSum float64
	for i := 4; i <= len(origData)-1; i++ {

		//解析原始观察值。
		observed, err := strconv.ParseFloat(origData[i][1], 64)
		if err != nil {
			log.Fatal(err)
		}

		//解析原始日期。
		date, err := strconv.ParseFloat(origData[i][0], 64)
		if err != nil {
			log.Fatal(err)
		}


		//获取累积总和直到索引 转换后的预测。
		cumSum += transPredictions[i-4]

		////计算反向转换的预测。
		predicted := math.Exp(math.Log(observed) + cumSum)

		mAE += math.Abs(observed-predicted) / float64(len(transPredictions))

		//填入绘图点。
		ptsObs[i-4].X = date
		ptsPred[i-4].X = date
		ptsObs[i-4].Y = observed
		ptsPred[i-4].Y = predicted
	}

	fmt.Printf("\nMAE = %0.2f\n\n", mAE)

	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.X.Label.Text = "time"
	p.Y.Label.Text = "passengers"
	p.Add(plotter.NewGrid())

	//添加时间序列的折线图点。
	lObs, err := plotter.NewLine(ptsObs)
	if err != nil {
		log.Fatal(err)
	}
	lObs.LineStyle.Width = vg.Points(1)

	lPred, err := plotter.NewLine(ptsPred)
	if err != nil {
		log.Fatal(err)
	}
	lPred.LineStyle.Width = vg.Points(1)
	lPred.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}

	p.Add(lObs, lPred)
	p.Legend.Add("Observed", lObs)
	p.Legend.Add("Predicted", lPred)
	if err := p.Save(10*vg.Inch, 4*vg.Inch, "passengers_ts.png"); err != nil {
		log.Fatal(err)
	}
}
