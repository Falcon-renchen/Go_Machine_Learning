package main

import (
	"log"
	"os"

	"github.com/pachyderm/pachyderm/src/client"
)
//让我们继续将属性文件放入attributes存储库中，并将diabetes.csv训练数据集放入training存储库中
func main() {


	c, err := client.NewFromAddress("0.0.0.0:30650")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	//在“ master”分支上的“ attributes”数据仓库中开始提交。
	commit, err := c.StartCommit("attributes", "master")
	if err != nil {
		log.Fatal(err)
	}

	// 打开属性JSON文件之一。
	f, err := os.Open("1.json")
	if err != nil {
		log.Fatal(err)
	}

	// 将包含属性的文件放入数据存储库。
	if _, err := c.PutFile("attributes", commit.ID, "1.json", f); err != nil {
		log.Fatal(err)
	}

	// 完成提交。
	if err := c.FinishCommit("attributes", commit.ID); err != nil {
		log.Fatal(err)
	}

	// 在“ master”分支的“ training”数据仓库中开始提交。
	commit, err = c.StartCommit("training", "master")
	if err != nil {
		log.Fatal(err)
	}

	// 打开训练数据集。
	f, err = os.Open("diabetes.csv")
	if err != nil {
		log.Fatal(err)
	}

	// 将包含训练数据集的文件放入数据存储库。
	if _, err := c.PutFile("training", commit.ID, "diabetes.csv", f); err != nil {
		log.Fatal(err)
	}

	// 完成提交。
	if err := c.FinishCommit("training", commit.ID); err != nil {
		log.Fatal(err)
	}
}
