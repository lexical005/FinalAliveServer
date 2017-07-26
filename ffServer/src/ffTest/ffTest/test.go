package main

import (
	"ffCommon/log/log"
	"ffCommon/util"

	"regexp"
	"sync"
)

func testWG() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Add(1)
	wg.Done()
	wg.Done()

	wg.Wait()

	log.RunLogger.Println("ok")

	wg = sync.WaitGroup{}
	wg.Add(1)
	wg.Add(1)
	wg.Done()
	wg.Done()

	wg.Wait()
	log.RunLogger.Println("ok")
}

func testSlice() {
	b := make([]int, 0, 10)

	a := b[0:8]

	c := a[0:10]

	log.RunLogger.Println(len(b), cap(b))
	log.RunLogger.Println(len(a), cap(a))
	log.RunLogger.Println(len(c), cap(c))

	a = append(a, 1)

	log.RunLogger.Println(len(b), cap(b))
	log.RunLogger.Println(len(a), cap(a))
	log.RunLogger.Println(len(c), cap(c))

}

var ch1 = make(chan bool, 1)

func f1() bool {
	return <-ch1
}

func testFile() {
	err := util.RemovePath("test.txt")
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	err = util.CreatePath("test")
	if err != nil {
		log.RunLogger.Println(err)
	}
}

func testc() {
	var i it
	i = &c2{
		c1: &c1{},
	}
	log.RunLogger.Println(i.Type())
}

func testregexp() {
	s := "[123, 321,]"
	r := regexp.MustCompile(`[\d]+`)
	log.RunLogger.Println(r.FindAllString(s, -1))
}
