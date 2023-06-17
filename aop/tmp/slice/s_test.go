package slice

import (
	"fmt"
	"sort"
	"testing"
)

func TestSlice(t *testing.T) {
	S()
	//v := [3]int{}
	//fmt.Println(reflect.TypeOf(v).Kind())
}

type AA struct {
	A uint64
}

func TestSort(t *testing.T) {
	data := []AA{{1}, {2}, {5}, {31}, {3}}
	sort.Slice(data, func(i, j int) bool {
		return data[i].A < data[j].A
	})

	fmt.Println(data)
}
