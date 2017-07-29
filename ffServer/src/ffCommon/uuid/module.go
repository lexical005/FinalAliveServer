package uuid

import "fmt"

// // NewGeneratorUnsafe 返回一个非多协程安全的Generator
// // 	requester: 使用者自定义的12位无符号数, 可根据项目需要, 全局规划这12位
// func NewGeneratorUnsafe(requester uint64) (Generator, error) {
// 	if requester > requesterBitMask {
// 		return nil, fmt.Errorf("uuid.NewGeneratorUnsafe: invalid requester, must between[0, %v]", requesterBitMask)
// 	}

// 	return &uuidGeneratorUnsafe{
// 		requester: requester << requesterBitOffset,
// 	}, nil
// }

// // NewGeneratorSafe 返回一个多协程安全的Generator
// // 	requester: 使用者自定义的12位无符号数, 可根据项目需要, 全局规划这12位
// func NewGeneratorSafe(requester uint64) (Generator, error) {
// 	if requester > requesterBitMask {
// 		return nil, fmt.Errorf("uuid.NewGeneratorSafe: invalid requester, must between[0, %v]", requesterBitMask)
// 	}

// 	return &uuidGeneratorSafe{
// 		Generator: &uuidGeneratorUnsafe{
// 			requester: requester << requesterBitOffset,
// 		},
// 	}, nil
// }

// NewGeneratorSafe 返回一个多协程安全的Generator
// 	requester: 使用者自定义的12位无符号数, 可根据项目需要, 全局规划这12位
func NewGeneratorSafe(requester uint64) (Generator, error) {
	if requester > requesterBitMask {
		return nil, fmt.Errorf("uuid.NewGeneratorSafe: invalid requester, must between[0, %v]", requesterBitMask)
	}

	return &uuidGeneratorSafe{
		Generator: &uuidGeneratorUnsafe{
			requester: requester << requesterBitOffset,
		},
	}, nil
}
