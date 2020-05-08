package process

import (
	"encoding/json"
	"fmt"
	"github.com/conversation/client/utils"
	"github.com/conversation/common"
)

type SmsProcess struct {
}

// 群聊消息发送
func (sp *SmsProcess) SendGroupMes(content string) (err error) {

	// 创建mes
	var mes common.Message
	mes.Type = common.SmsMesType

	// 创建SmsMes
	var smsMes common.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurrUser.UserId
	smsMes.UserStatus = CurrUser.UserStatus

	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SentGropMes json.Marshal sms fail = ", err)
		return

	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SentGropMes json.Marshal  mes fail = ", err)
		return
	}

	tf := &utils.Transfer{
		Conn: CurrUser.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMes err =", err.Error())
		return
	}
	return
}
