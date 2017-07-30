package ffProto

import (
	"ffCommon/log/log"
	"ffCommon/pool"

	"fmt"
)

var arraySizeLimit = []int{16, 32, 64, 128, 256, 512, 1024, 2048, 4096}
var arrayInitCount = []int{100, 100, 100, 100, 50, 50, 20, 10, 5}

var bufferPool []*pool.Pool

func init() {
	l1, l2 := len(arraySizeLimit), len(arrayInitCount)
	if l1 != l2 {
		panic(fmt.Sprintf("ffProto.proto_pool len(arraySizeLimit)[%d] != len(arrayInitCount)[%d]", l1, l2))
	}

	for index := 1; index < l1; index++ {
		if arraySizeLimit[index-1] >= arraySizeLimit[index] {
			panic(fmt.Sprintf("ffProto.proto_pool invalid arraySizeLimit at index %d:%d: size must increase", index-1, index))
		} else if arraySizeLimit[index] > protoMaxLength {
			panic(fmt.Sprintf("ffProto.proto_pool invalid arraySizeLimit: size must letter than protoMaxLength[%d]", protoMaxLength))
		}

		if arrayInitCount[index-1] < arrayInitCount[index] {
			panic(fmt.Sprintf("ffProto.proto_pool invalid arrayInitCount at index %d:%d: init count must decrease", index-1, index))
		}
	}

	bufferPool = make([]*pool.Pool, l1, l1)
	for i, v := range arrayInitCount {
		bufferPool[i] = newBufferPool(arraySizeLimit[i], v)
	}
}

func newBufferPool(bufSize, initCount int) *pool.Pool {
	creator := func() interface{} {
		return make([]byte, 0, bufSize)
	}

	return pool.New(fmt.Sprintf("ffProto.proto_pool.bufferPool_%v", bufSize), true, creator, initCount, 50)
}

// findBufferIndex 根据bufLengthLimit寻找合适的bufferPool下标
// 	bufLengthLimit: >=0 限定协议缓冲区大小; <0 协议缓冲区越大越好
func findBufferIndex(bufLengthLimit int) int {
	index := len(arraySizeLimit) - 1
	if bufLengthLimit >= 0 {
		for i, v := range arraySizeLimit {
			if bufLengthLimit <= v {
				index = i
				break
			}
		}
	}
	return index
}

// upperProtoBufferLength 将实际所需的缓冲区bufLengthLimit转换未更优的大小
func upperProtoBufferLength(bufLengthLimit int) int {
	for _, limit := range arraySizeLimit {
		if bufLengthLimit <= limit {
			return limit
		}
	}
	log.FatalLogger.Printf("upperProtoBufferLength bufLengthLimit[%v] is too large", bufLengthLimit)
	return bufLengthLimit
}

func applyBuffer(bufLengthLimit int) []byte {
	index := findBufferIndex(bufLengthLimit)
	buf, _ := bufferPool[index].Apply().([]byte)
	return buf
}

func backBuffer(buf []byte) {
	index := findBufferIndex(cap(buf))
	bufferPool[index].Back(buf)
}
