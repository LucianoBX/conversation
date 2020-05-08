package process

import (
	"fmt"
	"github.com/conversation/client/model"
	"github.com/conversation/common"
)

// 客户段维护的map
var onlineUsers map[int]*common.User = make(map[int]*common.User, 16)
var CurrUser model.CurrUser

// 在客户端显示当前在线的用户
func outputOnlineUser() {

	fmt.Println("当前在线用户列表")
	// 遍历输出
	for id, _ := range onlineUsers {
		if id == CurrUser.UserId {
			continue
		}
		fmt.Println("用户ID：\t", id)
	}
}

// 编写一个方法，处理返回的NotifyuserStatusMes
func updateUserStatus(notifyUserStatusMes *common.NotifyUserStatusMes) {
	//
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		user = &common.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}

	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user

	outputOnlineUser()
}
