package main

import (
	"fmt"
	"github.com/kniren/gota/dataframe"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"log"
	"os"
)

func main() {
	irisFile, err := os.Open("dataform/iris.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer irisFile.Close()
	irisDF := dataframe.ReadCSV(irisFile)

	//创建一个列，根据名字
	//为数据集中的每个要素列创建直方图。
	for _, colName := range irisDF.Names() {
		//如果该列是要素列之一，让我们创建值的直方图。

		//如果名字不是special
		if colName != "special" {
			//创建一个plotter.Values值，并用数据框各个列中的值填充它。
			v := make(plotter.Values, irisDF.Nrow())
			for i,floatVal := range irisDF.Col(colName).Float() {
				v[i] = floatVal
			}
			//创建一个直方图
			p, err := plot.New()
			if err != nil {
				log.Fatal(err)
			}
			p.Title.Text = fmt.Sprintf("Histogram of a %s", colName)

			//创建从标准法线绘制的值的直方图。
			h, err := plotter.NewHist(v, 16)
			if err != nil {
				log.Fatal(err)
			}
			//标准化直方图
			h.Normalize(1)
			//将直方图添加到图中
			p.Add(h)

			//将绘图保存到PNG文件。
			if err := p.Save(4*vg.Inch, 4*vg.Inch, colName+"_hist.png"); err != nil {
				log.Fatal(err)
			}
		}
	}

}
