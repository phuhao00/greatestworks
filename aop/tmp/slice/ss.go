package slice

import "fmt"

// s[:1:4] 每个数字前都有个冒号， slice内容为data从0到第1位，长度len为1，最大扩充项cap设置为4

func S() {
	var s = []int{1, 2, 3, 5, 8}
	fmt.Println(cap(s))

	newne := s[:1:4]
	fmt.Println(newne)
	fmt.Println(cap(newne))

}
