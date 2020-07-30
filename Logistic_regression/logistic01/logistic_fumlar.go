package main
//逻辑函数：
//    f(x) = 1 / (1+e^-x)
import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
	"math"
)

func logistic(x float64) float64 {
	return 1 / (1+math.Exp(-x))
}

func main() {
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.Title.Text = "Logistic Regression"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "f(x)"

	//创建绘图仪函数。
	logisticPlotter := plotter.NewFunction(func(x float64) float64 {
		return logistic(x)
	})
	logisticPlotter.Color = color.RGBA{B:255,A:255}

	p.Add(logisticPlotter)


	//设置轴范围。与其他数据集不同，由于函数不一定具有x和y值的有限范围，因此函数不会自动设置轴范围。
	p.X.Min = -10
	p.X.Max = 10
	p.Y.Min = -0.1
	p.Y.Max = 1.1

	if err := p.Save(4*vg.Inch,4*vg.Inch,"logistic.png"); err!=nil {
		log.Fatal(err)
	}
	
}
