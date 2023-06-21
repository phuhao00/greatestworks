package task

type ITarget interface {
	GetProgress()
	GetTotalProgress()
	CheckDone()
	OnAccept()
	OnEvent(interface{})
	ConfigVerify()
	Notify()
}
