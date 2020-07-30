package main

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

//对于经常访问的数据，就要设置缓存，暂时缓存到本地，或在应用程序运行的时候可以进行访问。

//例如要获得政府的api，用于人口普查，由于访问会有高延迟，所以需要进行缓存


func main() {

	//创建一个默认过期时间为5分钟的缓存，其中
	//每30秒清除一次过期的项目
	c := cache.New(5*time.Minute, 30*time.Second)

	////将键和值放入缓存中
	c.Set("mykey","myvalue",cache.DefaultExpiration)

	//进行健全性检查。 在缓存中输出键和值以标准化。
	v, found := c.Get("mykey")
	if found {
		fmt.Printf("key: mykey, value: %s\n", v)
	}
}
