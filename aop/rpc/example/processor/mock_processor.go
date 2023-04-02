package processor

import "fmt"

type MockProcessor int

type MockParam struct {
	Tag string
}

func (p MockProcessor) Print(req, rsp *string) error {
	fmt.Println(req)
	tmp := "hi,world"
	rsp = &tmp
	return nil
}

func (p MockProcessor) Print2(req, rsp *MockParam) error {
	fmt.Println(req)

	rsp.Tag = "abc"
	return nil
}
