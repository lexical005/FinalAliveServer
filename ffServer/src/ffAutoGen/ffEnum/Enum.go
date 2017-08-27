package ffEnum

import (
	"fmt"
)

// EActorAttr EActorAttr
type EActorAttr int32

const (
	internalEActorAttrActor  EActorAttr = EActorAttr(0) // 角色外观属性-角色
	internalEActorAttrMask   EActorAttr = EActorAttr(1) // 面具
	internalEActorAttrPants  EActorAttr = EActorAttr(2) // 裤子
	internalEActorAttrShoes  EActorAttr = EActorAttr(3) // 鞋子
	internalEActorAttrShirt  EActorAttr = EActorAttr(4) // 衬衫
	internalEActorAttrBelt   EActorAttr = EActorAttr(5) // 腰带
	internalEActorAttrGloves EActorAttr = EActorAttr(6) // 手套
	internalEActorAttrJacket EActorAttr = EActorAttr(7) // 外衣
	internalEActorAttrHead   EActorAttr = EActorAttr(8) // 帽子/头盔
	internalEActorAttrVest   EActorAttr = EActorAttr(9) // 防弹衣
)

type internalEActorAttrInfo struct {
	value EActorAttr
	toml  string
	desc  string
}

var allEActorAttrInfo = []*internalEActorAttrInfo{
	&internalEActorAttrInfo{
		value: internalEActorAttrActor,
		toml:  "EActorAttr.Actor",
		desc:  "角色外观属性-角色",
	},
	&internalEActorAttrInfo{
		value: internalEActorAttrMask,
		toml:  "EActorAttr.Mask",
		desc:  "面具",
	},
	&internalEActorAttrInfo{
		value: internalEActorAttrPants,
		toml:  "EActorAttr.Pants",
		desc:  "裤子",
	},
	&internalEActorAttrInfo{
		value: internalEActorAttrShoes,
		toml:  "EActorAttr.Shoes",
		desc:  "鞋子",
	},
	&internalEActorAttrInfo{
		value: internalEActorAttrShirt,
		toml:  "EActorAttr.Shirt",
		desc:  "衬衫",
	},
	&internalEActorAttrInfo{
		value: internalEActorAttrBelt,
		toml:  "EActorAttr.Belt",
		desc:  "腰带",
	},
	&internalEActorAttrInfo{
		value: internalEActorAttrGloves,
		toml:  "EActorAttr.Gloves",
		desc:  "手套",
	},
	&internalEActorAttrInfo{
		value: internalEActorAttrJacket,
		toml:  "EActorAttr.Jacket",
		desc:  "外衣",
	},
	&internalEActorAttrInfo{
		value: internalEActorAttrHead,
		toml:  "EActorAttr.Head",
		desc:  "帽子/头盔",
	},
	&internalEActorAttrInfo{
		value: internalEActorAttrVest,
		toml:  "EActorAttr.Vest",
		desc:  "防弹衣",
	},
}

var mapCodeToEActorAttrInfo = map[string]*internalEActorAttrInfo{
	allEActorAttrInfo[int(internalEActorAttrActor)].toml:  allEActorAttrInfo[int(internalEActorAttrActor)],
	allEActorAttrInfo[int(internalEActorAttrMask)].toml:   allEActorAttrInfo[int(internalEActorAttrMask)],
	allEActorAttrInfo[int(internalEActorAttrPants)].toml:  allEActorAttrInfo[int(internalEActorAttrPants)],
	allEActorAttrInfo[int(internalEActorAttrShoes)].toml:  allEActorAttrInfo[int(internalEActorAttrShoes)],
	allEActorAttrInfo[int(internalEActorAttrShirt)].toml:  allEActorAttrInfo[int(internalEActorAttrShirt)],
	allEActorAttrInfo[int(internalEActorAttrBelt)].toml:   allEActorAttrInfo[int(internalEActorAttrBelt)],
	allEActorAttrInfo[int(internalEActorAttrGloves)].toml: allEActorAttrInfo[int(internalEActorAttrGloves)],
	allEActorAttrInfo[int(internalEActorAttrJacket)].toml: allEActorAttrInfo[int(internalEActorAttrJacket)],
	allEActorAttrInfo[int(internalEActorAttrHead)].toml:   allEActorAttrInfo[int(internalEActorAttrHead)],
	allEActorAttrInfo[int(internalEActorAttrVest)].toml:   allEActorAttrInfo[int(internalEActorAttrVest)],
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *EActorAttr) UnmarshalText(data []byte) error {
	key := string(data)
	v, ok := mapCodeToEActorAttrInfo[key]
	if !ok {
		return fmt.Errorf("EActorAttr.UnmarshalText failed: invalid EActorAttr[%v]", key)
	}
	*e = v.value
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e EActorAttr) MarshalText() ([]byte, error) {
	return []byte(allEActorAttrInfo[e].toml), nil
}

func (e EActorAttr) String() string {
	return allEActorAttrInfo[e].toml
}

// EAmmunitionType EAmmunitionType
type EAmmunitionType int32

const (
	internalEAmmunitionTypeAmmoMagnum300 EAmmunitionType = EAmmunitionType(0) // 弹夹-.300马格兰
	internalEAmmunitionTypeAmmoACP45     EAmmunitionType = EAmmunitionType(1) // .45
	internalEAmmunitionTypeAmmoGauge12   EAmmunitionType = EAmmunitionType(2) // 12号口径
	internalEAmmunitionTypeAmmo5d56mm    EAmmunitionType = EAmmunitionType(3) // 5.56mm
	internalEAmmunitionTypeAmmo7d62mm    EAmmunitionType = EAmmunitionType(4) // 7.62mm
	internalEAmmunitionTypeAmmo9mm       EAmmunitionType = EAmmunitionType(5) // 9mm
)

type internalEAmmunitionTypeInfo struct {
	value EAmmunitionType
	toml  string
	desc  string
}

