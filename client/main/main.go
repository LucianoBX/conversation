package main

import (
	"fmt"
	"github.com/conversation/client/process"
	"os"
)

//用户ID，与用户密码
var (
	userId            int
	userName, userPwd string
)

func main() {
	// 接受用户选择
	var key int

	for true {
		fmt.Println("------------登录多人聊天系统------------")
		fmt.Println("\t\t\t 1 登录聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 请选择（1-3）：")

		fmt.Scanf("%d\n", &key)

		// 菜单选择
		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户id")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码")
			fmt.Scanf("%s\n", &userPwd)

			// 完成登录
			up := &process.UserProcess{}
			up.Login(userId, userPwd)

		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码:")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户名字(nickname):")
			fmt.Scanf("%s\n", &userName)

			// 完成注册
			up := &process.UserProcess{}
			up.Register(userId, userPwd, userName)

		case 3:
			fmt.Println("退出系统")
			os.Exit(0)

		default:
			fmt.Println("你输入有误，请重新输入")
		}
	}
}
