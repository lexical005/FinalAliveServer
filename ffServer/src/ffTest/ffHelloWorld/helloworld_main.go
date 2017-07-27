package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World! 你好，世界!")

	// buffer length and capcity
	a := make([]byte, 0, 10)
	a1 := a[0:0]
	a2 := a[0:10]
	println(len(a), cap(a))   // 0, 10
	println(len(a1), cap(a1)) // 0, 10
	println(len(a2), cap(a2)) // 0, 10

	a3 := a1[0:10]
	println(len(a3), cap(a3)) // 10, 10

	a4 := a2[0:0]
	println(len(a4), cap(a4)) // 0, 10
}