var allEAmmunitionTypeInfo = []*internalEAmmunitionTypeInfo{
	&internalEAmmunitionTypeInfo{
		value: internalEAmmunitionTypeAmmoMagnum300,
		toml:  "EAmmunitionType.AmmoMagnum300",
		desc:  "弹夹-.300马格兰",
	},
	&internalEAmmunitionTypeInfo{
		value: internalEAmmunitionTypeAmmoACP45,
		toml:  "EAmmunitionType.AmmoACP45",
		desc:  ".45",
	},
	&internalEAmmunitionTypeInfo{
		value: internalEAmmunitionTypeAmmoGauge12,
		toml:  "EAmmunitionType.AmmoGauge12",
		desc:  "12号口径",
	},
	&internalEAmmunitionTypeInfo{
		value: internalEAmmunitionTypeAmmo5d56mm,
		toml:  "EAmmunitionType.Ammo5d56mm",
		desc:  "5.56mm",
	},
	&internalEAmmunitionTypeInfo{
		value: internalEAmmunitionTypeAmmo7d62mm,
		toml:  "EAmmunitionType.Ammo7d62mm",
		desc:  "7.62mm",
	},
	&internalEAmmunitionTypeInfo{
		value: internalEAmmunitionTypeAmmo9mm,
		toml:  "EAmmunitionType.Ammo9mm",
		desc:  "9mm",
	},
}

var mapCodeToEAmmunitionTypeInfo = map[string]*internalEAmmunitionTypeInfo{
	allEAmmunitionTypeInfo[int(internalEAmmunitionTypeAmmoMagnum300)].toml: allEAmmunitionTypeInfo[int(internalEAmmunitionTypeAmmoMagnum300)],
	allEAmmunitionTypeInfo[int(internalEAmmunitionTypeAmmoACP45)].toml:     allEAmmunitionTypeInfo[int(internalEAmmunitionTypeAmmoACP45)],
	allEAmmunitionTypeInfo[int(internalEAmmunitionTypeAmmoGauge12)].toml:   allEAmmunitionTypeInfo[int(internalEAmmunitionTypeAmmoGauge12)],
	allEAmmunitionTypeInfo[int(internalEAmmunitionTypeAmmo5d56mm)].toml:    allEAmmunitionTypeInfo[int(internalEAmmunitionTypeAmmo5d56mm)],
	allEAmmunitionTypeInfo[int(internalEAmmunitionTypeAmmo7d62mm)].toml:    allEAmmunitionTypeInfo[int(internalEAmmunitionTypeAmmo7d62mm)],
	allEAmmunitionTypeInfo[int(internalEAmmunitionTypeAmmo9mm)].toml:       allEAmmunitionTypeInfo[int(internalEAmmunitionTypeAmmo9mm)],
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *EAmmunitionType) UnmarshalText(data []byte) error {
	key := string(data)
	v, ok := mapCodeToEAmmunitionTypeInfo[key]
	if !ok {
		return fmt.Errorf("EAmmunitionType.UnmarshalText failed: invalid EAmmunitionType[%v]", key)
	}
	*e = v.value
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e EAmmunitionType) MarshalText() ([]byte, error) {
	return []byte(allEAmmunitionTypeInfo[e].toml), nil
}

func (e EAmmunitionType) String() string {
	return allEAmmunitionTypeInfo[e].toml
}

// EAttachmentType EAttachmentType
type EAttachmentType int32

const (
	internalEAttachmentTypeMuzzleMod EAttachmentType = EAttachmentType(0) // 配件-枪口
	internalEAttachmentTypeLowerRail EAttachmentType = EAttachmentType(1) // 握把
	internalEAttachmentTypeUpperRail EAttachmentType = EAttachmentType(2) // 准镜
	internalEAttachmentTypeMagazine  EAttachmentType = EAttachmentType(3) // 弹夹
	internalEAttachmentTypeStock     EAttachmentType = EAttachmentType(4) // 枪托/子弹袋
)

type internalEAttachmentTypeInfo struct {
	value EAttachmentType
	toml  string
	desc  string
}

var allEAttachmentTypeInfo = []*internalEAttachmentTypeInfo{
	&internalEAttachmentTypeInfo{
		value: internalEAttachmentTypeMuzzleMod,
		toml:  "EAttachmentType.MuzzleMod",
		desc:  "配件-枪口",
	},
	&internalEAttachmentTypeInfo{
		value: internalEAttachmentTypeLowerRail,
		toml:  "EAttachmentType.LowerRail",
		desc:  "握把",
	},
	&internalEAttachmentTypeInfo{
		value: internalEAttachmentTypeUpperRail,
		toml:  "EAttachmentType.UpperRail",
		desc:  "准镜",
	},
	&internalEAttachmentTypeInfo{
		value: internalEAttachmentTypeMagazine,
		toml:  "EAttachmentType.Magazine",
		desc:  "弹夹",
	},
	&internalEAttachmentTypeInfo{
		value: internalEAttachmentTypeStock,
		toml:  "EAttachmentType.Stock",
		desc:  "枪托/子弹袋",
	},
}

var mapCodeToEAttachmentTypeInfo = map[string]*internalEAttachmentTypeInfo{
	allEAttachmentTypeInfo[int(internalEAttachmentTypeMuzzleMod)].toml: allEAttachmentTypeInfo[int(internalEAttachmentTypeMuzzleMod)],
	allEAttachmentTypeInfo[int(internalEAttachmentTypeLowerRail)].toml: allEAttachmentTypeInfo[int(internalEAttachmentTypeLowerRail)],
	allEAttachmentTypeInfo[int(internalEAttachmentTypeUpperRail)].toml: allEAttachmentTypeInfo[int(internalEAttachmentTypeUpperRail)],
	allEAttachmentTypeInfo[int(internalEAttachmentTypeMagazine)].toml:  allEAttachmentTypeInfo[int(internalEAttachmentTypeMagazine)],
	allEAttachmentTypeInfo[int(internalEAttachmentTypeStock)].toml:     allEAttachmentTypeInfo[int(internalEAttachmentTypeStock)],
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *EAttachmentType) UnmarshalText(data []byte) error {
	key := string(data)
	v, ok := mapCodeToEAttachmentTypeInfo[key]
	if !ok {
		return fmt.Errorf("EAttachmentType.UnmarshalText failed: invalid EAttachmentType[%v]", key)
	}
	*e = v.value
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e EAttachmentType) MarshalText() ([]byte, error) {
	return []byte(allEAttachmentTypeInfo[e].toml), nil
}

