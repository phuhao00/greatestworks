package main

func (w *World) Reload() {
	//TODO implement me
	panic("implement me")
}

func (w *World) Init() {
	//TODO implement me
	panic("implement me")
}

func (w *World) Start() {
	w.HandlerRegister()
	go w.Server.Run()
	go w.pm.Run()
}

func (w *World) Stop() {

}
