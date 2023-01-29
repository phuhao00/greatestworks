package broker

type Component interface {
	Resolve(opCh Operation)
	Launch()
	Stop()
}

type BaseComponent struct {
	stopCh chan struct{}
	opCh   chan Operation
}

func NewBaseComponent() *BaseComponent {
	return &BaseComponent{
		stopCh: make(chan struct{}),
		opCh:   make(chan Operation),
	}
}

func (b *BaseComponent) Resolve(op Operation) {
	b.opCh <- op
}

func (b *BaseComponent) Launch() {
	go func() {
		for {
			select {
			case <-b.stopCh:
				break
			case op := <-b.opCh:
				b.dealOp(op)
			}
		}
	}()
}

func (b *BaseComponent) Stop() {
	b.stopCh <- struct{}{}
}

func (b *BaseComponent) dealOp(operation Operation) {
	fn := func() {
		operation.CB()
		operation.Ret <- struct{}{}
	}
	if !operation.IsAsynchronous {
		fn()

	} else {
		go func() {
			fn()
		}()
	}
}