func (e EAttachmentType) String() string {
	return allEAttachmentTypeInfo[e].toml
}

// EAttr EAttr
type EAttr int32

const (
	internalEAttrHealth               EAttr = EAttr(0)  // 属性-血量
	internalEAttrBodyAttack           EAttr = EAttr(1)  // 身体攻击
	internalEAttrHeadAttack           EAttr = EAttr(2)  // 头部攻击
	internalEAttrBodyDefence          EAttr = EAttr(3)  // 身体防御
	internalEAttrHeadDefence          EAttr = EAttr(4)  // 头部防御
	internalEAttrDurable              EAttr = EAttr(5)  // 耐久
	internalEAttrClip                 EAttr = EAttr(6)  // 弹夹容量
	internalEAttrReload               EAttr = EAttr(7)  // 换弹时间-单位毫秒
	internalEAttrReloadReduce         EAttr = EAttr(8)  // 换弹时间缩短比例-百分比
	internalEAttrShootInterval        EAttr = EAttr(9)  // 射击间隔-单位毫秒
	internalEAttrShootIntervalReduce  EAttr = EAttr(10) // 射击间隔缩短比例-百分比
	internalEAttrHorPitch             EAttr = EAttr(11) // 水平准心跳动-单位0.001
	internalEAttrHorPitchReduce       EAttr = EAttr(12) // 水平准心跳动降低-百分比
	internalEAttrVerPitch             EAttr = EAttr(13) // 垂直准心跳动-单位0.001
	internalEAttrVerPitchReduce       EAttr = EAttr(14) // 垂直准心跳动降低-百分比
	internalEAttrPitchGoTime          EAttr = EAttr(15) // 准心跳动偏离时间-单位毫秒
	internalEAttrPitchBackTime        EAttr = EAttr(16) // 准心跳动恢复时间-单位毫秒
	internalEAttrCrossRangeMin        EAttr = EAttr(17) // 最小准心范围-单位0.001
	internalEAttrCrossRangeMax        EAttr = EAttr(18) // 最大准心范围-单位0.001
	internalEAttrCrossRange           EAttr = EAttr(19) // 准心扩散-单位0.001
	internalEAttrCrossOutTime         EAttr = EAttr(20) // 准心扩散时间-单位毫秒
	internalEAttrCrossInTime          EAttr = EAttr(21) // 准心收缩时间-单位毫秒
	internalEAttrCrossReduce          EAttr = EAttr(22) // 准心扩散降低-百分比
	internalEAttrCameraShakeRange     EAttr = EAttr(23) // 摄像机抖动位移-单位0.001
	internalEAttrCameraShakeTime      EAttr = EAttr(24) // 摄像机抖动时间-单位毫秒
	internalEAttrCameraShakeReduce    EAttr = EAttr(25) // 摄像机抖动降低-百分比
	internalEAttrShotRadius           EAttr = EAttr(26) // 散弹半径-单位0.001
	internalEAttrShotRadiusReduce     EAttr = EAttr(27) // 散弹半径降低-百分比
	internalEAttrMirrorMultiple       EAttr = EAttr(28) // 准镜倍数-单位0.001
	internalEAttrOpenMirrorTime       EAttr = EAttr(29) // 开镜时间-单位毫秒
	internalEAttrOpenMirrorTimeReduce EAttr = EAttr(30) // 开镜时间降低-百分比
)

type internalEAttrInfo struct {
	value EAttr
	toml  string
	desc  string
}

