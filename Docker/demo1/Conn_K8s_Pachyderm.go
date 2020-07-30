package main

import (
	"log"

	"github.com/pachyderm/pachyderm/src/client"
)

func main() {
	//连接Pachyderm

	//使用Kubernetes集群的IP连接到Pachyderm。当您在本地运行k8和/或将Pachyderm端口转发到本地主机时，
	//在这里我们将使用localhost来模拟sceneario。默认情况下，Pachyderm将在端口30650上公开。
	c, err := client.NewFromAddress("0.0.0.0:30650")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	//创建一个名为“培训”的数据存储库。
	if err := c.CreateRepo("training"); err != nil {
		log.Fatal(err)
	}

	//创建一个名为“属性”的数据存储库。
	if err := c.CreateRepo("attributes"); err != nil {
		log.Fatal(err)
	}

	//现在，我们将列出Pachyderm集群上的所有当前数据存储库，以进行完整性检查。
	//现在，我们应该有两个数据存储库。
	repos, err := c.ListRepo()
	if err != nil {
		log.Fatal(err)
	}

	//检查回购的数量是否符合我们的预期。
	if len(repos) != 2 {
		log.Fatal("Unexpected number of data repositories")
	}

	//检查存储库的名称是否符合我们的期望。
	if repos[0].Repo.Name != "attributes" || repos[1].Repo.Name != "training" {
		log.Fatal("Unexpected data repository name")
	}
}
