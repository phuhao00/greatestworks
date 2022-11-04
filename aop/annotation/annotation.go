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
	a.Prototypes[key] = of
}
