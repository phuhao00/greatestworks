package time

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	tmp := time.Now()
	fmt.Println(tmp.Unix())
	time.Sleep(time.Second * 3)
	tmp.Add(time.Second * 3)
	fmt.Println(tmp.Unix())
}
