package annotation

import (
	"google.golang.org/protobuf/proto"
	"reflect"
)

type Annotation struct {
	Prototypes map[uint16]reflect.Type
}

func (a *Annotation) Register(key uint16, pm proto.Message) {
	if _, ok := a.Prototypes[key]; ok {
		return
	}
	of := reflect.TypeOf(pm)
	a.Prototypes[key] = of.Elem()
}

func (a *Annotation) GetProtoType(id uint16) reflect.Type {
	return a.Prototypes[id]
}

func (a *Annotation) UnMarshal(id uint16, data []byte) proto.Message {
	var ret proto.Message
	if protoType := a.GetProtoType(id); protoType != nil {
		ret = reflect.New(protoType).Interface().(proto.Message)
		err := proto.Unmarshal(data, ret)
		if err != nil {
			return nil
		}
		return ret
	}
	return nil
}
