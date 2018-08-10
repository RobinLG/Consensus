package BLC

import (
	"math/rand"
	"time"
)

type Node struct {
	Name string // 节点名字
	Votes int // 被选举票数
}

// 创建数组保存节点
var n = make([]*Node, 5)

// 创建节点
func CreateNode() {

	//创建随机种子
	rand.Seed(time.Now().Unix())
	//随机票数
	n[0] = &Node{"node1", rand.Intn(10)}
	n[1] = &Node{"node2", rand.Intn(10)}
	n[2] = &Node{"node3", rand.Intn(10)}
	n[3] = &Node{"node4", rand.Intn(10)}
	n[4] = &Node{"node5", rand.Intn(10)}

}


