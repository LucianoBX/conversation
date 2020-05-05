package model

import (
	"net"
	"github.com/conversation/common"
)

type CurrUser struct {
	Conn net.Conn
	common.User
}
