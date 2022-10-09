package task

type Config struct {
	Id         uint32        `json:"id"`
	Name       string        `json:"name"`
	DropId     uint32        `json:"dropId"` //
	Category   int           `json:"category"`
	Targets    []*TargetConf `json:"targets"`
	SubmitType int           `json:"submitType"` //自动提交，手动提交
	AcceptType int           `json:"acceptType"`
}

type TargetConf struct {
}
