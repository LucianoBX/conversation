package process

import (
	"fmt"
)

type UserProcess struct {
	Conn   net.Conn
	UserId int
}

// 通知所有在线用户
func (up *UserProcess) NotifyOthersOnlineUser(userId int) {

	// 遍历onlineUsers，然后一个一个的发送NotifyUserStatusMes
	for id, up := range userMgr.onlineUsers {
		// 跳过自身
		if id == userId {
			continue
		}

		// 通知
		up.NotifyMeOnline(userId)
	}
}

// 群发上线消息
func (up *UserProcess) NotifyMeOnline(userId int) {

	//
	var mes common.Message
	mes.Type = common.NotifyUserStatusMesType

	var notifyUserStatusMes common.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = common.UserOnline

	//序列化消息
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marhal err=", err)
		return
	}

	// 发送消息
	tf := &utils.Transfer{
		Conn: up.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err=", err)
		return
	}
}

// 注册
func (up *UserProcess) ServerProcessRegister(mes *common.Message) (err error) {

	//1.取出mes 中的Data 并直接序列化程RegisterMes
	var rgMes common.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &rgMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}

	// 声明一个resMes
	var resMes common.Message
	resMes.Type = common.RegisterResMesType
	var rgResMes message.RegisterResMes

	// 用Redis数据库完成注册
	err = model.MyUserDao.Register(&rgMes.User)

	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误..."
		}
	} else {
		registerResMes.Code = 200
	}

	data, err := json.Marshal(rgResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	// 将Data 赋值给resMes
	resMes.Data = string(data)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	// 发送数据包
	tf := &utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)
	return
}

//编写专门函数处理登录请求

//编写一个函数serverProcessLogin函数， 专门处理登录请求
func (up *UserProcess) ServerProcessLogin(mes *common.Message) (err error) {

	//1. 先从mes 中取出 mes.Data ，并直接  反序列化成LoginMes
	var loginMes common.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	//1先声明一个 resMes
	var resMes common.Message
	resMes.Type = common.LoginResMesType
	//2在声明一个 LoginResMes，并完成赋值
	var loginResMes common.LoginResMes

	//我们需要到redis数据库去完成验证.
	//1.使用model.MyUserDao 到redis去验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)

	if err != nil {

		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误..."
		}

	} else {
		loginResMes.Code = 200
		//这里，因为用户登录成功，我们就把该登录成功的用放入到userMgr中
		//将登录成功的用户的userId 赋给 up
		up.UserId = loginMes.UserId
		userMgr.AddOnlineUser(up)
		//通知其它的在线用户， 我上线了
		up.NotifyOthersOnlineUser(loginMes.UserId)
		//将当前在线用户的id 放入到loginResMes.UsersId
		//遍历 userMgr.onlineUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		fmt.Println(user, "登录成功")
	}
	// //如果用户id= 100， 密码=123456, 认为合法，否则不合法

	// if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
	// 	//合法
	// 	loginResMes.Code = 200

	// } else {
	// 	//不合法
	// 	loginResMes.Code = 500 // 500 状态码，表示该用户不存在
	// 	loginResMes.Error = "该用户不存在, 请注册再使用..."
	// }

	//3将 loginResMes 序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	//4. 将data 赋值给 resMes
	resMes.Data = string(data)

	//5. 对resMes 进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}
	//6. 发送data, 我们将其封装到writePkg函数
	//因为使用分层模式(mvc), 我们先创建一个Transfer 实例，然后读取
	tf := &utils.Transfer{
		Conn: up.Conn,
	}
	err = tf.WritePkg(data)
	return
}
