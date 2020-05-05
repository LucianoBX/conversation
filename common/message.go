package common

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

// 用户状态常量
const (
	UserOffline = 0
	UserOnline  = 1 << iota
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// 各类性的消息
type LoginMes struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

type LoginResMes struct {
	Code   int `json:"code"`
	UserId []int
	Error  string `json:"error"`
}

type RegisterMes struct {
	User User `json:"user"`
}

type RegisterResMes struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

// 配合服务端推送用户状态变化消息
type NotifyUserStatusMes struct {
	UserId int `json:"userId"`
	Status int `json:"status"`
}

type SmsMes struct {
	Content string `json:"content"`
	User           // 匿名结构体
}
