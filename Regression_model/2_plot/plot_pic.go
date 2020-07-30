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
	f, err := os.Open("Advertising.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	advertDF := dataframe.ReadCSV(f)

	for _, colName := range advertDF.Names() {
		plotVals := make(plotter.Values, advertDF.Nrow()) //行    空数组
		for i, floatVal := range advertDF.Col(colName).Float() {
			plotVals[i] = floatVal
		}
		//创建直方图
		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.Title.Text = fmt.Sprintf("Histogram of a %s", colName)
		h, err := plotter.NewHist(plotVals, 16)
		if err != nil {
			log.Fatal(err)
		}
		h.Normalize(1)
		p.Add(h)
		if err := p.Save(4*vg.Inch,4*vg.Inch,colName + "_hist.png"); err!=nil {
			log.Fatal(err)
		}
	}
}
