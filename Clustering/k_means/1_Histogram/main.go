package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"log"
	"os"

	"github.com/kniren/gota/dataframe"
)

func main() {
	driverDataFile, err := os.Open("fleet_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer driverDataFile.Close()

	driverDF := dataframe.ReadCSV(driverDataFile)

	driverSummary := driverDF.Describe()
	fmt.Println(driverSummary)

	for _, colName := range driverDF.Names() {

		//创建一个plotter.Values值，并用数据框相应列中的值填充它。
		plotVals := make(plotter.Values, driverDF.Nrow())

		for i, floatVal := range driverDF.Col(colName).Float() {
			plotVals[i] = floatVal
		}

		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.Title.Text = fmt.Sprintf("1_Histogram of %s", colName)

		//创建我们的值的直方图。
		h, err := plotter.NewHist(plotVals,16)
		if err != nil {
			log.Fatal(err)
		}
		// 标准化直方图。
		h.Normalize(1)
		p.Add(h)

		if err := p.Save(4*vg.Inch, 4*vg.Inch, colName+"_hist.png"); err != nil {
			log.Fatal(err)
		}
	}

}
