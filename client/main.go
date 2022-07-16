package main

func main() {
	c := NewClient()
	c.InputHandlerRegister()
	c.Run()
	select {}
}