var allEAttrInfo = []*internalEAttrInfo{
	&internalEAttrInfo{
		value: internalEAttrHealth,
		toml:  "EAttr.Health",
		desc:  "属性-血量",
	},
	&internalEAttrInfo{
		value: internalEAttrBodyAttack,
		toml:  "EAttr.BodyAttack",
		desc:  "身体攻击",
	},
	&internalEAttrInfo{
		value: internalEAttrHeadAttack,
		toml:  "EAttr.HeadAttack",
		desc:  "头部攻击",
	},
	&internalEAttrInfo{
		value: internalEAttrBodyDefence,
		toml:  "EAttr.BodyDefence",
		desc:  "身体防御",
	},
	&internalEAttrInfo{
		value: internalEAttrHeadDefence,
		toml:  "EAttr.HeadDefence",
		desc:  "头部防御",
	},
	&internalEAttrInfo{
		value: internalEAttrDurable,
		toml:  "EAttr.Durable",
		desc:  "耐久",
	},
	&internalEAttrInfo{
		value: internalEAttrClip,
		toml:  "EAttr.Clip",
		desc:  "弹夹容量",
	},
	&internalEAttrInfo{
		value: internalEAttrReload,
		toml:  "EAttr.Reload",
		desc:  "换弹时间-单位毫秒",
	},
	&internalEAttrInfo{
		value: internalEAttrReloadReduce,
		toml:  "EAttr.ReloadReduce",
		desc:  "换弹时间缩短比例-百分比",
	},
	&internalEAttrInfo{
		value: internalEAttrShootInterval,
		toml:  "EAttr.ShootInterval",
		desc:  "射击间隔-单位毫秒",
	},
	&internalEAttrInfo{
		value: internalEAttrShootIntervalReduce,
		toml:  "EAttr.ShootIntervalReduce",
		desc:  "射击间隔缩短比例-百分比",
	},
	&internalEAttrInfo{
		value: internalEAttrHorPitch,
		toml:  "EAttr.HorPitch",
		desc:  "水平准心跳动-单位0.001",
	},
	&internalEAttrInfo{
		value: internalEAttrHorPitchReduce,
		toml:  "EAttr.HorPitchReduce",
		desc:  "水平准心跳动降低-百分比",
	},
	&internalEAttrInfo{
		value: internalEAttrVerPitch,
		toml:  "EAttr.VerPitch",
		desc:  "垂直准心跳动-单位0.001",
	},
	&internalEAttrInfo{
		value: internalEAttrVerPitchReduce,
		toml:  "EAttr.VerPitchReduce",
		desc:  "垂直准心跳动降低-百分比",
	},
	&internalEAttrInfo{
		value: internalEAttrPitchGoTime,
		toml:  "EAttr.PitchGoTime",
		desc:  "准心跳动偏离时间-单位毫秒",
	},
	&internalEAttrInfo{
		value: internalEAttrPitchBackTime,
		toml:  "EAttr.PitchBackTime",
		desc:  "准心跳动恢复时间-单位毫秒",
	},
	&internalEAttrInfo{
		value: internalEAttrCrossRangeMin,
		toml:  "EAttr.CrossRangeMin",
		desc:  "最小准心范围-单位0.001",
	},
	&internalEAttrInfo{
		value: internalEAttrCrossRangeMax,
		toml:  "EAttr.CrossRangeMax",
		desc:  "最大准心范围-单位0.001",
	},
	&internalEAttrInfo{
		value: internalEAttrCrossRange,
		toml:  "EAttr.CrossRange",
		desc:  "准心扩散-单位0.001",
	},
	&internalEAttrInfo{
		value: internalEAttrCrossOutTime,
		toml:  "EAttr.CrossOutTime",
		desc:  "准心扩散时间-单位毫秒",
	},
	&internalEAttrInfo{
		value: internalEAttrCrossInTime,
		toml:  "EAttr.CrossInTime",
		desc:  "准心收缩时间-单位毫秒",
	},
	&internalEAttrInfo{
		value: internalEAttrCrossReduce,
		toml:  "EAttr.CrossReduce",
		desc:  "准心扩散降低-百分比",
	},
	&internalEAttrInfo{
		value: internalEAttrCameraShakeRange,
		toml:  "EAttr.CameraShakeRange",
		desc:  "摄像机抖动位移-单位0.001",
	},
	&internalEAttrInfo{
		value: internalEAttrCameraShakeTime,
		toml:  "EAttr.CameraShakeTime",
		desc:  "摄像机抖动时间-单位毫秒",
	},
	&internalEAttrInfo{
		value: internalEAttrCameraShakeReduce,
		toml:  "EAttr.CameraShakeReduce",
		desc:  "摄像机抖动降低-百分比",
	},
	&internalEAttrInfo{
		value: internalEAttrShotRadius,
		toml:  "EAttr.ShotRadius",
		desc:  "散弹半径-单位0.001",
	},
	&internalEAttrInfo{
		value: internalEAttrShotRadiusReduce,
		toml:  "EAttr.ShotRadiusReduce",
		desc:  "散弹半径降低-百分比",
	},
	&internalEAttrInfo{
		value: internalEAttrMirrorMultiple,
		toml:  "EAttr.MirrorMultiple",
		desc:  "准镜倍数-单位0.001",
	},
	&internalEAttrInfo{
		value: internalEAttrOpenMirrorTime,
		toml:  "EAttr.OpenMirrorTime",
		desc:  "开镜时间-单位毫秒",
	},
	&internalEAttrInfo{
		value: internalEAttrOpenMirrorTimeReduce,
		toml:  "EAttr.OpenMirrorTimeReduce",
		desc:  "开镜时间降低-百分比",
	},
}

var mapCodeToEAttrInfo = map[string]*internalEAttrInfo{
	allEAttrInfo[int(internalEAttrHealth)].toml:               allEAttrInfo[int(internalEAttrHealth)],
	allEAttrInfo[int(internalEAttrBodyAttack)].toml:           allEAttrInfo[int(internalEAttrBodyAttack)],
	allEAttrInfo[int(internalEAttrHeadAttack)].toml:           allEAttrInfo[int(internalEAttrHeadAttack)],
	allEAttrInfo[int(internalEAttrBodyDefence)].toml:          allEAttrInfo[int(internalEAttrBodyDefence)],
	allEAttrInfo[int(internalEAttrHeadDefence)].toml:          allEAttrInfo[int(internalEAttrHeadDefence)],
	allEAttrInfo[int(internalEAttrDurable)].toml:              allEAttrInfo[int(internalEAttrDurable)],
	allEAttrInfo[int(internalEAttrClip)].toml:                 allEAttrInfo[int(internalEAttrClip)],
	allEAttrInfo[int(internalEAttrReload)].toml:               allEAttrInfo[int(internalEAttrReload)],
	allEAttrInfo[int(internalEAttrReloadReduce)].toml:         allEAttrInfo[int(internalEAttrReloadReduce)],
	allEAttrInfo[int(internalEAttrShootInterval)].toml:        allEAttrInfo[int(internalEAttrShootInterval)],
	allEAttrInfo[int(internalEAttrShootIntervalReduce)].toml:  allEAttrInfo[int(internalEAttrShootIntervalReduce)],
	allEAttrInfo[int(internalEAttrHorPitch)].toml:             allEAttrInfo[int(internalEAttrHorPitch)],
	allEAttrInfo[int(internalEAttrHorPitchReduce)].toml:       allEAttrInfo[int(internalEAttrHorPitchReduce)],
	allEAttrInfo[int(internalEAttrVerPitch)].toml:             allEAttrInfo[int(internalEAttrVerPitch)],
	allEAttrInfo[int(internalEAttrVerPitchReduce)].toml:       allEAttrInfo[int(internalEAttrVerPitchReduce)],
	allEAttrInfo[int(internalEAttrPitchGoTime)].toml:          allEAttrInfo[int(internalEAttrPitchGoTime)],
	allEAttrInfo[int(internalEAttrPitchBackTime)].toml:        allEAttrInfo[int(internalEAttrPitchBackTime)],
	allEAttrInfo[int(internalEAttrCrossRangeMin)].toml:        allEAttrInfo[int(internalEAttrCrossRangeMin)],
	allEAttrInfo[int(internalEAttrCrossRangeMax)].toml:        allEAttrInfo[int(internalEAttrCrossRangeMax)],
	allEAttrInfo[int(internalEAttrCrossRange)].toml:           allEAttrInfo[int(internalEAttrCrossRange)],
	allEAttrInfo[int(internalEAttrCrossOutTime)].toml:         allEAttrInfo[int(internalEAttrCrossOutTime)],
	allEAttrInfo[int(internalEAttrCrossInTime)].toml:          allEAttrInfo[int(internalEAttrCrossInTime)],
	allEAttrInfo[int(internalEAttrCrossReduce)].toml:          allEAttrInfo[int(internalEAttrCrossReduce)],
	allEAttrInfo[int(internalEAttrCameraShakeRange)].toml:     allEAttrInfo[int(internalEAttrCameraShakeRange)],
	allEAttrInfo[int(internalEAttrCameraShakeTime)].toml:      allEAttrInfo[int(internalEAttrCameraShakeTime)],
	allEAttrInfo[int(internalEAttrCameraShakeReduce)].toml:    allEAttrInfo[int(internalEAttrCameraShakeReduce)],
	allEAttrInfo[int(internalEAttrShotRadius)].toml:           allEAttrInfo[int(internalEAttrShotRadius)],
	allEAttrInfo[int(internalEAttrShotRadiusReduce)].toml:     allEAttrInfo[int(internalEAttrShotRadiusReduce)],
	allEAttrInfo[int(internalEAttrMirrorMultiple)].toml:       allEAttrInfo[int(internalEAttrMirrorMultiple)],
	allEAttrInfo[int(internalEAttrOpenMirrorTime)].toml:       allEAttrInfo[int(internalEAttrOpenMirrorTime)],
	allEAttrInfo[int(internalEAttrOpenMirrorTimeReduce)].toml: allEAttrInfo[int(internalEAttrOpenMirrorTimeReduce)],
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *EAttr) UnmarshalText(data []byte) error {
	key := string(data)
	v, ok := mapCodeToEAttrInfo[key]
	if !ok {
		return fmt.Errorf("EAttr.UnmarshalText failed: invalid EAttr[%v]", key)
	}
	*e = v.value
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e EAttr) MarshalText() ([]byte, error) {
	return []byte(allEAttrInfo[e].toml), nil
}

