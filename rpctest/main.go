package main

import (
	"net/rpc"
	"net/http"
	"fmt"
)

// 声明传参类型
type Params struct {
	Width, Height int
}

// 声明矩形对象类型
type Rect struct {

}

// 计算Params的周长
func (r *Rect) Permiter(p Params, ret *int) error {
	*ret = (p.Width + p.Height) * 2
	return nil
}

func main() {
	// 注册服务
	rect := new(Rect)
	rpc.Register(rect)
	rpc.HandleHTTP()
	if err := http.ListenAndServe(":9000",nil); err != nil {
		fmt.Println(err)
	}

}


