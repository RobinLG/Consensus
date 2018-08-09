package main

import (
	"robin/PoS/BLC"
	"fmt"
)

func main() {

	BLC.InitNode()

	var block = BLC.NewBlock(1,[]byte{0},nil)

	fmt.Printf("\n%s\n", block.Validator.Address)
	fmt.Println(block.Validator.Tokens)

}
