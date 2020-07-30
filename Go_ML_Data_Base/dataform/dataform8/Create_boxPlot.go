package main

import (
	"github.com/kniren/gota/dataframe"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"log"
	"os"
)
////如果一个盒子比另一个盒子大，则表示分布偏斜.
//箱形图还包括两个尾巴或胡须,
//与包含大多数值（中间值为50％）的区域相比，这些可以使我们快速直观地看到分布范围。
func main() {
	irisFile, err := os.Open("dataform/iris.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer irisFile.Close()

	irisDF := dataframe.ReadCSV(irisFile)

	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.Title.Text = "Box plots"
	p.Y.Label.Text = "Values"
	//创建用于我们的数据的框。
	w := vg.Points(50)

	for idx, colName := range irisDF.Names() {
		//如果该列是要素列之一，让我们创建值的直方图。

		if colName != "special" {
			v := make(plotter.Values, irisDF.Nrow())
			for i, floatVal := range irisDF.Col(colName).Float() {
				v[i] = floatVal
			}
			//将数据添加到绘图中。
			b, err := plotter.NewBoxPlot(w,float64(idx),v)
			if err != nil {
				log.Fatal(err)
			}
			p.Add(b)
		}
	}

	//使用x = 0，x = 1等的给定名称，将图的X轴设置为标称值。
	p.NominalX("sepal_length", "sepal_width", "petal_length", "petal_width")
	//gonum在箱形图中包含了几个离群值（标记为圆形或点）。
	//许多绘图包都包括这些。它们表示与分布中位数至少相距一定距离的值。
	if err := p.Save(6*vg.Inch, 8*vg.Inch, "boxplots.png"); err != nil {
		log.Fatal(err)
	}
}
