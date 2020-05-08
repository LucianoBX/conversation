package process

import (
	"encoding/json"
	"fmt"
	"github.com/conversation/common"
	// "github.com/conversation/server/model"
	"github.com/conversation/server/utils"
	"net"
)

type SmsProcess struct {
}

// 转发群发消息
func (sp *SmsProcess) SendGroupMes(mes *common.Message) {

	// 遍历在线用户，
	// 取出消息，转发

	var smsMes common.SmsMes
	err := json.Unmarshal([]byte(mes.Data), smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	for id, up := range userMgr.onlineUsers {

		if id == smsMes.UserId {

			continue
		}
		sp.SendMesToEachOnlineUser(data, up.Conn)

	}
}

func (sp *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {

	// 创建Transfer实例，发送data
	tf := &utils.Transfer{

		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err = ", err)
	}
}
