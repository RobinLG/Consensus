package BLC

import (
	"math/big"
	"bytes"
	"encoding/binary"
	"log"
	"crypto/sha256"
	"fmt"
)

// 难度位
const targetBit = 20

type ProofOfWork struct {
	Block *Block // 当前要验证的区块
	target *big.Int // 大数据存储，区块目标难度
}

// 将int64转换成字节数组
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}


func (pow *ProofOfWork) perpareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(targetBit)),
			IntToHex(int64(nonce)),
			IntToHex(int64(pow.Block.Height)),
		},
		[]byte{},
	)

	return data
}

// 创建新工作量证明对象
func NewProofOfWork(block *Block) *ProofOfWork {
	// 创建一个初始值为1的target（比特币系统中计算出此target相对复杂，此处简化，直接取1）
	target := big.NewInt(1)

	// 左移256 - targetBit位
	target = target.Lsh(target, 256 - targetBit)

	return &ProofOfWork{block, target}
}

func (proofOfWork *ProofOfWork) run() ([]byte, int64) {
	// 1.拼接block属性
	// 2.生成Hash
	// 3.判断Hash有效性

	nonce := 0
	var hashInt big.Int //存储有效Hash
	var hash [32]byte

	for {
		//准备数据
		dataBytes := proofOfWork.perpareData(nonce)

		//生成Hash
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])

		// 判断hashInt是否小于Block里面的target
		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//   proofOfWork.target == x
		//	 hashInt == y
		if proofOfWork.target.Cmp(&hashInt) == 1 {
			break
		}

		nonce = nonce + 1
	}

	return hash[:], int64(nonce)
}