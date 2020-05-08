package process

import (
	"encoding/json"
	"fmt"
	"github.com/conversation/common"
)

// 显示群发消息
func outputGroupMes(mes *common.Message) {

	var smsMes common.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err.Error())

		return
	}

	info := fmt.Sprintf("用户ID：%d对大家说：%s",
		smsMes.UserId, smsMes.Content)
	fmt.Println(info)
	fmt.Println()

}
