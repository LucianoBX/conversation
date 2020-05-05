package model

import (
	"encoding/json"
	"fmt"
	"github.com/conversation/common"
	"github.com/gomodule/redigo/redis"
	//"github.com/garyburd/redigo/redis"
)

// 做一个全局userDao，需要时直接调用
var (
	MyUserDao *UserDao
)

type UserDao struct {
	pool *redis.Pool
}

// 使用工厂模式，创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {

	userDao = &UserDao{
		pool: pool,
	}
	return
}

// 为UserDao 创建需要的方法
// 根据Id取得用户实例
func (ud *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {

	// 通过给定id 去 redis 查询这个用户
	//	iface, err := conn.Do("HGet", "users", id)
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		if err == redis.ErrNil {
			//表示在 users 哈希中，没有找到对应id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &User{}

	// 反序列化-> User
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

//完成校验，并实现登录
func (ud *UserDao) Login(userId int, userPwd string) (user *User, err error) {

	// 获取conn
	conn := ud.pool.Get()
	defer conn.Close()

	user, err = ud.getUserById(conn, userId)
	if err != nil {
		return
	}

	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return

}

// 完成注册方法
func (ud *UserDao) Register(user *common.User) (err error) {

	// 去一个连接
	conn := ud.pool.Get()
	defer conn.Close()

	_, err = ud.getUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}

	// 至此该id还没由注册，完成注册
	data, err := json.Marshal(user)
	if err != nil {
		return
	}

	//入库
	_, err = conn.Do("HSet", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误， err=", err)
		return
	}
	return
}
