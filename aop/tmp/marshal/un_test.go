package marshal

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	var m = map[int]int{}
	m[1] = 2
	m[2] = 2
	m[3] = 2
	bytes, _ := json.Marshal(m)
	var m2 map[int]int
	json.Unmarshal(bytes, &m2) //需要加指针
	fmt.Println(m2)

}

func Result() (ret int) {
	defer func() {
		fmt.Println(ret)
	}()
	return 2
}

func TestRet(t *testing.T) {
	Result()
}
