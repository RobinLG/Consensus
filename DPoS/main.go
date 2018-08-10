package main

import (
	"robin/DPoS/BLC"
	"fmt"
)



func main() {
	BLC.CreateNode()

	var block = BLC.NewBlock(1,[]byte{0},nil)

	fmt.Printf("%x---%s", block.Hash, block.Delegate.Name)

}
