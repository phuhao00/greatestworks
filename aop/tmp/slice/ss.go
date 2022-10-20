package slice

import "fmt"

func S() {
	var s = []int{1, 2, 3, 5, 8}
	fmt.Println(cap(s))

	newne := s[:1:4]
	fmt.Println(newne)
	fmt.Println(cap(newne))

}
