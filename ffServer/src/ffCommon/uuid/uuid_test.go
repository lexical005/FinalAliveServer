package uuid

import (
	"ffCommon/log/log"
	"ffCommon/util"

	"sync"
	"testing"
)

// func Test_Bug_Gen_Duplicate_uuid(t *testing.T) {
// 	id := 71838114186780

// 	requester := uint64(1)
// 	u1, _ := New(requester)

// 	// u.requester + (u.sn << sn_bit_count) + now

// 	v := id - u1.requester
// 	log.RunLogger.Println("u1.requester[%v] id - u1.requester[%v]", u1.requester, v)
// }

func Test_UUID_1(t *testing.T) {
	count := 100

	requester := uint64(1)

	u1, _ := NewGeneratorSafe(requester)
	result := make(map[UUID]string, count)
	for i := 0; i < count; i++ {
		id := u1.Gen()
		if v, ok := result[id]; ok {
			t.Errorf("Test_UUID_1 requester gen duplicate uuid[%v] [%v->%v]\n", id, v, u1.String())
		} else {
			result[id] = u1.String()
		}
	}

	log.RunLogger.Println("Test_UUID_1 passed")
}

func Test_UUID_2(t *testing.T) {
	count := 100

	requester1, requester2 := uint64(1), uint64(2)

	u1, _ := NewGeneratorSafe(requester1)
	u2, _ := NewGeneratorSafe(requester2)
	result1 := make(map[UUID]string, count)
	result2 := make(map[UUID]string, count)
	for i := 0; i < count; i++ {
		id1 := u1.Gen()
		if v, ok := result1[id1]; ok {
			t.Errorf("Test_UUID_2 requester1 gen duplicate uuid[%v] [%v->%v]\n", id1, v, u1.String())
		} else {
			result1[id1] = u1.String()
		}

		id2 := u2.Gen()
		if v, ok := result2[id2]; ok {
			t.Errorf("Test_UUID_2 requester2 gen duplicate uuid[%v] [%v->%v]\n", id2, v, u2.String())
		} else {
			result2[id2] = u2.String()
		}
	}

	for id := range result1 {
		if _, ok := result2[id]; ok {
			t.Errorf("Test_UUID_2 multi requester[%v, %v] gen duplicate uuid[%v]", requester1, requester2, id)
		}
	}

	log.RunLogger.Println("Test_UUID_2 passed")
}

func genUUIDAsync(params ...interface{}) {
	t, _ := params[0].(*testing.T)
	u, _ := params[1].(Generator)
	loopCount, _ := params[2].(int)
	result, _ := params[3].(map[UUID]string)
	wg, _ := params[4].(*sync.WaitGroup)
	muLock, _ := params[5].(*sync.Mutex)

	resultTmp := make(map[UUID]string, loopCount)

	for i := 0; i < loopCount; i++ {
		id := u.Gen()
		if v, ok := resultTmp[id]; ok {
			t.Errorf("genUUIDAsync gen duplicate uuid[%v] [%v->%v]\n", id, v, u.String())
		} else {
			resultTmp[id] = u.String()
		}
	}

	muLock.Lock()
	defer muLock.Unlock()

	for id1, v1 := range resultTmp {
		if v2, ok := result[id1]; ok {
			t.Errorf("genUUIDAsync multi goroutines gen duplicate uuid[%v] [%v->%v]\n", id1, v1, v2)
		} else {
			result[id1] = v2
		}
	}

	wg.Done()
}

func Test_UUIDAsync_1(t *testing.T) {
	goCount := 10
	loopCount := 100

	requester := uint64(11)

	u1, _ := NewGeneratorSafe(requester)
	result1 := make(map[UUID]string, loopCount)

	var muLock sync.Mutex

	var wg sync.WaitGroup
	for i := 0; i < goCount; i++ {
		wg.Add(1)
		go util.SafeGo(genUUIDAsync, t, u1, loopCount, result1, &wg, &muLock)
	}

	wg.Wait()

	log.RunLogger.Println("Test_UUIDAsync_1 passed")
}
