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
	// 代理人
	Delegate *Node
}

func NewBlock(height int64, pervhash []byte, data []byte)  *Block {

	// 创建DPoS对象
	dpos := NewPos(1, []byte{0}, nil)

	// 先选出代理矿工再创建新区块
	// 采用DPoS挖矿
	block := dpos.generateNextBlock()

	return block
}
