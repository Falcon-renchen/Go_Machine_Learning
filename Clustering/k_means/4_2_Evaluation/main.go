package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gonum/floats"
	"github.com/kniren/gota/dataframe"
)

//为了更定量地评估聚类，我们可以计算聚类中的点与聚类质心之间的聚类内平均距离
func main() {
	f, err := os.Open("fleet_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	driverDF := dataframe.ReadCSV(f)

	distances := driverDF.Col("Distance_Feature").Float()

	// clusterOne和clusterTwo将保存用于绘制的值。
	var clusterOne [][]float64
	var clusterTwo [][]float64

	// Fill the clusters with data.
	for i, speed := range driverDF.Col("Speeding_Feature").Float() {
		distanceOne := floats.Distance([]float64{distances[i], speed}, []float64{50.05, 8.83}, 2)
		distanceTwo := floats.Distance([]float64{distances[i], speed}, []float64{180.02, 18.29}, 2)
		if distanceOne < distanceTwo {
			clusterOne = append(clusterOne, []float64{distances[i], speed})
			continue
		}
		clusterTwo = append(clusterTwo, []float64{distances[i], speed})
	}

	// 输出我们的集群内指标
	fmt.Printf("\nCluster 1 Metric: %0.2f\n", withinClusterMean(clusterOne, []float64{50.05, 8.83}))
	fmt.Printf("\nCluster 2 Metric: %0.2f\n", withinClusterMean(clusterTwo, []float64{180.02, 18.29}))

}


// insideClusterMean计算之间的平均距离
//指向聚类和聚类的质心。
func withinClusterMean(cluster [][]float64, centroid []float64) float64 {

	//
	/// meanDistance将保留我们的结果。
	var meanDistance float64

	//
	//遍历群集中的点。
	for _, point := range cluster {
		meanDistance += floats.Distance(point, centroid, 2) / float64(len(cluster))
	}

	return meanDistance
}