func (e EAttr) String() string {
	return allEAttrInfo[e].toml
}

// EConsumableType EConsumableType
type EConsumableType int32

const (
	internalEConsumableTypeAdrenalineSyringe EConsumableType = EConsumableType(0) // 战场消耗品-肾上腺素
	internalEConsumableTypeBandage           EConsumableType = EConsumableType(1) // 绷带
	internalEConsumableTypeEnergyDrink       EConsumableType = EConsumableType(2) // 能量饮料
	internalEConsumableTypeFirstAidKit       EConsumableType = EConsumableType(3) // 急救包
	internalEConsumableTypeMedKit            EConsumableType = EConsumableType(4) // 医疗箱
)

type internalEConsumableTypeInfo struct {
	value EConsumableType
	toml  string
	desc  string
}

var allEConsumableTypeInfo = []*internalEConsumableTypeInfo{
	&internalEConsumableTypeInfo{
		value: internalEConsumableTypeAdrenalineSyringe,
		toml:  "EConsumableType.AdrenalineSyringe",
		desc:  "战场消耗品-肾上腺素",
	},
	&internalEConsumableTypeInfo{
		value: internalEConsumableTypeBandage,
		toml:  "EConsumableType.Bandage",
		desc:  "绷带",
	},
	&internalEConsumableTypeInfo{
		value: internalEConsumableTypeEnergyDrink,
		toml:  "EConsumableType.EnergyDrink",
		desc:  "能量饮料",
	},
	&internalEConsumableTypeInfo{
		value: internalEConsumableTypeFirstAidKit,
		toml:  "EConsumableType.FirstAidKit",
		desc:  "急救包",
	},
	&internalEConsumableTypeInfo{
		value: internalEConsumableTypeMedKit,
		toml:  "EConsumableType.MedKit",
		desc:  "医疗箱",
	},
}

var mapCodeToEConsumableTypeInfo = map[string]*internalEConsumableTypeInfo{
	allEConsumableTypeInfo[int(internalEConsumableTypeAdrenalineSyringe)].toml: allEConsumableTypeInfo[int(internalEConsumableTypeAdrenalineSyringe)],
	allEConsumableTypeInfo[int(internalEConsumableTypeBandage)].toml:           allEConsumableTypeInfo[int(internalEConsumableTypeBandage)],
	allEConsumableTypeInfo[int(internalEConsumableTypeEnergyDrink)].toml:       allEConsumableTypeInfo[int(internalEConsumableTypeEnergyDrink)],
	allEConsumableTypeInfo[int(internalEConsumableTypeFirstAidKit)].toml:       allEConsumableTypeInfo[int(internalEConsumableTypeFirstAidKit)],
	allEConsumableTypeInfo[int(internalEConsumableTypeMedKit)].toml:            allEConsumableTypeInfo[int(internalEConsumableTypeMedKit)],
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *EConsumableType) UnmarshalText(data []byte) error {
	key := string(data)
	v, ok := mapCodeToEConsumableTypeInfo[key]
	if !ok {
		return fmt.Errorf("EConsumableType.UnmarshalText failed: invalid EConsumableType[%v]", key)
	}
	*e = v.value
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e EConsumableType) MarshalText() ([]byte, error) {
	return []byte(allEConsumableTypeInfo[e].toml), nil
}

func (e EConsumableType) String() string {
	return allEConsumableTypeInfo[e].toml
}

// EEquipmentType EEquipmentType
type EEquipmentType int32

const (
	internalEEquipmentTypeHelmet EEquipmentType = EEquipmentType(0) // 战场防具-头盔
	internalEEquipmentTypeVest   EEquipmentType = EEquipmentType(1) // 防弹衣
)

type internalEEquipmentTypeInfo struct {
	value EEquipmentType
	toml  string
	desc  string
}

