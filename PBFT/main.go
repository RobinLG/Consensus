package main

import "fmt"

type Node struct {
	Name   string
	Status int // 以拜占庭将军举例，1为攻，0为不攻
	Votes []*Node
}

// 保存4个node
var nodes = make([]*Node, 0)

func createNodes() {
	A := Node{"A", 1, make([]*Node, 0)}
	B := Node{"B", 1, make([]*Node, 0)}
	C := Node{"C", 1, make([]*Node, 0)}
	D := Node{"D", 0, make([]*Node, 0)}
	// 按照拜占庭协议，这次仗可以打
	nodes = append(nodes, &A)
	nodes = append(nodes, &B)
	nodes = append(nodes, &C)
	nodes = append(nodes, &D)
}

// 互相转达
func Votes() {
	for i :=0; i < len(nodes); i++ {
		// 获取每个将军状态
		fmt.Println(nodes[i].Status)

		// 将当前将军状态分发给其它将军
		for j := 0; j < len(nodes); j++ {
			nodes[j].Votes = append(nodes[j].Votes, nodes[i])
		}
	}
}

// 判断本次进攻是否可行，判断是否满足：全部将军 > 3*作恶将军 + 1
func isValid() bool {
	// 在数组中取出最后一个对象
	node := nodes[len(nodes)-1]
	votes := node.Votes

	cnt := 0 // 作恶将军数
	for _, n := range votes {
		fmt.Println(n.Status)
		if n.Status == 0 {
			cnt++
		}
	}

	// f 为作恶将军
	// 判断是否满足：全部将军 > 3f + 1 。 若不满足则进攻成功 。
	if len(nodes) > 3 * cnt + 1 {
		return false
	}

	return true
}

func main() {
	createNodes()
	Votes()
	fmt.Println(isValid())

}