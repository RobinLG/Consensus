package main

import (
	"net/rpc"
	"fmt"
)

type Param struct {
	Width, Height int
}

func main() {
	// 调用rpc服务器
	rp, err := rpc.DialHTTP("tcp", "127.0.0.1:9000")
	if err != nil {
		fmt.Println(err)
	}

	ret := 0
	er := rp.Call("Rect.Permiter", Param{100, 100}, &ret)
	if er != nil {
		fmt.Println(er)
	}
	fmt.Println(ret)

}
