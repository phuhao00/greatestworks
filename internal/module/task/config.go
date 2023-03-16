package task

import (
	"github.com/phuhao00/greatestworks-proto/gen/messageId"
	"github.com/phuhao00/network"
)

type Status int

const (
	ACCEPT Status = iota + 1
	ING
	FINISH
	SUBMIT
)

const (
	defaultLoopNum     = 50
	defaultMonitor     = 100
	defaultChanInSize  = 1000
	defaultChanOutSize = 500
)

type ModuleConfig struct {
	LoopNum    int
	MonitorNum int
	ChInSize   int
	ChOutSize  int
}

type Config struct {
	Id              uint32        `json:"id"`
	Name            string        `json:"name"`
	DropId          uint32        `json:"dropId"` //
	Category        int           `json:"category"`
	Targets         []*TargetConf `json:"targets"`
	SubmitType      int           `json:"submitType"` //自动提交，手动提交
	AcceptType      int           `json:"acceptType"`
	CompleteNtf     int           `json:"completeNtf"` //完成是否推送
	UnlockCondition int           `json:"unlockCondition"`
	Module          string        `json:"module"`
}

type TargetConf struct {
	Id            uint32
	DropId        uint32
	Name          string
	CompleteParam string
}

type PlayerActionParam struct {
	MessageId messageId.MessageId
	Player    Player
	Packet    *network.Message
}
