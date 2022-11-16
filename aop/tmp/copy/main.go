package main

import "fmt"

func main() {
	var b = make([]uint32, 10)
	fmt.Printf("%p \n", b)
	BB(b)
	fmt.Println(b)
}

func BB(bb []uint32) {
	fmt.Printf("%p \n", bb)

	for i := 0; i < 10; i++ {
		copy(bb, []uint32{uint32(i)})
		fmt.Printf("%p \n", bb)

		bb = bb[1:]
	}
	fmt.Printf("%p \n", bb)

}
