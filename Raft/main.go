package main

/*
   因为fmt.print是io的操作，io操作总会有等待的过程，所以goroutine会与fmt.print切换，意味着显示至控制台的信息有时候并不意味着实际执行顺序
 */

import (
	"fmt"
	"math/rand"
	"flag"
	"time"
	"net"
	"strconv"
	"log"
	"strings"
	"net/http"
)

const(
	LEADER = iota
	CANDIDATE
	FOLLOWER
)

type Addr struct {
	Host string // IP
	Port int // Port
	Addr string // 地址描述
}

type RaftSever struct {
	// 选票
	Votes int
	// 角色 : follower,candidate,leader
	Role int
	// 存放每个节点的地址信息
	Node []Addr
	// 判断当前节点是否处于选举中
	isElecting bool
	// 设置选举的间隔时间(超时时间 )
	Timeout int
	// 设置两个通道变量
	ElectChan chan bool
	// 控制leader心跳信号
	HeartBeatChan chan bool
	// 实现网页接收数据的传参
	CustomerMsg chan string
	// 服务器端口
	Port int
}

func (rs *RaftSever) changeRole(role int) {
	switch role {
	case LEADER:
		fmt.Println("I become Leader")
	case CANDIDATE:
		fmt.Println("I became candidate")
	case FOLLOWER:
		fmt.Println("I become follower")
	}
	rs.Role = role
}

func (rs *RaftSever) resetTIMEOUT() int {
	rand.Seed(time.Now().Unix())
	// Raft系统中一般为1500 - 3000 ms 选一次
	return 1500 + rand.Intn(1500)
}

func main() {

	// 获取参数的方法
	port := flag.Int("p", 5000, "port")
	flag.Parse()
	fmt.Println("输入的端口号为：", *port)

	// 创建新对象
	rs := RaftSever{}

	// 监听http协议
	go rs.setHttpServer()

	rs.Votes = 0
	rs.Role = FOLLOWER
	// 默认可以投票
	rs.isElecting = true
	rs.Timeout = rs.resetTIMEOUT()
	rs.ElectChan = make(chan bool)
	rs.HeartBeatChan = make(chan bool)
	rs.CustomerMsg = make(chan string)
	rs.Node = []Addr{
		{"127.0.0.1", 5000, "5000"},
		{"127.0.0.1", 5001, "5001"},
	}

	// 设置服务器端口
	rs.Port = *port
	// 运行rs
	rs.Run()

	//for {;}
}

// 运行服务器程序
func (rs *RaftSever) Run() {
	// rs监听,通过tcp协议
	netListen, _ := net.Listen("tcp", ":" + strconv.Itoa(rs.Port))

	// 给其它节点发送数据
	go rs.elect()
	// 控制发送数据的间隔时间
	go rs.electTimeDuration()
	// 查状态
	go rs.printRole()
	// leader发送心跳信号
	go rs.sendHeartBeat()
	// leader给其它服务器发送信息
	go rs.sendDataToOtherNodes()

	for {
		// 等待其它服务器的链接
		conn, _ := netListen.Accept()
		// 监听服务器信息
		go func() {
			for  {
				bts := make([]byte, 1024)
				n, _ := conn.Read(bts)
				fmt.Println("收到的消息: ", string(bts[:n]))
				if strings.HasPrefix(string(bts[:n]), "agree") {
					// 说明这是投票数据，有服务器给此服务器投票
					rs.Votes++
					data := "当前服务器: " + strconv.Itoa(rs.Port) + ", 当前票数为: " + strconv.Itoa(rs.Votes)
					fmt.Println(data)
					// 判断票数等于指定的值，则leader选择成功
					if VotesSuccess(rs.Votes, len(rs.Node)/2) {
						msg := "服务器" + strconv.Itoa(rs.Port) + "被选举为leader"
						fmt.Println(msg)
						// 通知其它服务器，停止选举工作，并且其它节点退回follower状态
						rs.writeToOthers("stopVotes")
						rs.isElecting = false
						rs.changeRole(LEADER)


					}
				}

				if strings.HasPrefix(string(bts[:n]), "stopVotes") {

					rs.isElecting = false
					rs.changeRole(FOLLOWER)
				}
			}
		}()

	}
}

func VotesSuccess(votes int, target int) bool {
	if votes == target {
		return true
	}
	return false
}

func (rs *RaftSever) writeToOthers(data string) {
	// 向其它服务器发送数据
	for _, v := range rs.Node {
		if v.Port != rs.Port {
			netAddr, err := net.ResolveTCPAddr("tcp4", ":" + strconv.Itoa(v.Port))
			if err != nil {
				log.Panic(err)
			}
			//fmt.Println(netAddr)
			conn, err := net.DialTCP("tcp", nil, netAddr)
			if err != nil {
				fmt.Println("无其它服务器开启")
			}else{
				data = data + " 发送的服务器为: " + strconv.Itoa(rs.Port)
				conn.Write([]byte(data))
			}

		}
	}
}

// 给别的服务器发送数据
func (rs *RaftSever) elect() {
	for {
		// 通过通道，确定现在可以投票，才给别人投票
		<- rs.ElectChan
		rs.writeToOthers("agree")
		// rs.ElectChan <- false
		if rs.Role != LEADER {
			// 服务器变为candidate状态
			rs.changeRole(CANDIDATE)
		}else{

		}


	}
}

// 选举时间间隔
func (rs *RaftSever) electTimeDuration() {
	for {
		if rs.isElecting{
			time.Sleep(time.Duration(rs.Timeout) * time.Millisecond)
			rs.ElectChan <- true
		}
	}
}

// 各服务器状态
func (rs *RaftSever) printRole() {
	for {
		time.Sleep(time.Second)
		fmt.Println(strconv.Itoa(rs.Port) + " 状态为 " + strconv.Itoa(rs.Role))
	}
}

// leader发送心跳信号
func (rs *RaftSever) sendHeartBeat() {
	// 每隔一秒发送一次心跳
	for {
		time.Sleep(time.Second)
		if rs.Role == LEADER {
			// 发送消息
			rs.writeToOthers("heart beating....")
		}
	}
}

// 通过leader给其它服务器发信息
func (rs *RaftSever) sendDataToOtherNodes() {
	for {
		// 每隔十秒由主节点给子节点发送数据
		//time.Sleep(10 * time.Second)
		msg := <- rs.CustomerMsg
		if rs.Role == LEADER {
			rs.writeToOthers(msg)
		}
	}
}

// 设置服务器
func (rs *RaftSever) setHttpServer() {
	http.HandleFunc("/req", rs.request)
	if err := http.ListenAndServe("127.0.0.1:1234", nil); err != nil {
		fmt.Println(err)
	}
}

func (rs *RaftSever) request (writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	// 接收网页信息
	if len(request.Form["data"]) > 0 {
		fmt.Println(request.Form["data"][0])
		rs.CustomerMsg <- request.Form["data"][0]
	}
}