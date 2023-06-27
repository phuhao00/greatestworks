package nsq

import (
	_ "github.com/phuhao00/sugar"
	"greatestworks/aop/fn"
)

var osUserName = fn.GetUser()

var (
	PublicChat  = "Chat" + "-" + osUserName
	PrivateChat = "PrivateChat" + "-" + osUserName
	SystemMsg   = "SystemMsg" + "-" + osUserName
	Complex     = "Complex" + "-" + osUserName
	World       = "World" + "-" + osUserName
)

const (
	ChatNSQ  uint32 = 0
	LogicNSQ uint32 = 1
)
