package BLC

type Block struct {
	// 区块高度
	Height int64
	// 上一区块Hash
	PrevHash []byte
	// 时间戳
	Timestamp int64
	// 交易数据
	Data []byte
	// 当前区块Hash
	Hash []byte
	// 记录挖矿节点地址
	Validator *Node
}

//func CreateGenesisBlock(data []byte) *Block {
//
//	block := &Block{1, []byte{0}, time.Now().Unix(), data, nil, nil}
//
//	pos := NewPos(block)
//	block.Hash = pos.CalculateHash()
//	return block
//}

func NewBlock(height int64, prevHash []byte, data []byte) *Block{

	// 创建PoS对象
	pow := NewPos(height, prevHash, data)

	// 先选出矿工再创建新区块
	// 采用PoS挖矿
	block := pow.generateNextBlock()

	return block
	}


