package main

import "robin/PoW/BLC"

func main() {
	//height int64, prevHash []byte, data []byte
	BLC.NewBlock(1, []byte{}, nil)
}

