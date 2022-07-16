package main

func main() {
	c := NewClient()
	c.Run()
	select {}
}
