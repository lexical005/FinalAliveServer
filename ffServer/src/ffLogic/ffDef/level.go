package ffDef

// ILevelMgr 关卡管理接口
type ILevelMgr interface {
	// StoryLevel 剧情关卡接口
	StoryLevel() IStoryLevel

	// Init 初始化
	Init()

	// Clear 清除自身, 以待重用
	Clear()
}

// IStoryLevel 剧情关卡数据接口
// levelID: 关卡完整编号(含类型,难度,区域,序号等信息)
// areaID: 关卡区域
// difficutyID: 关卡难度
type IStoryLevel interface {
	// BuyCount 累计购买次数与当天购买次数
	BuyCount() (int, int)

	// PassedCount 累计通关次数和关卡当天通关次数
	PassedCount() (int, int)

	// Score 累计通关评级和关卡当天通关评级
	Score() (int, int)

	// LevelPassedCount 特定关卡累计通关次数和关卡当天通关次数
	LevelPassedCount(levelID int) (int, int)

	// AreaPassedCount 特定区域累计通关次数和关卡当天通关次数
	AreaPassedCount(areaID int) (int, int)

	// DifficutyPassedCount 特定难度累计通关次数和关卡当天通关次数
	DifficutyPassedCount(difficutyID int) (int, int)

	// LevelScore 特定关卡通关评级
	LevelScore(levelID int) int

	// AreaScore 特定区域总通关评级
	AreaScore(areaID int) int

	// DifficutyScore 特定难度总通关评级
	DifficutyScore(difficutyID int) int

	// LevelPass 通关关卡
	LevelPass(levelID, score int)

	// Init 初始化
	Init()

	// Clear 清除自身, 以待重用
	Clear()
}
