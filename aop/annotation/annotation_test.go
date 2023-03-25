package annotation

import (
	"fmt"
	"github.com/phuhao00/greatestworks-proto/player"
	"google.golang.org/protobuf/proto"
	"reflect"
	"testing"
)

func TestAnnotation(t *testing.T) {
	a := Annotation{Prototypes: map[uint16]reflect.Type{}}
	a.Register(11, &player.CSLogin{})
	//
	data, _ := proto.Marshal(&player.CSLogin{
		UserName: "phuhao",
		Password: "123",
	})
	message := a.UnMarshal(11, data)
	fmt.Println(message)
}
