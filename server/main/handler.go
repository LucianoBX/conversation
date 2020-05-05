package main

import (
	"fmt"
	"github.com/coversation/common"
	"github.com/coversation/server/process"
	"github.com/coversation/server/utlis"
	"io"
	"net"
)

//先创建一个Handler 的结构体体
type Handler struct {
	Conn net.Conn
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
		up := &process2.UserProcess{
			Conn: p.Conn,
		}
		err = up.ServerProcessLogin(mes)

	case common.RegisterMesType:
		//处理注册
		up := &process2.UserProcess{
			Conn: p.Conn,
		}
		err = up.ServerProcessRegister(mes) // type : data

	case common.SmsMesType:
		//创建一个SmsProcess实例完成转发群聊消息.
		smsProcess := &process2.SmsProcess{}
		smsProcess.SendGroupMes(mes)

	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}

// 持续发送消息
func (p *Handler) process2() (err error) {

	//循环的向客户端发送的信息
	for {
		//这里我们将读取数据包，直接封装成一个函数readPkg(), 返回Message, Err
		//创建一个Transfer 实例完成读包任务
		tf := &utils.Transfer{
			Conn: p.Conn,
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

		err = p.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}

}
