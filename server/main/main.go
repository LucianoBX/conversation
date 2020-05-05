package main

import (
	"fmt"
	"github.com/conversation/server/model"
	"net"
	"time"
)

// comunicate with clients
func process(conn net.Conn) {
	// delay the close of conn
	defer conn.Close()

	// 调用总控
	processcor := &Processor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误， err=", err)
	}
}
