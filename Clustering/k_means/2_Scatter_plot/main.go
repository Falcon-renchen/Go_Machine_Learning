package main

import (
	"github.com/kniren/gota/dataframe"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
	"os"
)

func main() {
	f, err := os.Open("fleet_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	driverDF := dataframe.ReadCSV(f)

	// 提取距离列
	yVals := driverDF.Col("Distance_Feature").Float()

	// pts将保留绘图的值
	pts := make(plotter.XYs, driverDF.Nrow())

	//用数据填充点数。
	for i, floatVal := range driverDF.Col("Speeding_Feature").Float() {
		pts[i].X = floatVal
		pts[i].Y = yVals[i]
	}

	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.X.Label.Text = "Speeding"
	p.Y.Label.Text = "Distance"
	//创建网格
	p.Add(plotter.NewGrid())

	//创建散点图
	s, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatal(err)
	}
	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
	s.GlyphStyle.Radius = vg.Points(3)

	p.Add(s)
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "fleet_data_scatter.png"); err != nil {
		log.Fatal(err)
	}
}
