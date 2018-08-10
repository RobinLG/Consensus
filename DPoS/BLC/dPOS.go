package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

type DPoS struct {
	Block *Block // 当前要验证区块
	Node *Node // 代理矿工
}

func (dpos *DPoS) perpareData() []byte {
	data := bytes.Join(
		[][]byte{
			dpos.Block.Data,
			dpos.Block.PrevHash,
			IntToHex(dpos.Block.Height),
			IntToHex(dpos.Block.Timestamp),
		},
		[]byte{},
	)

	return data
}

// 计算区块Hash
func (dpos *DPoS) CalculateHash() []byte {
	//准备数据
	record := dpos.perpareData()

	//生成Hash
	hash := sha256.Sum256(record)
	//fmt.Printf("\r%x", hash)

	return hash[:]
}

// 创建DPoS对象
func NewPos(height int64, prevHash []byte, data []byte) *DPoS {

	dpos := new(DPoS)
	fmt.Println(dpos.sortNodes())// 选出的三个矿工

	// DPoS中选出票数最高的前n位，作为代理矿工,模拟轮流选出矿工，仅做示例展示。实际上每次生产块都是一个循环到节点[0]
	dpos.Node = dpos.sortNodes()[0]
	dpos.Node = dpos.sortNodes()[1]
	dpos.Node = dpos.sortNodes()[2]

	// 创建区块对象
	dpos.Block = &Block{height, prevHash, time.Now().Unix(), data, nil, nil}

	return &DPoS{dpos.Block, dpos.Node}
}

// DPoS中选出票数最高的前n位
func (dpos *DPoS) sortNodes() []*Node {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4-i; j++ {
			if n[j].Votes < n[j+1].Votes {
				t := n[j]
				n[j] = n[j+1]
				n[j+1] = t
			}
		}
	}

	return n[:3]
}

// 设置代理人方法
func (dpos *DPoS) setDelete (node *Node) {
	dpos.Block.Delegate = node
}

func (dpos *DPoS) generateNextBlock() *Block {

	// 设置当前区块挖矿
	dpos.Block.Delegate = dpos.Node

	// 设置当前区块挖矿
	dpos.Block.Hash = dpos.CalculateHash()

	return dpos.Block
}