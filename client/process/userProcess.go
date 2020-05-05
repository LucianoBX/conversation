package process

import (
	"enconding/binary"
	"enconding/json"
	"fmt"
	"github.com/conversation/client/utils"
	"github.com/conversation/common"
	"net"
	"os"
)

type UserProcess struct{}

// 注册用户
func (up *UserProcess) Register(userId int, userPwd, username string) (err error) {
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err= ", err)
		return
	}

	// 延时关闭
	defer conn.Close()

	//2.send message
	var mes common.Message
	mes.Type = common.RegisterMesType

	// registerMes Object
	var rgMes common.RegisterMes
	rgMes.User.UserId = userId
	rgMes.User.UserPwd = userPwd
	rgmes.User.UserName = userName

	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err= ", err)
		return
	}

	// mes process
	mes.Data = string(data)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("mes json.Marshal err= ", err)
		return
	}

	// 发送消息给服务器
	tf := &utils.Transfer{
		Conn: conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("readPkg(conn) err= ", err)
		return
	}

	// 读取返回包
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) err= ", err)
		return
	}

	var rgResMes common.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功，请登录")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exist(0)
	}
	return
}

// 登录帐号
func (up *UserProcess) Login(userId int, userPwd string) (err error) {
	//
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}

	defer conn.Close()

	// 准备相关结构体
	var mes common.Message
	mes.Type = common.LoginMesType

	var loginMes common.LoginMes
	loginMes.UserId = userId
	loginMes.Userpwd = userPwd

	// 序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err= ", err)
		return
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err= ", err)
	}

	// 先处理长度信息
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:4], pkgLen)

	n, err := conn.Write(buf[:4], pkgLen)
	if n != 4 || err != nil {
		fmt.Println("conn.Write fail err= ", err)
		return
	}

	fmt.Printf("客户端发送长度=%d 内容= %s", len(data), string(data))

	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write mess fail err= ", err)
		return
	}

	// 建立链接，发送消息
	tf := &utils.Transfer{
		Conn: conn,
	}

	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) err= ", err)
		return
	}

	// 将mes的Data部分反序列化程LoginResMes
	var loginResMes common.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResmes)
	if loginResMes.Code == 200 {
		// 初始化CurrUser
		CurrUser.Conn = conn
		CurrUser.UserId = userId
		CurrUser.Userstatus = common.UserOnline

		// 显示当前用户列表
		fmt.Println("当前用户列表如下")
		for _, v := range loginResMes.UserId {
			if v == userId {
				continue
			}

			fmt.Println("userId:\t", v)

			// onlineuser 完成初始化
			user := &common.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")

		// 开一个协程，保持和服务端的通讯，当服务段推送
		// 更新本地用户列表
		go severProcessMes(conn)

		// 显示登录成功的菜单
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}
