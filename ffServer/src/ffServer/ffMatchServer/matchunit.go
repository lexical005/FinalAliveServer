package main

// iMatchUnit 匹配单元
type iMatchUnit interface {
	// AllReady 匹配单元是不是已经全部准备
	AllReady() bool

	// Count 匹配单元内有多少matchPlayer
	Count() int

	// MatchSuccess 进入了准备组, 匹配完成
	MatchSuccess()
}
