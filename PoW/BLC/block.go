package BLC

import (
	"time"
	"fmt"
)

type Block struct {
	// 区块高度
	Height int64
	// 上一区块哈希
	PrevHash []byte
	// 区块数据
	Data []byte
	// 时间戳
	Timestamp int64
	// Hash
	Hash []byte
	// Nonce
	Nonce int64
}

func NewBlock(height int64, prevHash []byte, data []byte) *Block {

	// 创建新区块
	block := &Block{height, prevHash, data, time.Now().Unix(), nil, 0}

	//调用工作量证明返回有效Hash和Nonce
	pow := NewProofOfWork(block)

	//挖矿验证
	hash, nonce := pow.run()

	block.Hash = hash
	block.Nonce = nonce

	fmt.Println()

	return block
}