var allEEquipmentTypeInfo = []*internalEEquipmentTypeInfo{
	&internalEEquipmentTypeInfo{
		value: internalEEquipmentTypeHelmet,
		toml:  "EEquipmentType.Helmet",
		desc:  "战场防具-头盔",
	},
	&internalEEquipmentTypeInfo{
		value: internalEEquipmentTypeVest,
		toml:  "EEquipmentType.Vest",
		desc:  "防弹衣",
	},
}

var mapCodeToEEquipmentTypeInfo = map[string]*internalEEquipmentTypeInfo{
	allEEquipmentTypeInfo[int(internalEEquipmentTypeHelmet)].toml: allEEquipmentTypeInfo[int(internalEEquipmentTypeHelmet)],
	allEEquipmentTypeInfo[int(internalEEquipmentTypeVest)].toml:   allEEquipmentTypeInfo[int(internalEEquipmentTypeVest)],
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *EEquipmentType) UnmarshalText(data []byte) error {
	key := string(data)
	v, ok := mapCodeToEEquipmentTypeInfo[key]
	if !ok {
		return fmt.Errorf("EEquipmentType.UnmarshalText failed: invalid EEquipmentType[%v]", key)
	}
	*e = v.value
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e EEquipmentType) MarshalText() ([]byte, error) {
	return []byte(allEEquipmentTypeInfo[e].toml), nil
}

func (e EEquipmentType) String() string {
	return allEEquipmentTypeInfo[e].toml
}

// EGunWeaponType EGunWeaponType
type EGunWeaponType int32

const (
	internalEGunWeaponTypeSniperRifle             EGunWeaponType = EGunWeaponType(0) // 枪械-狙击枪-awm
	internalEGunWeaponTypePistol                  EGunWeaponType = EGunWeaponType(1) // 手枪-p1911
	internalEGunWeaponTypeSubmachineGun           EGunWeaponType = EGunWeaponType(2) // 冲锋枪-ump9
	internalEGunWeaponTypeLightMachineGun         EGunWeaponType = EGunWeaponType(3) // 轻机枪-m249
	internalEGunWeaponTypeAssaultRifle            EGunWeaponType = EGunWeaponType(4) // 自动步枪-akm
	internalEGunWeaponTypeShotgun                 EGunWeaponType = EGunWeaponType(5) // 霰弹枪-s1897
	internalEGunWeaponTypeDesignatedMarksmanRifle EGunWeaponType = EGunWeaponType(6) // 精确射击步枪-sks
)

type internalEGunWeaponTypeInfo struct {
	value EGunWeaponType
	toml  string
	desc  string
}

var allEGunWeaponTypeInfo = []*internalEGunWeaponTypeInfo{
	&internalEGunWeaponTypeInfo{
		value: internalEGunWeaponTypeSniperRifle,
		toml:  "EGunWeaponType.SniperRifle",
		desc:  "枪械-狙击枪-awm",
	},
	&internalEGunWeaponTypeInfo{
		value: internalEGunWeaponTypePistol,
		toml:  "EGunWeaponType.Pistol",
		desc:  "手枪-p1911",
	},
	&internalEGunWeaponTypeInfo{
		value: internalEGunWeaponTypeSubmachineGun,
		toml:  "EGunWeaponType.SubmachineGun",
		desc:  "冲锋枪-ump9",
	},
	&internalEGunWeaponTypeInfo{
		value: internalEGunWeaponTypeLightMachineGun,
		toml:  "EGunWeaponType.LightMachineGun",
		desc:  "轻机枪-m249",
	},
	&internalEGunWeaponTypeInfo{
		value: internalEGunWeaponTypeAssaultRifle,
		toml:  "EGunWeaponType.AssaultRifle",
		desc:  "自动步枪-akm",
	},
	&internalEGunWeaponTypeInfo{
		value: internalEGunWeaponTypeShotgun,
		toml:  "EGunWeaponType.Shotgun",
		desc:  "霰弹枪-s1897",
	},
	&internalEGunWeaponTypeInfo{
		value: internalEGunWeaponTypeDesignatedMarksmanRifle,
		toml:  "EGunWeaponType.DesignatedMarksmanRifle",
		desc:  "精确射击步枪-sks",
	},
}

var mapCodeToEGunWeaponTypeInfo = map[string]*internalEGunWeaponTypeInfo{
	allEGunWeaponTypeInfo[int(internalEGunWeaponTypeSniperRifle)].toml:             allEGunWeaponTypeInfo[int(internalEGunWeaponTypeSniperRifle)],
	allEGunWeaponTypeInfo[int(internalEGunWeaponTypePistol)].toml:                  allEGunWeaponTypeInfo[int(internalEGunWeaponTypePistol)],
	allEGunWeaponTypeInfo[int(internalEGunWeaponTypeSubmachineGun)].toml:           allEGunWeaponTypeInfo[int(internalEGunWeaponTypeSubmachineGun)],
	allEGunWeaponTypeInfo[int(internalEGunWeaponTypeLightMachineGun)].toml:         allEGunWeaponTypeInfo[int(internalEGunWeaponTypeLightMachineGun)],
	allEGunWeaponTypeInfo[int(internalEGunWeaponTypeAssaultRifle)].toml:            allEGunWeaponTypeInfo[int(internalEGunWeaponTypeAssaultRifle)],
	allEGunWeaponTypeInfo[int(internalEGunWeaponTypeShotgun)].toml:                 allEGunWeaponTypeInfo[int(internalEGunWeaponTypeShotgun)],
	allEGunWeaponTypeInfo[int(internalEGunWeaponTypeDesignatedMarksmanRifle)].toml: allEGunWeaponTypeInfo[int(internalEGunWeaponTypeDesignatedMarksmanRifle)],
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *EGunWeaponType) UnmarshalText(data []byte) error {
	key := string(data)
	v, ok := mapCodeToEGunWeaponTypeInfo[key]
	if !ok {
		return fmt.Errorf("EGunWeaponType.UnmarshalText failed: invalid EGunWeaponType[%v]", key)
	}
	*e = v.value
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e EGunWeaponType) MarshalText() ([]byte, error) {
	return []byte(allEGunWeaponTypeInfo[e].toml), nil
}

