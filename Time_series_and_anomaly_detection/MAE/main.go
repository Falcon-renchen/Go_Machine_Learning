package main

import (
	"fmt"
	"log"

	"github.com/lytics/anomalyzer"
)

//我们可能并不总是对预测时间序列感兴趣。我们可能想要检测时间序列中的异常行为。
//例如，我们可能想知道何时网络上出现异常流量，或者当异常用户尝试在应用程序中进行某些操作时，我们可能想要发出警报。这些事件可能与安全性相关，或者仅用于调整我们的基础结构或应用程序设置。


//首先，InfluxDB（https:www.influxdata.com/）和Prometheus（https://prometheus.io/）生态系统具有多种异常检测选项。
// InfluxDB和Prometheus均提供基于Go的开源时间序列数据库和相关工具。它们对监视基础结构和应用程序很有用
func main() {

	//使用配置（例如我们要使用的异常检测方法）初始化AnomalyzerConf值。
	conf := &anomalyzer.AnomalyzerConf{
		Sensitivity: 0.1,
		UpperBound:  5,
		LowerBound:  anomalyzer.NA, // 忽略下界
		ActiveSize:  1,
		NSeasons:    4,
		Methods:     []string{"diff", "fence", "highrank", "lowrank", "magnitude"},
	}


	//将定期观测的时间序列创建为浮点数切片。如先前示例中所使用的，这可能来自数据库或文件。
	ts := []float64{0.1, 0.2, 0.5, 0.12, 0.38, 0.9, 0.74}


	//根据现有的时间序列值和配置创建一个新的异常分析器。
	anom, err := anomalyzer.NewAnomalyzer(conf, ts)
	if err != nil {
		log.Fatal(err)
	}


	//向异常分析器提供新的观测值。
	//异常分析器将参考序列中的现有值分析该值，并输出该值异常的可能性。
	prob := anom.Push(15.2)
	fmt.Printf("Probability of 15.2 being anomalous: %0.2f\n", prob)

	prob = anom.Push(0.43)
	fmt.Printf("Probability of 0.33 being anomalous: %0.2f\n", prob)
}
