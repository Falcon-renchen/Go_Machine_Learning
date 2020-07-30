package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

/*
如果您的CSV文件没有用逗号分隔，并且/或者如果CSV文件包含带注释的行，
则可以利用csv.Reader.Comma和csv.Reader.Comment字段来正确处理格式独特的CSV文件。
如果CSV文件中的字段用单引号引起来，则可能需要添加一个辅助函数来修剪单引号并解析值。
 */

func main() {
	//打开文件
	f, err := os.Open("dataform/iris.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)

	//全部读取
	reader.FieldsPerRecord = -1 	//假如我们不知道每一行有多少个数据，设置-1

	//如果要打印输循环输出一条数据，则把以下信息替换成后面的code
	//将所有数据都打印出来
	rawCSVData, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	fmt.Println(rawCSVData)

	/*//循环打印一条数据
	var rawCSVData2 [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		rawCSVData = append(rawCSVData, record)
	}
	fmt.Println(rawCSVData2)
	*/

}
