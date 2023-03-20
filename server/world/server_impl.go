package main

func (w *World) Loop() {
	//TODO implement me
	panic("implement me")
}

func (w *World) Monitor() {
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
