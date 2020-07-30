package main
/**
第一列Driver_ID包括特定驱动程序的各种匿名标识。
第二和第三列是我们将在集群中使用的属性。
该Distance_Feature列是每个数据驱动的平均距离，
并且Speeding_Feature是驾驶员以每小时超过速度限制的速度行驶5英里以上的平均时间百分比。
 */
import (
	"github.com/go-gota/gota/dataframe"
	"log"
	"os"
)
func main() {
	f, err := os.Open("fleet_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fleetDF := dataframe.ReadCSV(f)

	fleetDF.Describe()

}
