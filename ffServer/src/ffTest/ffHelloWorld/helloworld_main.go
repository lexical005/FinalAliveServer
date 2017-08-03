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

	a = append(a, 1)
	a = append(a, 2)
	a = append(a, 3)
	a = append(a, 4)
	a = append(a, 5)
	b := a

	a = a[:0]
	a = append(a, b[2:]...)
	fmt.Printf("%v", a)

	a = a[:0]
	a = append(a, 1)
	a = append(a, 2)
	a = append(a, 3)
	a = append(a, 4)
	a = append(a, 5)

	s := &st1{
		a: a,
	}
	s.test1()

	testCC()

	testChan1()

	gotest()
}

type st1 struct {
	a []byte
}

func (s *st1) test1() {
	s.test2(s.a)
}

func (s *st1) test2(buf []byte) {
	s.a = s.a[:0]
	s.a = append(s.a, buf[2:]...)
	fmt.Printf("%v", s.a)
}

func testCC() {
	c1 := make(chan int, 2)
	c2 := make(chan int, 1)
	c1 <- 1
	c2 <- <-c1
	println(<-c2)
}

// select 自上而下求值(必然进行!), 若有多个case满足, 则随机挑选一个case执行
func testChan1() {
	println("testChan1")

	f := func() int {
		c1 := make(chan int, 1)
		c2 := make(chan int, 1)
		c1 <- 1
		c2 <- 1

		select {
		case <-c1:
			break
		case <-c2:
			break
		}
		if len(c1) == 0 {
			return -1
		}
		return 1
	}
	v := 0
	for i := 0; i < 105; i++ {
		v += f()
	}
	println("result:")
	println(v)
}

func gotest() {
	i1, i2 := 1, 1
	defer func() {
		println("gotest main 1")
		if err := recover(); err != nil {
			println(err)
		}
		println(i1)
	}()
	defer func() {
		println("gotest main 2")
		if err := recover(); err != nil {
			println(err)
		}
		i1 = i1 / (i2 - 1)
	}()
	println("gotest main")
	i2 = i1 / (i1 - 1)
}
