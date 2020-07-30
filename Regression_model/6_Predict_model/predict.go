package main

import (
	"github.com/kniren/gota/dataframe"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"log"
	"os"
)

func predict(tv float64) float64 {
	return 7.0688 + tv*0.0489
}

func main() {
	f, err := os.Open("Advertising.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	advertDF := dataframe.ReadCSV(f)

	// y轴为 Sales的值
	yVals := advertDF.Col("Sales").Float()
	// 画散点图
	pts := make(plotter.XYs, advertDF.Nrow())
	// ptsPred将保存绘图的预测值。 应该是直线
	ptsPre := make(plotter.XYs, advertDF.Nrow())

	for i,floatVal := range advertDF.Col("TV").Float() {
		//散点图
		pts[i].X = floatVal
		pts[i].Y = yVals[i]
		// y为Sales ，  x为tv的值
		ptsPre[i].X = floatVal
		ptsPre[i].Y = predict(floatVal)
	}
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.X.Label.Text = "TV"
	p.Y.Label.Text = "Sales"

	p.Add(plotter.NewGrid())
	s, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatal(err)
	}
	s.GlyphStyle.Radius = vg.Points(3)

	l, err := plotter.NewLine(ptsPre)
	if err != nil {
		log.Fatal(err)
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5),vg.Points(5)}

	p.Add(s,l)

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "regression_line.png"); err != nil {
		log.Fatal(err)
	}
}
