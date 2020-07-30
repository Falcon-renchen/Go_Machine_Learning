package main

import (
	"fmt"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

func main() {
	observed := []float64{
		260.0, //此数字是没有定期锻炼的观察到的数字。
		135.0, // 该数字是在零星运动中观察到的数字。
		105.0, // 该数字是定期运动观察到的数字。
	}

	//定义观察到的总数
	totalObserved := 500.0

	//计算预期的频率（再次假设零假设）。
	expected := []float64{
		totalObserved * 0.60,
		totalObserved * 0.25,
		totalObserved * 0.15,
	}

	// 计算ChiSquare检验统计量。
	chiSquare := stat.ChiSquare(observed, expected)

	fmt.Println(chiSquare)

	//创建具有K个自由度的卡方分布。
	//在这种情况下，我们有K = 3-1 = 2，因为卡方分布的自由度是可能类别的数量减去一。
	chiDist := distuv.ChiSquared{
		K: 2.0,
		Src: nil,
	}

	// 计算我们的特定测试统计量的p值。
	pValue := chiDist.Prob(chiSquare)

	//==0。0001    如果大于5%，则需要代替假设
	fmt.Printf("p-value: %0.4f\n\n", pValue)
}
