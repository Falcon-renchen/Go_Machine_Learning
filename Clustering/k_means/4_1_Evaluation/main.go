package main

import (
	"github.com/gonum/floats"
	"github.com/kniren/gota/dataframe"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
	"os"
)

//评估集群的合法性

func main() {
	f, err := os.Open("fleet_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	driverDF := dataframe.ReadCSV(f)

	yVals := driverDF.Col("Distance_Feature").Float()

	// clusterOne和clusterTwo将保存用于绘制的值。
	var clusterOne [][]float64
	var clusterTwo [][]float64

	// 用数据填充群集。
	for i, xVal := range driverDF.Col("Speeding_Feature").Float() {
		//算出数据和质心的距离，然后比较，看数据距离哪个质心最近，就归为哪个点
		distanceOne := floats.Distance([]float64{yVals[i], xVal}, []float64{50.05, 8.83}, 2)
		distanceTwo := floats.Distance([]float64{yVals[i], xVal}, []float64{180.02, 18.29}, 2)
		if distanceOne < distanceTwo {
			clusterOne = append(clusterOne, []float64{xVal, yVals[i]})
			continue
		}
		clusterTwo = append(clusterTwo, []float64{xVal, yVals[i]})
	}
	// pts *将保存绘图值
	ptsOne := make(plotter.XYs, len(clusterOne))
	ptsTwo := make(plotter.XYs, len(clusterTwo))

	// 用数据填充点数。
	for i, point := range clusterOne {
		ptsOne[i].X = point[0]
		ptsOne[i].Y = point[1]
	}

	for i, point := range clusterTwo {
		ptsTwo[i].X = point[0]
		ptsTwo[i].Y = point[1]
	}

	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.X.Label.Text = "Speeding"
	p.Y.Label.Text = "Distance"
	p.Add(plotter.NewGrid())

	sOne, err := plotter.NewScatter(ptsOne)
	if err != nil {
		log.Fatal(err)
	}

	//将两个切片设置成不一样的颜色
	sOne.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
	sOne.GlyphStyle.Radius = vg.Points(3)

	sTwo, err := plotter.NewScatter(ptsTwo)
	if err != nil {
		log.Fatal(err)
	}
	sTwo.GlyphStyle.Color = color.RGBA{B: 255, A: 255}
	sTwo.GlyphStyle.Radius = vg.Points(3)

	// Save the plot to a PNG file.
	p.Add(sOne, sTwo)
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "fleet_data_clusters.png"); err != nil {
		log.Fatal(err)
	}
}
