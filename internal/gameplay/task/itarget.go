package task

type Target interface {
	GetProgress()
	GetTotalProgress()
	CheckDone()
	OnAccept()
	OnEvent(interface{})
	ConfigVerify()
	Notify()
}
