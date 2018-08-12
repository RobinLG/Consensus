package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"net"
	"bufio"
	"fmt"
	"io"
	"strconv"
	"math/rand"
)

// 创建区块
type Block struct {
	Height int
	Timestamp int
	PrevHash string
	Hash string
	Data string
	// 终端地址
	Validator string
}


// 计算某个字符串的hash
func calculateHash(record string) string {
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

// 计算block的hash
func calculateBlockHash(block Block) string {
	record := string(block.Height) + string(block.Timestamp) + block.PrevHash + block.Data
	hashCode := calculateHash(record)
	return hashCode
}

// 生成新区块链
var Blockchain []Block

func generateNextBlock(oldBlock Block, data string, vald string) Block {
	var newBlock Block
	// 设置区块高度
	newBlock.Height = oldBlock.Height + 1
	newBlock.Timestamp = int(time.Now().Unix())
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Data = data
	newBlock.Hash = calculateBlockHash(newBlock)
	newBlock.Validator = vald
	// 添加到区块链
	Blockchain = append(Blockchain, newBlock)
	return  newBlock
}

// 创建创世区块
func genesisiBlock() Block {
	var genesisBlock = Block{0, int(time.Now().Unix()), "", "", "", ""}
	// 计算genesisBlock的hash值
	genesisBlock.Hash = calculateBlockHash(genesisBlock)
	// 将创世区块添加到数组
	Blockchain = append(Blockchain, genesisBlock)
	return genesisBlock
}

// 创建conn终端链接的数组
var connAddr []net.Conn

// 创建节点类型
type Node struct {
	// 终端地址
	Address string
	// 币龄
	Coins int
}

// 保存终端的对象
var nodes []Node

// 通道实现线程通信
var announcements = make(chan string)

func main() {

	// 测试代码
	// genesisBlock := genesisiBlock()
	// newBlock := generateNextBlock(genesisBlock, "1")
	// generateNextBlock(newBlock, "2")
	// 区块链中是否有三个区块
	// fmt.Println(Blockchain)

	genesisBlock := genesisiBlock()

	// 通过终端连接上代码
	// 终端telnet
	// 创建监听
	netListen, _ := net.Listen("tcp", "127.0.0.1:1234")
	defer netListen.Close()

	go func() {
		// 通过矿工实现区块的挖矿
		for {
			// 此代码须等announcements有值后才会执行
			w := <- announcements
			// 将新的区块，利用w矿工，添加到数组中
			block := generateNextBlock(genesisBlock, "100", w)
			fmt.Println(block)
			// UTXO

		}
	}()

	go func() {
		// 每隔开10s选择一次矿工
		for {
			time.Sleep(10 * time.Second)
			winner := pickWinner()
			fmt.Println("系统通过PoS帮您选出的矿工为", winner)
			// 将矿工放入到通信中
			announcements <- winner
		}
	}()

	// 等待链接
	for {
		conn, _ := netListen.Accept()
		// 将所有的链接保存到数组
		connAddr = append(connAddr, conn)
		// 创建缓存区
		scanbalance := bufio.NewScanner(conn)

		io.WriteString(conn, "Please input Coins:")

		go func() {
			// 在协程中扫描终端内容
			for scanbalance.Scan() {
				txt := scanbalance.Text()
				// 打印终端输入的信息
				fmt.Println("您刚才从终端输入的币龄为：", txt)

				// 通过时间戳创建地址
				addr := calculateHash(time.Now().String())
				cons, _ := strconv.Atoi(txt)
				node := Node{addr, cons}
				// 将链接终端对象存放到数组
				nodes = append(nodes, node)
				fmt.Println(nodes)
			}
		}()

	}
}

// 通过PoS共识算法选择矿工
func pickWinner() string {
	// 选择矿工，利用PoS共识算法选择矿工
	var lottyPool []string

	// 根据币龄把对应的矿工地址存放到数组中
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		for j := 0; j < node.Coins; j++ {
			lottyPool = append(lottyPool, node.Address)
		}
	}

	if len(lottyPool) != 0 {
		// 通过随机值，找到准备挖矿的矿工
		rand.Seed(time.Now().Unix())
		r := rand.Intn(len(lottyPool))
		workerAddress := lottyPool[r]
		// 返回矿工地址
		return workerAddress
	}

	return ""

}

