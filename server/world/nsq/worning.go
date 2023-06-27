package nsq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"greatestworks/aop/fn"
	"greatestworks/aop/logger"
	"net/http"
	"time"
)

var messageSendTime int64

type PanicDingDingNotify struct {
	MsgType string        `json:"msgtype"`
	Text    *DingDingText `json:"text"`
	At      *DingDingAt   `json:"at"`
}
type DingDingText struct {
	Content string `json:"content"`
}

type DingDingAt struct {
	IsAtAll bool `json:"isAtAll"`
}

func SendWarningMessage(ip string, errType string, err interface{}, stackMsg string) {
	now := time.Now().Unix()
	if now-messageSendTime < 100 {
		return
	}

	messageSendTime = now

	userList := []string{"calvin"}
	user := fn.GetUser()

	find := false
	for _, userName := range userList {
		if userName == user {
			find = true
			break
		}
	}

	if !find {
		return
	}

	msg := &PanicDingDingNotify{
		MsgType: "text",
		Text:    &DingDingText{},
		At:      &DingDingAt{IsAtAll: true},
	}

	content := fmt.Sprintf("ip:%v user:%v errType:%v error:%v\n content:%v", ip, user, errType, err, stackMsg)
	msg.Text.Content = content

	url := "https://oapi.dingtalk.com/robot/send?access_token="
	jsonValue, err := json.Marshal(msg)
	if err != nil {
		logger.Error("err:%v", err)
		return
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {

	}
}
