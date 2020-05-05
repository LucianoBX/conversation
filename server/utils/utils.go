package utils

import (
	"encoding/binary"
	"enconing/json"
	"errors"
	"fmt"
	"github.com/conversation/common"
	"net"
)

type Transfer struct {
	// 连接字段
	Conn net.Conn
	// buffer
	Buf [8096]byte
}

func (t *Transfer) ReadPkg() (mes message.Message, err errors) {

	fmt.Println("读取客户端发送消息……")

	// read the first 4 byte
	_, err = t.Conn.Read(t.Buf[:4])
	if err != nil {
		return
	}

	// 将buf[:4]转成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(t.Buf[:4])

	// read message accord to pkgLen
	n, err := t.Conn.Read(t.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		return
	}

	// pkg反序列化
	//
	err = json.Unmarhal(t.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}
	return
}

func (t *Transfer) WritePkg(data []byte) (err error) {

	// 先发送长度用于判断
	// 获得数据长度
	var pkgLen uint32
	pkgLen = uint32(len(data))

	// 写入长度数据
	binary.BigEndian.PutUint32(t.Buf[:4], pkgLen)

	// 发送长度
	n, err := t.Conn.Write(t.Buf[:4])
	if err != nil || n != 4 {
		fmt.Print("conn.Write(bytes) fall")
		return
	}

	// 发送数据本身
	n, err = t.Conn.Write(data)
	if err != nil || n != int(pkgLen) {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	return
}
