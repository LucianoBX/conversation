package main

import (
	"fmt"
	"github.com/conversation/common"
	"github.com/conversation/server/model"
	"github.com/conversation/server/process"
	"github.com/conversation/server/utils"
	"github.com/gomodule/redigo/redis"
	"io"
	"net"
	"time"
)

// 定义一个全局的pool
var pool *redis.Pool

//先创建一个Handler 的结构体体
type Handler struct {
	Conn net.Conn
}

func init() {
	//当服务器启动时，我们就去初始化我们的redis的连接池
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()
}

//初始化连接池
func initPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}

//这里我们编写一个函数，完成对UserDao的初始化任务
func initUserDao() {
	//这里的pool 本身就是一个全局的变量
	//这里需要注意一个初始化顺序问题
	//initPool, 在 initUserDao
	model.MyUserDao = model.NewUserDao(pool)
}

// comunicate with clients
func connProcess(conn net.Conn) {
	// delay the close of conn
	defer conn.Close()

	// 调用总控
	handler := &Handler{
		Conn: conn,
	}

	err := handler.sendProcess()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误， err=", err)
	}
}

//编写一个ServerProcessMes 函数
//功能：根据客户端发送消息种类不同，决定调用哪个函数来处理
func (h *Handler) serverProcessMes(mes *common.Message) (err error) {

	//看看是否能接收到客户端发送的群发的消息
	fmt.Println("mes=", mes)

	switch mes.Type {
	case common.LoginMesType:
		//处理登录登录
		//创建一个UserProcess实例
		up := &process.UserProcess{
			Conn: h.Conn,
		}
		err = up.ServerProcessLogin(mes)

	case common.RegisterMesType:
		//处理注册
		up := &process.UserProcess{
			Conn: h.Conn,
		}
		err = up.ServerProcessRegister(mes) // type : data

	case common.SmsMesType:
		//创建一个SmsProcess实例完成转发群聊消息.
		smsProcess := &process.SmsProcess{}
		smsProcess.SendGroupMes(mes)

	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}

// 持续发送消息
func (h *Handler) sendProcess() (err error) {

	//循环的向客户端发送的信息
	for {
		//这里我们将读取数据包，直接封装成一个函数readPkg(), 返回Message, Err
		//创建一个Transfer 实例完成读包任务
		tf := &utils.Transfer{
			Conn: h.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器端也退出..")
				return err
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}

		}

		err = h.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}

}
func main() {

	//提示信息
	fmt.Println("服务器[新的结构]在8889端口监听....")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	//一旦监听成功，就等待客户端来链接服务器
	for {
		fmt.Println("等待客户端来链接服务器.....")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}

		//一旦链接成功，则启动一个协程和客户端保持通讯。。
		go connProcess(conn)
	}
}
