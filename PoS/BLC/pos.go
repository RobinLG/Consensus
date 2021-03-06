package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

type PoS struct {
	Block *Block // 当前要验证的区块
	Node *Node //矿工
}

func (pos *PoS) perpareData() []byte {
	data := bytes.Join(
		[][]byte{
			pos.Block.PrevHash,
			pos.Block.Data,
			IntToHex(pos.Block.Timestamp),
			IntToHex(pos.Block.Height),
		},
		[]byte{},
	)

	return data
}

// 计算区块Hash
func (pos *PoS) CalculateHash() []byte {

	//准备数据
	record := pos.perpareData()

	//生成Hash
	hash := sha256.Sum256(record)
	fmt.Printf("\r%x",hash)
	//pos.Block.Hash = hash[:]
	return  hash[:]
}

// 创建Pos对象
func NewPos(height int64, prevHash []byte, data []byte) *PoS {

	pos := new(	PoS)

	// 创建五个矿工节点
	InitNode()
	// 将各节点地址加入数组，模拟选择矿工的概率
	cnt := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < (n[i].Tokens * n[i].Days); j++ {
			addr[cnt] = &n[i]
			cnt++
		}
	}

	// 选出挖矿节点

	// 设置随机种子
	rand.Seed(time.Now().Unix())
	//[0,15)随机数
	var rd = rand.Intn(15)
	//选出矿工
	pos.Node = addr[rd]

	//创建区块对象
	pos.Block = &Block{height, prevHash, time.Now().Unix(), data, nil, nil}


	return &PoS{pos.Block, pos.Node}
}

func (pos *PoS) generateNextBlock() *Block {


	// 设置当前区块挖矿
	pos.Block.Validator = pos.Node

	// 计算区块Hash
	pos.Block.Hash = pos.CalculateHash()

	// 挖矿节点矿工原有的币增加(此处注意！由于币数增加，成为矿工节点的概率也会改变，node文件下的[]*Node不再是15，币龄也会归零，所涉情况较复杂，此简化版不再优化)
	pos.Block.Validator.Tokens = pos.Block.Validator.Tokens + 1

	return pos.Block

}


