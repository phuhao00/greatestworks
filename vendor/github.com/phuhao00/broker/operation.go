package broker

type CallBackFn func()

type Operation struct {
	IsAsynchronous bool
	CB             CallBackFn
	Ret            chan interface{}
}
