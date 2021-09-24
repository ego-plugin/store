package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gotomicro/ego/core/econf"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ego-plugin/store/eqmgo"
)

type UserInfo struct {
	Name   string `bson:"name"`
	Age    uint16 `bson:"age"`
	Weight uint32 `bson:"weight"`
}

var userInfo = UserInfo{
	Name: "xm",
	Age: 7,
	Weight: 40,
}

func main() {
	var stopCh = make(chan bool)
	// 假设你配置的toml如下所示
	conf := `
[mongo]
	debug=true
	dsn="mongodb://user:password@localhost:27017,localhost:27018"
`
	// 加载配置文件
	err := econf.LoadFromReader(strings.NewReader(conf), toml.Unmarshal)
	if err != nil {
		panic("LoadFromReader fail," + err.Error())
	}

	// 初始化emongo组件
	cmp := eqmgo.Load("mongo").Build()
	coll := cmp.Client().Database("test").Collection("cells")
	findOne(coll)

	stopCh <- true
}

func findOne(coll *eqmgo.Collection) {
	one := UserInfo{}
	err := coll.Find(context.TODO(), bson.M{"name": userInfo.Name}).One(&one)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(one)
}