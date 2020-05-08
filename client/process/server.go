package process

import (
	"encoding/json"
	"fmt"
	"github.com/conversation/client/utils"
	"github.com/conversation/common"
	"net"
	"os"
)

// 显示登录后的界面
func ShowMenu() {
	fmt.Println("-------恭喜xxx登录成功---------")
	fmt.Println("-------1. 显示在线用户列表---------")
	fmt.Println("-------2. 发送消息---------")
	fmt.Println("-------3. 信息列表---------")
	fmt.Println("-------4. 退出系统---------")
	fmt.Println("请选择(1-4):")
	var (
		key     int
		content string
	)

	// 将SmsProcess建立在外部，优化
	smsProcess := &SmsProcess{}
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		//显示在线用户列表
		outputOnlineUser()
	case 2:
		fmt.Println("说些啥？/What U want to say?")
		fmt.Scanf("%s/n", &content)
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("你选择推出系统……")
		os.Exit(0)
	default:
		fmt.Println("你输入的选项不正确/wrong input")
	}
}

func serverProcessMes(conn net.Conn) {
	// 创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}

	for {
		fmt.Println("客户端正等待读取服务器消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err= ", err)
			return
		}

		// 读取到信息下一步处理
		switch mes.Type {
		case common.NotifyUserStatusMesType: // 有人上线了
			//1. 取出.NotifyUserStatusMes
			var notifyUserStatusMes common.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			//2. 把这个用户的信息，状态保存到客户map[int]User中
			updateUserStatus(&notifyUserStatusMes)

		//处理
		case common.SmsMesType: //有人群发消息
			outputGroupMes(&mes)
		default:
			fmt.Println("服务器端返回了未知的消息类型")
		}
	}
}