func (e EGunWeaponType) String() string {
	return allEGunWeaponTypeInfo[e].toml
}

// EItemType EItemType
type EItemType int32

const (
	internalEItemTypeGunWeapon    EItemType = EItemType(0) // 物品-战场枪械
	internalEItemTypeAmmunition   EItemType = EItemType(1) // 战场枪械弹药
	internalEItemTypeAttachment   EItemType = EItemType(2) // 战场枪械配件
	internalEItemTypeMelleeWeapon EItemType = EItemType(3) // 战场近战物理武器
	internalEItemTypeEquipment    EItemType = EItemType(4) // 战场防具
	internalEItemTypeConsumable   EItemType = EItemType(5) // 战场补给品
	internalEItemTypeThrowable    EItemType = EItemType(6) // 战场投掷物
)

type internalEItemTypeInfo struct {
	value EItemType
	toml  string
	desc  string
}

var allEItemTypeInfo = []*internalEItemTypeInfo{
	&internalEItemTypeInfo{
		value: internalEItemTypeGunWeapon,
		toml:  "EItemType.GunWeapon",
		desc:  "物品-战场枪械",
	},
	&internalEItemTypeInfo{
		value: internalEItemTypeAmmunition,
		toml:  "EItemType.Ammunition",
		desc:  "战场枪械弹药",
	},
	&internalEItemTypeInfo{
		value: internalEItemTypeAttachment,
		toml:  "EItemType.Attachment",
		desc:  "战场枪械配件",
	},
	&internalEItemTypeInfo{
		value: internalEItemTypeMelleeWeapon,
		toml:  "EItemType.MelleeWeapon",
		desc:  "战场近战物理武器",
	},
	&internalEItemTypeInfo{
		value: internalEItemTypeEquipment,
		toml:  "EItemType.Equipment",
		desc:  "战场防具",
	},
	&internalEItemTypeInfo{
		value: internalEItemTypeConsumable,
		toml:  "EItemType.Consumable",
		desc:  "战场补给品",
	},
	&internalEItemTypeInfo{
		value: internalEItemTypeThrowable,
		toml:  "EItemType.Throwable",
		desc:  "战场投掷物",
	},
}

var mapCodeToEItemTypeInfo = map[string]*internalEItemTypeInfo{
	allEItemTypeInfo[int(internalEItemTypeGunWeapon)].toml:    allEItemTypeInfo[int(internalEItemTypeGunWeapon)],
	allEItemTypeInfo[int(internalEItemTypeAmmunition)].toml:   allEItemTypeInfo[int(internalEItemTypeAmmunition)],
	allEItemTypeInfo[int(internalEItemTypeAttachment)].toml:   allEItemTypeInfo[int(internalEItemTypeAttachment)],
	allEItemTypeInfo[int(internalEItemTypeMelleeWeapon)].toml: allEItemTypeInfo[int(internalEItemTypeMelleeWeapon)],
	allEItemTypeInfo[int(internalEItemTypeEquipment)].toml:    allEItemTypeInfo[int(internalEItemTypeEquipment)],
	allEItemTypeInfo[int(internalEItemTypeConsumable)].toml:   allEItemTypeInfo[int(internalEItemTypeConsumable)],
	allEItemTypeInfo[int(internalEItemTypeThrowable)].toml:    allEItemTypeInfo[int(internalEItemTypeThrowable)],
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *EItemType) UnmarshalText(data []byte) error {
	key := string(data)
	v, ok := mapCodeToEItemTypeInfo[key]
	if !ok {
		return fmt.Errorf("EItemType.UnmarshalText failed: invalid EItemType[%v]", key)
	}
	*e = v.value
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e EItemType) MarshalText() ([]byte, error) {
	return []byte(allEItemTypeInfo[e].toml), nil
}

func (e EItemType) String() string {
	return allEItemTypeInfo[e].toml
}

// EMelleeWeaponType EMelleeWeaponType
type EMelleeWeaponType int32

const (
	internalEMelleeWeaponTypeCrowbar EMelleeWeaponType = EMelleeWeaponType(0) // 战场近战物理武器-撬棒
	internalEMelleeWeaponTypePan     EMelleeWeaponType = EMelleeWeaponType(1) // 平底锅
	internalEMelleeWeaponTypeSickle  EMelleeWeaponType = EMelleeWeaponType(2) // 镰刀
	internalEMelleeWeaponTypeMachete EMelleeWeaponType = EMelleeWeaponType(3) // 砍刀
)

type internalEMelleeWeaponTypeInfo struct {
	value EMelleeWeaponType
	toml  string
	desc  string
}

var allEMelleeWeaponTypeInfo = []*internalEMelleeWeaponTypeInfo{
	&internalEMelleeWeaponTypeInfo{
		value: internalEMelleeWeaponTypeCrowbar,
		toml:  "EMelleeWeaponType.Crowbar",
		desc:  "战场近战物理武器-撬棒",
	},
	&internalEMelleeWeaponTypeInfo{
		value: internalEMelleeWeaponTypePan,
		toml:  "EMelleeWeaponType.Pan",
		desc:  "平底锅",
	},
	&internalEMelleeWeaponTypeInfo{
		value: internalEMelleeWeaponTypeSickle,
		toml:  "EMelleeWeaponType.Sickle",
		desc:  "镰刀",
	},
	&internalEMelleeWeaponTypeInfo{
		value: internalEMelleeWeaponTypeMachete,
		toml:  "EMelleeWeaponType.Machete",
		desc:  "砍刀",
	},
}

