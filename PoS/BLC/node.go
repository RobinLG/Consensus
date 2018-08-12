package BLC

type Node struct {
	Tokens int // 持币数
	Days int // 币龄
	Address []byte // 地址
}

// 创建五个节点
var n = make([]Node, 5)
// 存储每个节点，以此模仿选中节点的概率事件
var addr = make([]*Node, 15)

func InitNode() {
	n[0] = Node{1, 1, []byte("0x1")}
	n[1] = Node{2, 1, []byte("0x2")}
	n[2] = Node{3, 1, []byte("0x3")}
	n[3] = Node{4, 1, []byte("0x4")}
	n[4] = Node{5, 1, []byte("0x5")}
}

