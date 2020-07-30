package main

import (
	"fmt"
	"github.com/kniren/gota/dataframe"
	"gonum.org/v1/gonum/floats"
	"log"
	"os"
)
//dfFloatRow从DataFrame检索一个float值切片
//在给定的索引和给定的列名。
func dfFloatRow(df dataframe.DataFrame, names []string, idx int) []float64 {
	var row []float64
	for _, name := range names {
		row = append(row, df.Col(name).Float()[idx])
	}
	return row
}

type centroid []float64
func main() {
	irisFile, err := os.Open("iris.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer irisFile.Close()
	irisDF := dataframe.ReadCSV(irisFile)

	//定义CSV文件中包含的三个独立物种的名称。
	speciesName := []string{
		"Iris-setosa",
		"Iris-versicolor",
		"Iris-virginica",
	}
	//创建一个地图以保存我们的质心信息
	centroids := make(map[string]centroid)
	// 创建一个映射，以保存每个集群的过滤后的数据框。
	clusters := make(map[string]dataframe.DataFrame)

	//将数据集过滤到三个单独的数据帧中，每个数据帧都对应一个物种
	for _, species := range speciesName {
		//过滤原始数据集。
		filter := dataframe.F{
			Colname:    "species",
			Comparator: "==",
			Comparando: species,
		}
		filtered := irisDF.Filter(filter)

		////将过滤后的数据框添加到群集图。
		clusters[species] = filtered

		//计算特征的均值。
		summaryDF := filtered.Describe()

		//将每个维度的均值放入相应的质心。
		var c centroid
		for _, feature := range summaryDF.Names() {
			//跳过不相关的列。
			if feature=="column" || feature=="species" {
				continue
			}
			c = append(c, summaryDF.Col(feature).Float()[0])
		}
		//将此质心添加到我们的地图。
		centroids[species] = c
		////clusters[species] = filtered
		////作为健全性检查，输出我们的质心。
		//for _, species := range speciesName {
		//	fmt.Printf("%s centroid: %v\n", species, centroids[species])
		//}
	}
	//将我们的标签转换成字符串切片并创建一个切片
	//为方便起见，使用浮点列名称。
	labels := irisDF.Col("species").Records()
	floatColumns := []string{
		"sepal_length",
		"sepal_width",
		"petal_length",
		"petal_width",
	}

	//循环记录累积平均轮廓系数的记录。
	var silhouette float64

	for idx, label := range labels {

		// a将存储a的累计值。
		var a float64

		// 循环遍历同一群集中的数据点。
		for i := 0; i < clusters[label].Nrow(); i++ {

			// 获取数据点进行比较。
			current := dfFloatRow(irisDF, floatColumns, idx)
			other := dfFloatRow(clusters[label], floatColumns, i)

			// Add to a.
			a += floats.Distance(current, other, 2) / float64(clusters[label].Nrow())
		}

		// 确定最近的其他群集。
		var otherCluster string
		var distanceToCluster float64
		for _, species := range speciesName {

			// 跳过包含数据点的群集。
			if species == label {
				continue
			}

			// 计算到当前群集到群集的距离。
			distanceForThisCluster := floats.Distance(centroids[label], centroids[species], 2)

			// 如果相关，请替换当前集群。
			if distanceToCluster == 0.0 || distanceForThisCluster < distanceToCluster {
				otherCluster = species
				distanceToCluster = distanceForThisCluster
			}
		}

		// b将存储b的累计值。
		var b float64

		// 循环访问最近的其他群集中的数据点。
		for i := 0; i < clusters[otherCluster].Nrow(); i++ {

			// 获取数据点进行比较。
			current := dfFloatRow(irisDF, floatColumns, idx)
			other := dfFloatRow(clusters[otherCluster], floatColumns, i)

			// Add to b.
			b += floats.Distance(current, other, 2) / float64(clusters[otherCluster].Nrow())
		}

		// 添加到平均轮廓系数。
		if a > b {
			silhouette += ((b - a) / a) / float64(len(labels))
		}
		silhouette += ((b - a) / b) / float64(len(labels))
	}

	// 将最终的平均轮廓系数输出到stdout。
	fmt.Printf("\nAverage Silhouette Coefficient: %0.2f\n\n", silhouette)

}
