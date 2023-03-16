package buff

type Base struct {
	Level     int
	IsForever bool
	ConfId    uint32
	Id        int64
	StartTime int64
	EndTime   int64
	Desc      string
}

func (b *Base) OnStart() {
	//TODO implement me
	panic("implement me")
}

func (b *Base) OnEnd() {
	//TODO implement me
	panic("implement me")
}