var mapCodeToEMelleeWeaponTypeInfo = map[string]*internalEMelleeWeaponTypeInfo{
	allEMelleeWeaponTypeInfo[int(internalEMelleeWeaponTypeCrowbar)].toml: allEMelleeWeaponTypeInfo[int(internalEMelleeWeaponTypeCrowbar)],
	allEMelleeWeaponTypeInfo[int(internalEMelleeWeaponTypePan)].toml:     allEMelleeWeaponTypeInfo[int(internalEMelleeWeaponTypePan)],
	allEMelleeWeaponTypeInfo[int(internalEMelleeWeaponTypeSickle)].toml:  allEMelleeWeaponTypeInfo[int(internalEMelleeWeaponTypeSickle)],
	allEMelleeWeaponTypeInfo[int(internalEMelleeWeaponTypeMachete)].toml: allEMelleeWeaponTypeInfo[int(internalEMelleeWeaponTypeMachete)],
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *EMelleeWeaponType) UnmarshalText(data []byte) error {
	key := string(data)
	v, ok := mapCodeToEMelleeWeaponTypeInfo[key]
	if !ok {
		return fmt.Errorf("EMelleeWeaponType.UnmarshalText failed: invalid EMelleeWeaponType[%v]", key)
	}
	*e = v.value
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e EMelleeWeaponType) MarshalText() ([]byte, error) {
	return []byte(allEMelleeWeaponTypeInfo[e].toml), nil
}

func (e EMelleeWeaponType) String() string {
	return allEMelleeWeaponTypeInfo[e].toml
}

// EShootMode EShootMode
type EShootMode int32

const (
	internalEShootModeAutoInterval EShootMode = EShootMode(0) // 射击模式-连续射击-时间间隔
	internalEShootModeAutoAction   EShootMode = EShootMode(1) // 连续射击-拉栓动作
	internalEShootModeAutoThree    EShootMode = EShootMode(2) // 三连射
)

type internalEShootModeInfo struct {
	value EShootMode
	toml  string
	desc  string
}

var allEShootModeInfo = []*internalEShootModeInfo{
	&internalEShootModeInfo{
		value: internalEShootModeAutoInterval,
		toml:  "EShootMode.AutoInterval",
		desc:  "射击模式-连续射击-时间间隔",
	},
	&internalEShootModeInfo{
		value: internalEShootModeAutoAction,
		toml:  "EShootMode.AutoAction",
		desc:  "连续射击-拉栓动作",
	},
	&internalEShootModeInfo{
		value: internalEShootModeAutoThree,
		toml:  "EShootMode.AutoThree",
		desc:  "三连射",
	},
}

var mapCodeToEShootModeInfo = map[string]*internalEShootModeInfo{
	allEShootModeInfo[int(internalEShootModeAutoInterval)].toml: allEShootModeInfo[int(internalEShootModeAutoInterval)],
	allEShootModeInfo[int(internalEShootModeAutoAction)].toml:   allEShootModeInfo[int(internalEShootModeAutoAction)],
	allEShootModeInfo[int(internalEShootModeAutoThree)].toml:    allEShootModeInfo[int(internalEShootModeAutoThree)],
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *EShootMode) UnmarshalText(data []byte) error {
	key := string(data)
	v, ok := mapCodeToEShootModeInfo[key]
	if !ok {
		return fmt.Errorf("EShootMode.UnmarshalText failed: invalid EShootMode[%v]", key)
	}
	*e = v.value
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e EShootMode) MarshalText() ([]byte, error) {
	return []byte(allEShootModeInfo[e].toml), nil
}

func (e EShootMode) String() string {
	return allEShootModeInfo[e].toml
}

// EThrowableType EThrowableType
type EThrowableType int32

const (
	internalEThrowableTypeFragGrenade     EThrowableType = EThrowableType(0) // 投掷物-破片手雷
	internalEThrowableTypeMolotovCocktail EThrowableType = EThrowableType(1) // 燃烧弹
	internalEThrowableTypeSmokeGrenade    EThrowableType = EThrowableType(2) // 烟雾弹
	internalEThrowableTypeStunGrenade     EThrowableType = EThrowableType(3) // 闪光弹
)

type internalEThrowableTypeInfo struct {
	value EThrowableType
	toml  string
	desc  string
}

var allEThrowableTypeInfo = []*internalEThrowableTypeInfo{
	&internalEThrowableTypeInfo{
		value: internalEThrowableTypeFragGrenade,
		toml:  "EThrowableType.FragGrenade",
		desc:  "投掷物-破片手雷",
	},
	&internalEThrowableTypeInfo{
		value: internalEThrowableTypeMolotovCocktail,
		toml:  "EThrowableType.MolotovCocktail",
		desc:  "燃烧弹",
	},
	&internalEThrowableTypeInfo{
		value: internalEThrowableTypeSmokeGrenade,
		toml:  "EThrowableType.SmokeGrenade",
		desc:  "烟雾弹",
	},
	&internalEThrowableTypeInfo{
		value: internalEThrowableTypeStunGrenade,
		toml:  "EThrowableType.StunGrenade",
		desc:  "闪光弹",
	},
}

var mapCodeToEThrowableTypeInfo = map[string]*internalEThrowableTypeInfo{
	allEThrowableTypeInfo[int(internalEThrowableTypeFragGrenade)].toml:     allEThrowableTypeInfo[int(internalEThrowableTypeFragGrenade)],
	allEThrowableTypeInfo[int(internalEThrowableTypeMolotovCocktail)].toml: allEThrowableTypeInfo[int(internalEThrowableTypeMolotovCocktail)],
	allEThrowableTypeInfo[int(internalEThrowableTypeSmokeGrenade)].toml:    allEThrowableTypeInfo[int(internalEThrowableTypeSmokeGrenade)],
	allEThrowableTypeInfo[int(internalEThrowableTypeStunGrenade)].toml:     allEThrowableTypeInfo[int(internalEThrowableTypeStunGrenade)],
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *EThrowableType) UnmarshalText(data []byte) error {
	key := string(data)
	v, ok := mapCodeToEThrowableTypeInfo[key]
	if !ok {
		return fmt.Errorf("EThrowableType.UnmarshalText failed: invalid EThrowableType[%v]", key)
	}
	*e = v.value
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e EThrowableType) MarshalText() ([]byte, error) {
	return []byte(allEThrowableTypeInfo[e].toml), nil
}

func (e EThrowableType) String() string {
	return allEThrowableTypeInfo[e].toml
}
