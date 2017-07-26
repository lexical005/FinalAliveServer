package pool

import (
	"ffCommon/log/log"

	"testing"
)

func Test_FreeCache_1(t *testing.T) {
	value := 0
	creator := func() interface{} {
		value++
		return value
	}

	r := make([]interface{}, 0, 0)
	c := New("Test_FreeCache_1", true, creator, 0, 50)
	var v interface{}

	step := 1
	log.RunLogger.Printf("step[%v] c[%v]\n", step, c)

	// 正常获取/扩展
	r = append(r, c.Apply())
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	r = append(r, c.Apply())
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	r = append(r, c.Apply())
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	r = append(r, c.Apply())
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	log.RunLogger.Println("")

	// 归还
	v = r[len(r)-1]
	r = r[:len(r)-1]
	c.Back(v)
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	v = r[len(r)-1]
	r = r[:len(r)-1]
	c.Back(v)
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	v = r[len(r)-1]
	r = r[:len(r)-1]
	c.Back(v)
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	v = r[len(r)-1]
	r = r[:len(r)-1]
	c.Back(v)
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	// 超量归还，不报错
	c.Back(v)
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	log.RunLogger.Println("Test_FreeCache_1 passed")
}

func Test_LimitCache_1(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			log.RunLogger.Println(err)

			log.RunLogger.Println("Test_LimitCache_1 passed")
		}
	}()

	value := 0
	creator := func() interface{} {
		value++
		return value
	}

	r := make([]interface{}, 0, 0)
	c := New("Test_LimitCache_1", false, creator, 0, 50)
	var v interface{}

	step := 1
	log.RunLogger.Printf("step[%v] c[%v]\n", step, c)

	// 正常获取/扩展
	r = append(r, c.Apply())
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	r = append(r, c.Apply())
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	r = append(r, c.Apply())
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	r = append(r, c.Apply())
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	log.RunLogger.Println("")

	// 归还
	v = r[len(r)-1]
	r = r[:len(r)-1]
	c.Back(v)
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	v = r[len(r)-1]
	r = r[:len(r)-1]
	c.Back(v)
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	v = r[len(r)-1]
	r = r[:len(r)-1]
	c.Back(v)
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	v = r[len(r)-1]
	r = r[:len(r)-1]
	c.Back(v)
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)

	// 报错
	c.Back(v)
	step++
	log.RunLogger.Printf("step[%v] r[%v] c[%v]\n", step, r, c)
}
