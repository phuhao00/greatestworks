package attribute

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type Category uint32

const (
	OnceAttribute    Category = iota + 0 // 不会重置
	DayAttribute                         // 每日重置
	WeekAttribute                        // 每周重置
	MonthAttribute                       // 每月重置
	YearAttribute                        // 每年重置
	WeekNatAttribute                     // 自然周重置
)

type Attribute struct {
	Value    string
	SetTime  int64
	AttrType uint32
}

type Attributes struct {
	Id                uint64
	dayAttributes     *sync.Map
	weekAttributes    *sync.Map
	monthAttributes   *sync.Map
	yearAttributes    *sync.Map
	onceAttributes    *sync.Map
	ClientAttributes  *sync.Map
	weekNatAttributes *sync.Map
}

func (attrs *Attributes) Clear() {
	attrs.dayAttributes = nil
	attrs.weekAttributes = nil
	attrs.monthAttributes = nil
	attrs.yearAttributes = nil
	attrs.onceAttributes = nil
	attrs.ClientAttributes = nil
	attrs.weekNatAttributes = nil
}

func (attrs *Attributes) LoadFromDB() {

}

func (attrs *Attributes) SaveDB() {

}

func (attrs *Attributes) Set(c Category, key string, v interface{}) {

}

func (attrs *Attributes) Get(c Category, key string) *Attribute {
	return nil
}

func ValueToString(i interface{}) string {

	switch v := i.(type) {

	case nil:
		return ""
	case string:
		return v
	case []byte:
		return string(v)
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(v).Uint(), 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', 3, 32)
	case float64:
		return strconv.FormatFloat(float64(v), 'f', 3, 64)
	case bool:
		return strconv.FormatBool(v)
	case []uint32:
		return NumberSliceToString(v, ",")
	case []int32:
		return NumberSliceToString(v, ",")
	case []uint64:
		return NumberSliceToString(v, ",")
	default:
		return ""
	}
	
}

func NumberSliceToString[T uint32 | int32 | int64 | uint64](s []T, delim string) string {
	return strings.Trim(strings.Join(strings.Split(fmt.Sprint(s), " "), delim), "[]")
}
