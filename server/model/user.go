package model

// 定义一个用户结构提
type User struct {
	// 注意tag信息与, 与结构体field 一致
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}
