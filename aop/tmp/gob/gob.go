package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Getter interface {
	Get() string
}

type Foo struct {
	Bar string
}

func (f Foo) Get() string {
	return f.Bar
}

func main() {
	buf := bytes.NewBuffer(nil) //创建一个缓冲区，且用空去初始化他，如果这个nil换成其他的buffer就相当于用传入的buffer赋值给新的buffer
	gob.Register(Foo{})         //注册Foo结构因为我们要编码interface

	//create a interface getter of Foo
	g := Getter(Foo{"zhr"})

	//encode interface getter
	enc := gob.NewEncoder(buf)
	enc.Encode(&g)

	fmt.Println("after encode is ", buf)
	//decode
	dec := gob.NewDecoder(buf)
	var gg Getter
	if err := dec.Decode(&gg); err != nil {
		panic(err)
	}
	fmt.Println("after decode is ", gg)
}
