package main

import (
	"net/http"
	"os"
	"fmt"
	"io"
)

// 声明节点信息，表示各小国家
type nodeInfo struct {
	// 标志
	id string
	// 准备访问方法
	path string
	// 服务器做出响应
	writer http.ResponseWriter
}

// 存放四个国家的地址
var nodeTable = make(map[string]string)

// 拜占庭容错在Fabric里面的使用
func main() {
	// 获取执行的参数
	userId := os.Args[1]
	fmt.Println(userId)

	// 创建四个国家地址
	nodeTable = map[string]string {
		"a":"localhost:1111",
		"b":"localhost:1112",
		"c":"localhost:1113",
		"d":"localhost:1114",
	}

	node := nodeInfo{userId, nodeTable[userId], nil}

	// http协议的回调函数
	// http://localhost:1111/req?warTime=8888
	http.HandleFunc("/req", node.request)
	http.HandleFunc("/prePrepare", node.prePrepare)
	http.HandleFunc("/prepare", node.prepare)
	http.HandleFunc("/commit", node.commit)


	// 启动服务器
	if err := http.ListenAndServe(node.path, nil); err != nil {

		fmt.Println(err)
	}

}

// 此函数是http访问时候req命令的请求回调函数
func (node *nodeInfo) request(writer http.ResponseWriter, request *http.Request) {
	// 设置允许解析接收参数
	request.ParseForm()
	//fmt.Println(request.Form["warTime"][0])
	if (len(request.Form["warTime"]) > 0) {
		node.writer = writer

		// 激活主节点后，广播给其它节点，通过a向其它节点做广播
		node.broadcast(request.Form["warTime"][0], "/prePrepare")
	}
}

// 向其它节点广播
func (node *nodeInfo) broadcast(msg string, path string) {
	// 遍历所有节点
	for nodeId, url := range nodeTable {
		if nodeId == node.id {
			continue
		}
		//http://localhost:1112/perPrepare?warTime=8888&nodeId=b
		http.Get("http://" + url + path + "?warTime=" + msg + "&nodeId" + node.id)
	}
}

func (node *nodeInfo) prePrepare(writer http.ResponseWriter, request *http.Request) {
	// 设置允许解析接收参数
	request.ParseForm()
	//fmt.Println("HelloWorld")
	// 再做分发
	if len(request.Form["warTime"]) > 0 {
		fmt.Println(node.id)
		node.broadcast(request.Form["warTime"][0], "/prepare")
	}
}

func (node *nodeInfo) prepare(writer http.ResponseWriter, request *http.Request) {
	// 设置允许解析接收参数
	request.ParseForm()
	// 调用验证
	if len(request.Form["warTime"]) > 0 {
		fmt.Println(request.Form["warTime"][0])
	}
	if len(request.Form["nodeId"]) > 0 {
		fmt.Println(request.Form["nodeId"][0])
	}

	node.authentication(request)
}

var authenticationSuccess = true
var authenticationMap = make(map[string] string)
// 获得除了本节点以外的其它节点数据
func (node *nodeInfo) authentication(request *http.Request) {
	// 设置允许解析接收参数
	request.ParseForm()

	if authenticationSuccess != false {
		if len(request.Form["nodeId"]) > 0 {
			authenticationMap[request.Form["nodeId"][0]] = "ok"
		}
	}

	if len(nodeTable) >= len(authenticationMap) * 3 + 1 {
		// 实现拜占庭容错，通过commit反馈给浏览器
		node.broadcast(request.Form["warTime"][0], "/commit")
	}
}

func (node *nodeInfo) commit(writer http.ResponseWriter, request *http.Request) {

	//给浏览器反馈相应
	io.WriteString(node.writer, "ok")
}