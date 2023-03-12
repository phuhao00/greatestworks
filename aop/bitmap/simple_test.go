package bitmap

import (
	"fmt"
	"testing"
)

func TestSimple(t *testing.T) {
	//a1 := byte('a')
	//fmt.Println(a1)
	bitMap := NewBitMap(20)
	bitMap.set(20)
	bitMap.set(1)
	bitMap.set(2)
	bitMap.set(3)
	bitMap.set(4)
	bitMap.set(5)
	fmt.Println("20:", bitMap.isExist(20))
	fmt.Println("5:", bitMap.isExist(5))
	bitMap.del(20)
	fmt.Println("20:", bitMap.isExist(20))
	fmt.Println("5:", bitMap.isExist(5))
	fmt.Println("1:", bitMap.isExist(1))
	fmt.Println("3:", bitMap.isExist(3))
	fmt.Println("6:", bitMap.isExist(6))
}

func TestToggle(t *testing.T) {
	bitMap := NewBitMap(20)
	bitMap.set(20)
}
