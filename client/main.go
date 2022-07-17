package main

func main() {
	c := NewClient()
	c.InputHandlerRegister()
	c.MessageHandlerRegister()
	c.Run()
	select {}
}
