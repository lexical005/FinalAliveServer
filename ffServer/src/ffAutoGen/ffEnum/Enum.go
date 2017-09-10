package ffEnum

import (
	"fmt"
)

// EActorAttr EActorAttr
type EActorAttr int32

const (
	// EActorAttrActor 角色外观属性-角色
	EActorAttrActor EActorAttr = EActorAttr(0)
	// EActorAttrMask 面具
	EActorAttrMask EActorAttr = EActorAttr(1)
	// EActorAttrPants 裤子
	EActorAttrPants EActorAttr = EActorAttr(2)
	// EActorAttrShoes 鞋子
	EActorAttrShoes EActorAttr = EActorAttr(3)
	// EActorAttrShirt 衬衫
	EActorAttrShirt EActorAttr = EActorAttr(4)
	// EActorAttrBelt 腰带
	EActorAttrBelt EActorAttr = EActorAttr(5)
	// EActorAttrGloves 手套
	EActorAttrGloves EActorAttr = EActorAttr(6)
	// EActorAttrJacket 外衣
	EActorAttrJacket EActorAttr = EActorAttr(7)
	// EActorAttrHead 帽子/头盔
	EActorAttrHead EActorAttr = EActorAttr(8)
	// EActorAttrVest 防弹衣
	EActorAttrVest EActorAttr = EActorAttr(9)
)

type internalEActorAttrInfo struct {
	value EActorAttr
	toml  string
	desc  string
}

var allEActorAttrInfo = []*internalEActorAttrInfo{
	&internalEActorAttrInfo{
		value: EActorAttrActor,
		toml:  "EActorAttr.Actor",
		desc:  "角色外观属性-角色",
	},
	&internalEActorAttrInfo{
		value: EActorAttrMask,
		toml:  "EActorAttr.Mask",
		desc:  "面具",
	},
	&internalEActorAttrInfo{
		value: EActorAttrPants,
		toml:  "EActorAttr.Pants",
		desc:  "裤子",
	},
	&internalEActorAttrInfo{
		value: EActorAttrShoes,
		toml:  "EActorAttr.Shoes",
		desc:  "鞋子",
	},
	&internalEActorAttrInfo{
		value: EActorAttrShirt,
		toml:  "EActorAttr.Shirt",
		desc:  "衬衫",
	},
	&internalEActorAttrInfo{
		value: EActorAttrBelt,
		toml:  "EActorAttr.Belt",
		desc:  "腰带",
	},
	&internalEActorAttrInfo{
		value: EActorAttrGloves,
		toml:  "EActorAttr.Gloves",
		desc:  "手套",
	},
	&internalEActorAttrInfo{
		value: EActorAttrJacket,
		toml:  "EActorAttr.Jacket",
		desc:  "外衣",
	},
	&internalEActorAttrInfo{
		value: EActorAttrHead,
		toml:  "EActorAttr.Head",
		desc:  "帽子/头盔",
	},
	&internalEActorAttrInfo{
		value: EActorAttrVest,
		toml:  "EActorAttr.Vest",
		desc:  "防弹衣",
	},
}

var mapCodeToEActorAttrInfo = map[string]*internalEActorAttrInfo{
	allEActorAttrInfo[int32(EActorAttrActor)].toml:  allEActorAttrInfo[int(EActorAttrActor)],
	allEActorAttrInfo[int32(EActorAttrMask)].toml:   allEActorAttrInfo[int(EActorAttrMask)],
	allEActorAttrInfo[int32(EActorAttrPants)].toml:  allEActorAttrInfo[int(EActorAttrPants)],
	allEActorAttrInfo[int32(EActorAttrShoes)].toml:  allEActorAttrInfo[int(EActorAttrShoes)],
	allEActorAttrInfo[int32(EActorAttrShirt)].toml:  allEActorAttrInfo[int(EActorAttrShirt)],
	allEActorAttrInfo[int32(EActorAttrBelt)].toml:   allEActorAttrInfo[int(EActorAttrBelt)],
	allEActorAttrInfo[int32(EActorAttrGloves)].toml: allEActorAttrInfo[int(EActorAttrGloves)],
	allEActorAttrInfo[int32(EActorAttrJacket)].toml: allEActorAttrInfo[int(EActorAttrJacket)],
	allEActorAttrInfo[int32(EActorAttrHead)].toml:   allEActorAttrInfo[int(EActorAttrHead)],
	allEActorAttrInfo[int32(EActorAttrVest)].toml:   allEActorAttrInfo[int(EActorAttrVest)],
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
	// EAmmunitionTypeAmmoMagnum300 弹夹-.300马格兰
	EAmmunitionTypeAmmoMagnum300 EAmmunitionType = EAmmunitionType(0)
	// EAmmunitionTypeAmmoACP45 .45
	EAmmunitionTypeAmmoACP45 EAmmunitionType = EAmmunitionType(1)
	// EAmmunitionTypeAmmoGauge12 12号口径
	EAmmunitionTypeAmmoGauge12 EAmmunitionType = EAmmunitionType(2)
	// EAmmunitionTypeAmmo5d56mm 5.56mm
	EAmmunitionTypeAmmo5d56mm EAmmunitionType = EAmmunitionType(3)
	// EAmmunitionTypeAmmo7d62mm 7.62mm
	EAmmunitionTypeAmmo7d62mm EAmmunitionType = EAmmunitionType(4)
	// EAmmunitionTypeAmmo9mm 9mm
	EAmmunitionTypeAmmo9mm EAmmunitionType = EAmmunitionType(5)
)

type internalEAmmunitionTypeInfo struct {
	value EAmmunitionType
	toml  string
	desc  string
}

var allEAmmunitionTypeInfo = []*internalEAmmunitionTypeInfo{
	&internalEAmmunitionTypeInfo{
		value: EAmmunitionTypeAmmoMagnum300,
		toml:  "EAmmunitionType.AmmoMagnum300",
		desc:  "弹夹-.300马格兰",
	},
	&internalEAmmunitionTypeInfo{
		value: EAmmunitionTypeAmmoACP45,
		toml:  "EAmmunitionType.AmmoACP45",
		desc:  ".45",
	},
	&internalEAmmunitionTypeInfo{
		value: EAmmunitionTypeAmmoGauge12,
		toml:  "EAmmunitionType.AmmoGauge12",
		desc:  "12号口径",
	},
	&internalEAmmunitionTypeInfo{
		value: EAmmunitionTypeAmmo5d56mm,
		toml:  "EAmmunitionType.Ammo5d56mm",
		desc:  "5.56mm",
	},
	&internalEAmmunitionTypeInfo{
		value: EAmmunitionTypeAmmo7d62mm,
		toml:  "EAmmunitionType.Ammo7d62mm",
		desc:  "7.62mm",
	},
	&internalEAmmunitionTypeInfo{
		value: EAmmunitionTypeAmmo9mm,
		toml:  "EAmmunitionType.Ammo9mm",
		desc:  "9mm",
	},
}

var mapCodeToEAmmunitionTypeInfo = map[string]*internalEAmmunitionTypeInfo{
	allEAmmunitionTypeInfo[int32(EAmmunitionTypeAmmoMagnum300)].toml: allEAmmunitionTypeInfo[int(EAmmunitionTypeAmmoMagnum300)],
	allEAmmunitionTypeInfo[int32(EAmmunitionTypeAmmoACP45)].toml:     allEAmmunitionTypeInfo[int(EAmmunitionTypeAmmoACP45)],
	allEAmmunitionTypeInfo[int32(EAmmunitionTypeAmmoGauge12)].toml:   allEAmmunitionTypeInfo[int(EAmmunitionTypeAmmoGauge12)],
	allEAmmunitionTypeInfo[int32(EAmmunitionTypeAmmo5d56mm)].toml:    allEAmmunitionTypeInfo[int(EAmmunitionTypeAmmo5d56mm)],
	allEAmmunitionTypeInfo[int32(EAmmunitionTypeAmmo7d62mm)].toml:    allEAmmunitionTypeInfo[int(EAmmunitionTypeAmmo7d62mm)],
	allEAmmunitionTypeInfo[int32(EAmmunitionTypeAmmo9mm)].toml:       allEAmmunitionTypeInfo[int(EAmmunitionTypeAmmo9mm)],
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

// EArmorType EArmorType
type EArmorType int32

const (
	// EArmorTypeHelmet 战场防具-头盔
	EArmorTypeHelmet EArmorType = EArmorType(0)
	// EArmorTypeVest 防弹衣
	EArmorTypeVest EArmorType = EArmorType(1)
)

type internalEArmorTypeInfo struct {
	value EArmorType
	toml  string
	desc  string
}

var allEArmorTypeInfo = []*internalEArmorTypeInfo{
	&internalEArmorTypeInfo{
		value: EArmorTypeHelmet,
		toml:  "EArmorType.Helmet",
		desc:  "战场防具-头盔",
	},
	&internalEArmorTypeInfo{
		value: EArmorTypeVest,
		toml:  "EArmorType.Vest",
		desc:  "防弹衣",
	},
}

var mapCodeToEArmorTypeInfo = map[string]*internalEArmorTypeInfo{
	allEArmorTypeInfo[int32(EArmorTypeHelmet)].toml: allEArmorTypeInfo[int(EArmorTypeHelmet)],
	allEArmorTypeInfo[int32(EArmorTypeVest)].toml:   allEArmorTypeInfo[int(EArmorTypeVest)],
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *EArmorType) UnmarshalText(data []byte) error {
	key := string(data)
	v, ok := mapCodeToEArmorTypeInfo[key]
	if !ok {
		return fmt.Errorf("EArmorType.UnmarshalText failed: invalid EArmorType[%v]", key)
	}
	*e = v.value
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e EArmorType) MarshalText() ([]byte, error) {
	return []byte(allEArmorTypeInfo[e].toml), nil
}

func (e EArmorType) String() string {
	return allEArmorTypeInfo[e].toml
}

// EAttachmentType EAttachmentType
type EAttachmentType int32

const (
	// EAttachmentTypeMuzzleMod 配件-枪口
	EAttachmentTypeMuzzleMod EAttachmentType = EAttachmentType(0)
	// EAttachmentTypeLowerRail 握把
	EAttachmentTypeLowerRail EAttachmentType = EAttachmentType(1)
	// EAttachmentTypeUpperRail 准镜
	EAttachmentTypeUpperRail EAttachmentType = EAttachmentType(2)
	// EAttachmentTypeMagazine 弹夹
	EAttachmentTypeMagazine EAttachmentType = EAttachmentType(3)
	// EAttachmentTypeStock 枪托/子弹袋
	EAttachmentTypeStock EAttachmentType = EAttachmentType(4)
)

type internalEAttachmentTypeInfo struct {
	value EAttachmentType
	toml  string
	desc  string
}

var allEAttachmentTypeInfo = []*internalEAttachmentTypeInfo{
	&internalEAttachmentTypeInfo{
		value: EAttachmentTypeMuzzleMod,
		toml:  "EAttachmentType.MuzzleMod",
		desc:  "配件-枪口",
	},
	&internalEAttachmentTypeInfo{
		value: EAttachmentTypeLowerRail,
		toml:  "EAttachmentType.LowerRail",
		desc:  "握把",
	},
	&internalEAttachmentTypeInfo{
		value: EAttachmentTypeUpperRail,
		toml:  "EAttachmentType.UpperRail",
		desc:  "准镜",
	},
	&internalEAttachmentTypeInfo{
		value: EAttachmentTypeMagazine,
		toml:  "EAttachmentType.Magazine",
		desc:  "弹夹",
	},
	&internalEAttachmentTypeInfo{
		value: EAttachmentTypeStock,
		toml:  "EAttachmentType.Stock",
		desc:  "枪托/子弹袋",
	},
}

var mapCodeToEAttachmentTypeInfo = map[string]*internalEAttachmentTypeInfo{
	allEAttachmentTypeInfo[int32(EAttachmentTypeMuzzleMod)].toml: allEAttachmentTypeInfo[int(EAttachmentTypeMuzzleMod)],
	allEAttachmentTypeInfo[int32(EAttachmentTypeLowerRail)].toml: allEAttachmentTypeInfo[int(EAttachmentTypeLowerRail)],
	allEAttachmentTypeInfo[int32(EAttachmentTypeUpperRail)].toml: allEAttachmentTypeInfo[int(EAttachmentTypeUpperRail)],
	allEAttachmentTypeInfo[int32(EAttachmentTypeMagazine)].toml:  allEAttachmentTypeInfo[int(EAttachmentTypeMagazine)],
	allEAttachmentTypeInfo[int32(EAttachmentTypeStock)].toml:     allEAttachmentTypeInfo[int(EAttachmentTypeStock)],
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
	// EAttrHealth 属性-血量
	EAttrHealth EAttr = EAttr(0)
	// EAttrBodyAttack 身体攻击
	EAttrBodyAttack EAttr = EAttr(1)
	// EAttrHeadAttack 头部攻击
	EAttrHeadAttack EAttr = EAttr(2)
	// EAttrBodyDefence 身体防御
	EAttrBodyDefence EAttr = EAttr(3)
	// EAttrHeadDefence 头部防御
	EAttrHeadDefence EAttr = EAttr(4)
	// EAttrDurable 耐久
	EAttrDurable EAttr = EAttr(5)
	// EAttrClip 弹夹容量
	EAttrClip EAttr = EAttr(6)
	// EAttrReload 换弹时间-单位毫秒
	EAttrReload EAttr = EAttr(7)
	// EAttrReloadReduce 换弹时间缩短比例-百分比
	EAttrReloadReduce EAttr = EAttr(8)
	// EAttrShootInterval 射击间隔-单位毫秒
	EAttrShootInterval EAttr = EAttr(9)
	// EAttrShootIntervalReduce 射击间隔缩短比例-百分比
	EAttrShootIntervalReduce EAttr = EAttr(10)
	// EAttrHorPitch 水平准心跳动-单位0.001
	EAttrHorPitch EAttr = EAttr(11)
	// EAttrHorPitchReduce 水平准心跳动降低-百分比
	EAttrHorPitchReduce EAttr = EAttr(12)
	// EAttrVerPitch 垂直准心跳动-单位0.001
	EAttrVerPitch EAttr = EAttr(13)
	// EAttrVerPitchReduce 垂直准心跳动降低-百分比
	EAttrVerPitchReduce EAttr = EAttr(14)
	// EAttrPitchGoTime 准心跳动偏离时间-单位毫秒
	EAttrPitchGoTime EAttr = EAttr(15)
	// EAttrPitchBackTime 准心跳动恢复时间-单位毫秒
	EAttrPitchBackTime EAttr = EAttr(16)
	// EAttrCrossRangeMin 最小准心范围-单位0.001
	EAttrCrossRangeMin EAttr = EAttr(17)
	// EAttrCrossRangeMax 最大准心范围-单位0.001
	EAttrCrossRangeMax EAttr = EAttr(18)
	// EAttrCrossMoveRangeMax 最大移动准心范围-单位0.001
	EAttrCrossMoveRangeMax EAttr = EAttr(19)
	// EAttrCrossRange 准心扩散-单位0.001
	EAttrCrossRange EAttr = EAttr(20)
	// EAttrCrossOutTime 准心扩散时间-单位毫秒
	EAttrCrossOutTime EAttr = EAttr(21)
	// EAttrCrossInTime 准心收缩时间-单位毫秒
	EAttrCrossInTime EAttr = EAttr(22)
	// EAttrCrossMoveOutTime 移动准心扩散时间
	EAttrCrossMoveOutTime EAttr = EAttr(23)
	// EAttrCrossReduce 准心扩散降低-百分比
	EAttrCrossReduce EAttr = EAttr(24)
	// EAttrCameraShakeRange 摄像机抖动位移-单位0.001
	EAttrCameraShakeRange EAttr = EAttr(25)
	// EAttrCameraShakeTime 摄像机抖动时间-单位毫秒
	EAttrCameraShakeTime EAttr = EAttr(26)
	// EAttrCameraShakeReduce 摄像机抖动降低-百分比
	EAttrCameraShakeReduce EAttr = EAttr(27)
	// EAttrShotRadius 散弹半径-单位0.001
	EAttrShotRadius EAttr = EAttr(28)
	// EAttrShotRadiusReduce 散弹半径降低-百分比
	EAttrShotRadiusReduce EAttr = EAttr(29)
	// EAttrMirrorMultiple 准镜倍数-单位0.001
	EAttrMirrorMultiple EAttr = EAttr(30)
	// EAttrOpenMirrorTime 开镜时间-单位毫秒
	EAttrOpenMirrorTime EAttr = EAttr(31)
	// EAttrOpenMirrorTimeReduce 开镜时间降低-百分比
	EAttrOpenMirrorTimeReduce EAttr = EAttr(32)
)

type internalEAttrInfo struct {
	value EAttr
	toml  string
	desc  string
}

var allEAttrInfo = []*internalEAttrInfo{
	&internalEAttrInfo{
		value: EAttrHealth,
		toml:  "EAttr.Health",
		desc:  "属性-血量",
	},
	&internalEAttrInfo{
		value: EAttrBodyAttack,
		toml:  "EAttr.BodyAttack",
		desc:  "身体攻击",
	},
	&internalEAttrInfo{
		value: EAttrHeadAttack,
		toml:  "EAttr.HeadAttack",
		desc:  "头部攻击",
	},
	&internalEAttrInfo{
		value: EAttrBodyDefence,
		toml:  "EAttr.BodyDefence",
		desc:  "身体防御",
	},
	&internalEAttrInfo{
		value: EAttrHeadDefence,
		toml:  "EAttr.HeadDefence",
		desc:  "头部防御",
	},
	&internalEAttrInfo{
		value: EAttrDurable,
		toml:  "EAttr.Durable",
		desc:  "耐久",
	},
	&internalEAttrInfo{
		value: EAttrClip,
		toml:  "EAttr.Clip",
		desc:  "弹夹容量",
	},
	&internalEAttrInfo{
		value: EAttrReload,
		toml:  "EAttr.Reload",
		desc:  "换弹时间-单位毫秒",
	},
	&internalEAttrInfo{
		value: EAttrReloadReduce,
		toml:  "EAttr.ReloadReduce",
		desc:  "换弹时间缩短比例-百分比",
	},
	&internalEAttrInfo{
		value: EAttrShootInterval,
		toml:  "EAttr.ShootInterval",
		desc:  "射击间隔-单位毫秒",
	},
	&internalEAttrInfo{
		value: EAttrShootIntervalReduce,
		toml:  "EAttr.ShootIntervalReduce",
		desc:  "射击间隔缩短比例-百分比",
	},
	&internalEAttrInfo{
		value: EAttrHorPitch,
		toml:  "EAttr.HorPitch",
		desc:  "水平准心跳动-单位0.001",
	},
	&internalEAttrInfo{
		value: EAttrHorPitchReduce,
		toml:  "EAttr.HorPitchReduce",
		desc:  "水平准心跳动降低-百分比",
	},
	&internalEAttrInfo{
		value: EAttrVerPitch,
		toml:  "EAttr.VerPitch",
		desc:  "垂直准心跳动-单位0.001",
	},
	&internalEAttrInfo{
		value: EAttrVerPitchReduce,
		toml:  "EAttr.VerPitchReduce",
		desc:  "垂直准心跳动降低-百分比",
	},
	&internalEAttrInfo{
		value: EAttrPitchGoTime,
		toml:  "EAttr.PitchGoTime",
		desc:  "准心跳动偏离时间-单位毫秒",
	},
	&internalEAttrInfo{
		value: EAttrPitchBackTime,
		toml:  "EAttr.PitchBackTime",
		desc:  "准心跳动恢复时间-单位毫秒",
	},
	&internalEAttrInfo{
		value: EAttrCrossRangeMin,
		toml:  "EAttr.CrossRangeMin",
		desc:  "最小准心范围-单位0.001",
	},
	&internalEAttrInfo{
		value: EAttrCrossRangeMax,
		toml:  "EAttr.CrossRangeMax",
		desc:  "最大准心范围-单位0.001",
	},
	&internalEAttrInfo{
		value: EAttrCrossMoveRangeMax,
		toml:  "EAttr.CrossMoveRangeMax",
		desc:  "最大移动准心范围-单位0.001",
	},
	&internalEAttrInfo{
		value: EAttrCrossRange,
		toml:  "EAttr.CrossRange",
		desc:  "准心扩散-单位0.001",
	},
	&internalEAttrInfo{
		value: EAttrCrossOutTime,
		toml:  "EAttr.CrossOutTime",
		desc:  "准心扩散时间-单位毫秒",
	},
	&internalEAttrInfo{
		value: EAttrCrossInTime,
		toml:  "EAttr.CrossInTime",
		desc:  "准心收缩时间-单位毫秒",
	},
	&internalEAttrInfo{
		value: EAttrCrossMoveOutTime,
		toml:  "EAttr.CrossMoveOutTime",
		desc:  "移动准心扩散时间",
	},
	&internalEAttrInfo{
		value: EAttrCrossReduce,
		toml:  "EAttr.CrossReduce",
		desc:  "准心扩散降低-百分比",
	},
	&internalEAttrInfo{
		value: EAttrCameraShakeRange,
		toml:  "EAttr.CameraShakeRange",
		desc:  "摄像机抖动位移-单位0.001",
	},
	&internalEAttrInfo{
		value: EAttrCameraShakeTime,
		toml:  "EAttr.CameraShakeTime",
		desc:  "摄像机抖动时间-单位毫秒",
	},
	&internalEAttrInfo{
		value: EAttrCameraShakeReduce,
		toml:  "EAttr.CameraShakeReduce",
		desc:  "摄像机抖动降低-百分比",
	},
	&internalEAttrInfo{
		value: EAttrShotRadius,
		toml:  "EAttr.ShotRadius",
		desc:  "散弹半径-单位0.001",
	},
	&internalEAttrInfo{
		value: EAttrShotRadiusReduce,
		toml:  "EAttr.ShotRadiusReduce",
		desc:  "散弹半径降低-百分比",
	},
	&internalEAttrInfo{
		value: EAttrMirrorMultiple,
		toml:  "EAttr.MirrorMultiple",
		desc:  "准镜倍数-单位0.001",
	},
	&internalEAttrInfo{
		value: EAttrOpenMirrorTime,
		toml:  "EAttr.OpenMirrorTime",
		desc:  "开镜时间-单位毫秒",
	},
	&internalEAttrInfo{
		value: EAttrOpenMirrorTimeReduce,
		toml:  "EAttr.OpenMirrorTimeReduce",
		desc:  "开镜时间降低-百分比",
	},
}

var mapCodeToEAttrInfo = map[string]*internalEAttrInfo{
	allEAttrInfo[int32(EAttrHealth)].toml:               allEAttrInfo[int(EAttrHealth)],
	allEAttrInfo[int32(EAttrBodyAttack)].toml:           allEAttrInfo[int(EAttrBodyAttack)],
	allEAttrInfo[int32(EAttrHeadAttack)].toml:           allEAttrInfo[int(EAttrHeadAttack)],
	allEAttrInfo[int32(EAttrBodyDefence)].toml:          allEAttrInfo[int(EAttrBodyDefence)],
	allEAttrInfo[int32(EAttrHeadDefence)].toml:          allEAttrInfo[int(EAttrHeadDefence)],
	allEAttrInfo[int32(EAttrDurable)].toml:              allEAttrInfo[int(EAttrDurable)],
	allEAttrInfo[int32(EAttrClip)].toml:                 allEAttrInfo[int(EAttrClip)],
	allEAttrInfo[int32(EAttrReload)].toml:               allEAttrInfo[int(EAttrReload)],
	allEAttrInfo[int32(EAttrReloadReduce)].toml:         allEAttrInfo[int(EAttrReloadReduce)],
	allEAttrInfo[int32(EAttrShootInterval)].toml:        allEAttrInfo[int(EAttrShootInterval)],
	allEAttrInfo[int32(EAttrShootIntervalReduce)].toml:  allEAttrInfo[int(EAttrShootIntervalReduce)],
	allEAttrInfo[int32(EAttrHorPitch)].toml:             allEAttrInfo[int(EAttrHorPitch)],
	allEAttrInfo[int32(EAttrHorPitchReduce)].toml:       allEAttrInfo[int(EAttrHorPitchReduce)],
	allEAttrInfo[int32(EAttrVerPitch)].toml:             allEAttrInfo[int(EAttrVerPitch)],
	allEAttrInfo[int32(EAttrVerPitchReduce)].toml:       allEAttrInfo[int(EAttrVerPitchReduce)],
	allEAttrInfo[int32(EAttrPitchGoTime)].toml:          allEAttrInfo[int(EAttrPitchGoTime)],
	allEAttrInfo[int32(EAttrPitchBackTime)].toml:        allEAttrInfo[int(EAttrPitchBackTime)],
	allEAttrInfo[int32(EAttrCrossRangeMin)].toml:        allEAttrInfo[int(EAttrCrossRangeMin)],
	allEAttrInfo[int32(EAttrCrossRangeMax)].toml:        allEAttrInfo[int(EAttrCrossRangeMax)],
	allEAttrInfo[int32(EAttrCrossMoveRangeMax)].toml:    allEAttrInfo[int(EAttrCrossMoveRangeMax)],
	allEAttrInfo[int32(EAttrCrossRange)].toml:           allEAttrInfo[int(EAttrCrossRange)],
	allEAttrInfo[int32(EAttrCrossOutTime)].toml:         allEAttrInfo[int(EAttrCrossOutTime)],
	allEAttrInfo[int32(EAttrCrossInTime)].toml:          allEAttrInfo[int(EAttrCrossInTime)],
	allEAttrInfo[int32(EAttrCrossMoveOutTime)].toml:     allEAttrInfo[int(EAttrCrossMoveOutTime)],
	allEAttrInfo[int32(EAttrCrossReduce)].toml:          allEAttrInfo[int(EAttrCrossReduce)],
	allEAttrInfo[int32(EAttrCameraShakeRange)].toml:     allEAttrInfo[int(EAttrCameraShakeRange)],
	allEAttrInfo[int32(EAttrCameraShakeTime)].toml:      allEAttrInfo[int(EAttrCameraShakeTime)],
	allEAttrInfo[int32(EAttrCameraShakeReduce)].toml:    allEAttrInfo[int(EAttrCameraShakeReduce)],
	allEAttrInfo[int32(EAttrShotRadius)].toml:           allEAttrInfo[int(EAttrShotRadius)],
	allEAttrInfo[int32(EAttrShotRadiusReduce)].toml:     allEAttrInfo[int(EAttrShotRadiusReduce)],
	allEAttrInfo[int32(EAttrMirrorMultiple)].toml:       allEAttrInfo[int(EAttrMirrorMultiple)],
	allEAttrInfo[int32(EAttrOpenMirrorTime)].toml:       allEAttrInfo[int(EAttrOpenMirrorTime)],
	allEAttrInfo[int32(EAttrOpenMirrorTimeReduce)].toml: allEAttrInfo[int(EAttrOpenMirrorTimeReduce)],
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

// EBornType EBornType
type EBornType int32

const (
	// EBornTypeRandRolePrepare 准备期间-角色出生点
	EBornTypeRandRolePrepare EBornType = EBornType(0)
	// EBornTypeRandItemPrepare 准备期间-物品出生点
	EBornTypeRandItemPrepare EBornType = EBornType(1)
	// EBornTypeRandItemBattle 战斗期间-物品出生点
	EBornTypeRandItemBattle EBornType = EBornType(2)
)

type internalEBornTypeInfo struct {
	value EBornType
	toml  string
	desc  string
}

var allEBornTypeInfo = []*internalEBornTypeInfo{
	&internalEBornTypeInfo{
		value: EBornTypeRandRolePrepare,
		toml:  "EBornType.RandRolePrepare",
		desc:  "准备期间-角色出生点",
	},
	&internalEBornTypeInfo{
		value: EBornTypeRandItemPrepare,
		toml:  "EBornType.RandItemPrepare",
		desc:  "准备期间-物品出生点",
	},
	&internalEBornTypeInfo{
		value: EBornTypeRandItemBattle,
		toml:  "EBornType.RandItemBattle",
		desc:  "战斗期间-物品出生点",
	},
}

var mapCodeToEBornTypeInfo = map[string]*internalEBornTypeInfo{
	allEBornTypeInfo[int32(EBornTypeRandRolePrepare)].toml: allEBornTypeInfo[int(EBornTypeRandRolePrepare)],
	allEBornTypeInfo[int32(EBornTypeRandItemPrepare)].toml: allEBornTypeInfo[int(EBornTypeRandItemPrepare)],
	allEBornTypeInfo[int32(EBornTypeRandItemBattle)].toml:  allEBornTypeInfo[int(EBornTypeRandItemBattle)],
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *EBornType) UnmarshalText(data []byte) error {
	key := string(data)
	v, ok := mapCodeToEBornTypeInfo[key]
	if !ok {
		return fmt.Errorf("EBornType.UnmarshalText failed: invalid EBornType[%v]", key)
	}
	*e = v.value
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e EBornType) MarshalText() ([]byte, error) {
	return []byte(allEBornTypeInfo[e].toml), nil
}

func (e EBornType) String() string {
	return allEBornTypeInfo[e].toml
}

// EConsumableType EConsumableType
type EConsumableType int32

const (
	// EConsumableTypeAdrenalineSyringe 战场消耗品-肾上腺素
	EConsumableTypeAdrenalineSyringe EConsumableType = EConsumableType(0)
	// EConsumableTypeBandage 绷带
	EConsumableTypeBandage EConsumableType = EConsumableType(1)
	// EConsumableTypeEnergyDrink 能量饮料
	EConsumableTypeEnergyDrink EConsumableType = EConsumableType(2)
	// EConsumableTypeFirstAidKit 急救包
	EConsumableTypeFirstAidKit EConsumableType = EConsumableType(3)
	// EConsumableTypeMedKit 医疗箱
	EConsumableTypeMedKit EConsumableType = EConsumableType(4)
)

type internalEConsumableTypeInfo struct {
	value EConsumableType
	toml  string
	desc  string
}

var allEConsumableTypeInfo = []*internalEConsumableTypeInfo{
	&internalEConsumableTypeInfo{
		value: EConsumableTypeAdrenalineSyringe,
		toml:  "EConsumableType.AdrenalineSyringe",
		desc:  "战场消耗品-肾上腺素",
	},
	&internalEConsumableTypeInfo{
		value: EConsumableTypeBandage,
		toml:  "EConsumableType.Bandage",
		desc:  "绷带",
	},
	&internalEConsumableTypeInfo{
		value: EConsumableTypeEnergyDrink,
		toml:  "EConsumableType.EnergyDrink",
		desc:  "能量饮料",
	},
	&internalEConsumableTypeInfo{
		value: EConsumableTypeFirstAidKit,
		toml:  "EConsumableType.FirstAidKit",
		desc:  "急救包",
	},
	&internalEConsumableTypeInfo{
		value: EConsumableTypeMedKit,
		toml:  "EConsumableType.MedKit",
		desc:  "医疗箱",
	},
}

var mapCodeToEConsumableTypeInfo = map[string]*internalEConsumableTypeInfo{
	allEConsumableTypeInfo[int32(EConsumableTypeAdrenalineSyringe)].toml: allEConsumableTypeInfo[int(EConsumableTypeAdrenalineSyringe)],
	allEConsumableTypeInfo[int32(EConsumableTypeBandage)].toml:           allEConsumableTypeInfo[int(EConsumableTypeBandage)],
	allEConsumableTypeInfo[int32(EConsumableTypeEnergyDrink)].toml:       allEConsumableTypeInfo[int(EConsumableTypeEnergyDrink)],
	allEConsumableTypeInfo[int32(EConsumableTypeFirstAidKit)].toml:       allEConsumableTypeInfo[int(EConsumableTypeFirstAidKit)],
	allEConsumableTypeInfo[int32(EConsumableTypeMedKit)].toml:            allEConsumableTypeInfo[int(EConsumableTypeMedKit)],
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

// EGunWeaponType EGunWeaponType
type EGunWeaponType int32

const (
	// EGunWeaponTypeSniperRifle 枪械-狙击枪-awm
	EGunWeaponTypeSniperRifle EGunWeaponType = EGunWeaponType(0)
	// EGunWeaponTypePistol 手枪-p1911
	EGunWeaponTypePistol EGunWeaponType = EGunWeaponType(1)
	// EGunWeaponTypeSubmachineGun 冲锋枪-ump9
	EGunWeaponTypeSubmachineGun EGunWeaponType = EGunWeaponType(2)
	// EGunWeaponTypeLightMachineGun 轻机枪-m249
	EGunWeaponTypeLightMachineGun EGunWeaponType = EGunWeaponType(3)
	// EGunWeaponTypeAssaultRifle 自动步枪-akm
	EGunWeaponTypeAssaultRifle EGunWeaponType = EGunWeaponType(4)
	// EGunWeaponTypeShotgun 霰弹枪-s1897
	EGunWeaponTypeShotgun EGunWeaponType = EGunWeaponType(5)
	// EGunWeaponTypeDesignatedMarksmanRifle 精确射击步枪-sks
	EGunWeaponTypeDesignatedMarksmanRifle EGunWeaponType = EGunWeaponType(6)
)

type internalEGunWeaponTypeInfo struct {
	value EGunWeaponType
	toml  string
	desc  string
}

var allEGunWeaponTypeInfo = []*internalEGunWeaponTypeInfo{
	&internalEGunWeaponTypeInfo{
		value: EGunWeaponTypeSniperRifle,
		toml:  "EGunWeaponType.SniperRifle",
		desc:  "枪械-狙击枪-awm",
	},
	&internalEGunWeaponTypeInfo{
		value: EGunWeaponTypePistol,
		toml:  "EGunWeaponType.Pistol",
		desc:  "手枪-p1911",
	},
	&internalEGunWeaponTypeInfo{
		value: EGunWeaponTypeSubmachineGun,
		toml:  "EGunWeaponType.SubmachineGun",
		desc:  "冲锋枪-ump9",
	},
	&internalEGunWeaponTypeInfo{
		value: EGunWeaponTypeLightMachineGun,
		toml:  "EGunWeaponType.LightMachineGun",
		desc:  "轻机枪-m249",
	},
	&internalEGunWeaponTypeInfo{
		value: EGunWeaponTypeAssaultRifle,
		toml:  "EGunWeaponType.AssaultRifle",
		desc:  "自动步枪-akm",
	},
	&internalEGunWeaponTypeInfo{
		value: EGunWeaponTypeShotgun,
		toml:  "EGunWeaponType.Shotgun",
		desc:  "霰弹枪-s1897",
	},
	&internalEGunWeaponTypeInfo{
		value: EGunWeaponTypeDesignatedMarksmanRifle,
		toml:  "EGunWeaponType.DesignatedMarksmanRifle",
		desc:  "精确射击步枪-sks",
	},
}

var mapCodeToEGunWeaponTypeInfo = map[string]*internalEGunWeaponTypeInfo{
	allEGunWeaponTypeInfo[int32(EGunWeaponTypeSniperRifle)].toml:             allEGunWeaponTypeInfo[int(EGunWeaponTypeSniperRifle)],
	allEGunWeaponTypeInfo[int32(EGunWeaponTypePistol)].toml:                  allEGunWeaponTypeInfo[int(EGunWeaponTypePistol)],
	allEGunWeaponTypeInfo[int32(EGunWeaponTypeSubmachineGun)].toml:           allEGunWeaponTypeInfo[int(EGunWeaponTypeSubmachineGun)],
	allEGunWeaponTypeInfo[int32(EGunWeaponTypeLightMachineGun)].toml:         allEGunWeaponTypeInfo[int(EGunWeaponTypeLightMachineGun)],
	allEGunWeaponTypeInfo[int32(EGunWeaponTypeAssaultRifle)].toml:            allEGunWeaponTypeInfo[int(EGunWeaponTypeAssaultRifle)],
	allEGunWeaponTypeInfo[int32(EGunWeaponTypeShotgun)].toml:                 allEGunWeaponTypeInfo[int(EGunWeaponTypeShotgun)],
	allEGunWeaponTypeInfo[int32(EGunWeaponTypeDesignatedMarksmanRifle)].toml: allEGunWeaponTypeInfo[int(EGunWeaponTypeDesignatedMarksmanRifle)],
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
	// EItemTypeGunWeapon 物品-战场枪械
	EItemTypeGunWeapon EItemType = EItemType(0)
	// EItemTypeAmmunition 战场枪械弹药
	EItemTypeAmmunition EItemType = EItemType(1)
	// EItemTypeAttachment 战场枪械配件
	EItemTypeAttachment EItemType = EItemType(2)
	// EItemTypeMelleeWeapon 战场近战物理武器
	EItemTypeMelleeWeapon EItemType = EItemType(3)
	// EItemTypeArmor 战场防具
	EItemTypeArmor EItemType = EItemType(4)
	// EItemTypeConsumable 战场补给品
	EItemTypeConsumable EItemType = EItemType(5)
	// EItemTypeThrowable 战场投掷物
	EItemTypeThrowable EItemType = EItemType(6)
	// EItemTypeRole 主角
	EItemTypeRole EItemType = EItemType(7)
	// EItemTypeBox 箱子
	EItemTypeBox EItemType = EItemType(8)
)

type internalEItemTypeInfo struct {
	value EItemType
	toml  string
	desc  string
}

var allEItemTypeInfo = []*internalEItemTypeInfo{
	&internalEItemTypeInfo{
		value: EItemTypeGunWeapon,
		toml:  "EItemType.GunWeapon",
		desc:  "物品-战场枪械",
	},
	&internalEItemTypeInfo{
		value: EItemTypeAmmunition,
		toml:  "EItemType.Ammunition",
		desc:  "战场枪械弹药",
	},
	&internalEItemTypeInfo{
		value: EItemTypeAttachment,
		toml:  "EItemType.Attachment",
		desc:  "战场枪械配件",
	},
	&internalEItemTypeInfo{
		value: EItemTypeMelleeWeapon,
		toml:  "EItemType.MelleeWeapon",
		desc:  "战场近战物理武器",
	},
	&internalEItemTypeInfo{
		value: EItemTypeArmor,
		toml:  "EItemType.Armor",
		desc:  "战场防具",
	},
	&internalEItemTypeInfo{
		value: EItemTypeConsumable,
		toml:  "EItemType.Consumable",
		desc:  "战场补给品",
	},
	&internalEItemTypeInfo{
		value: EItemTypeThrowable,
		toml:  "EItemType.Throwable",
		desc:  "战场投掷物",
	},
	&internalEItemTypeInfo{
		value: EItemTypeRole,
		toml:  "EItemType.Role",
		desc:  "主角",
	},
	&internalEItemTypeInfo{
		value: EItemTypeBox,
		toml:  "EItemType.Box",
		desc:  "箱子",
	},
}

var mapCodeToEItemTypeInfo = map[string]*internalEItemTypeInfo{
	allEItemTypeInfo[int32(EItemTypeGunWeapon)].toml:    allEItemTypeInfo[int(EItemTypeGunWeapon)],
	allEItemTypeInfo[int32(EItemTypeAmmunition)].toml:   allEItemTypeInfo[int(EItemTypeAmmunition)],
	allEItemTypeInfo[int32(EItemTypeAttachment)].toml:   allEItemTypeInfo[int(EItemTypeAttachment)],
	allEItemTypeInfo[int32(EItemTypeMelleeWeapon)].toml: allEItemTypeInfo[int(EItemTypeMelleeWeapon)],
	allEItemTypeInfo[int32(EItemTypeArmor)].toml:        allEItemTypeInfo[int(EItemTypeArmor)],
	allEItemTypeInfo[int32(EItemTypeConsumable)].toml:   allEItemTypeInfo[int(EItemTypeConsumable)],
	allEItemTypeInfo[int32(EItemTypeThrowable)].toml:    allEItemTypeInfo[int(EItemTypeThrowable)],
	allEItemTypeInfo[int32(EItemTypeRole)].toml:         allEItemTypeInfo[int(EItemTypeRole)],
	allEItemTypeInfo[int32(EItemTypeBox)].toml:          allEItemTypeInfo[int(EItemTypeBox)],
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
	// EMelleeWeaponTypeCrowbar 战场近战物理武器-撬棒
	EMelleeWeaponTypeCrowbar EMelleeWeaponType = EMelleeWeaponType(0)
	// EMelleeWeaponTypePan 平底锅
	EMelleeWeaponTypePan EMelleeWeaponType = EMelleeWeaponType(1)
	// EMelleeWeaponTypeSickle 镰刀
	EMelleeWeaponTypeSickle EMelleeWeaponType = EMelleeWeaponType(2)
	// EMelleeWeaponTypeMachete 砍刀
	EMelleeWeaponTypeMachete EMelleeWeaponType = EMelleeWeaponType(3)
)

type internalEMelleeWeaponTypeInfo struct {
	value EMelleeWeaponType
	toml  string
	desc  string
}

var allEMelleeWeaponTypeInfo = []*internalEMelleeWeaponTypeInfo{
	&internalEMelleeWeaponTypeInfo{
		value: EMelleeWeaponTypeCrowbar,
		toml:  "EMelleeWeaponType.Crowbar",
		desc:  "战场近战物理武器-撬棒",
	},
	&internalEMelleeWeaponTypeInfo{
		value: EMelleeWeaponTypePan,
		toml:  "EMelleeWeaponType.Pan",
		desc:  "平底锅",
	},
	&internalEMelleeWeaponTypeInfo{
		value: EMelleeWeaponTypeSickle,
		toml:  "EMelleeWeaponType.Sickle",
		desc:  "镰刀",
	},
	&internalEMelleeWeaponTypeInfo{
		value: EMelleeWeaponTypeMachete,
		toml:  "EMelleeWeaponType.Machete",
		desc:  "砍刀",
	},
}

var mapCodeToEMelleeWeaponTypeInfo = map[string]*internalEMelleeWeaponTypeInfo{
	allEMelleeWeaponTypeInfo[int32(EMelleeWeaponTypeCrowbar)].toml: allEMelleeWeaponTypeInfo[int(EMelleeWeaponTypeCrowbar)],
	allEMelleeWeaponTypeInfo[int32(EMelleeWeaponTypePan)].toml:     allEMelleeWeaponTypeInfo[int(EMelleeWeaponTypePan)],
	allEMelleeWeaponTypeInfo[int32(EMelleeWeaponTypeSickle)].toml:  allEMelleeWeaponTypeInfo[int(EMelleeWeaponTypeSickle)],
	allEMelleeWeaponTypeInfo[int32(EMelleeWeaponTypeMachete)].toml: allEMelleeWeaponTypeInfo[int(EMelleeWeaponTypeMachete)],
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

// ERoleAction ERoleAction
type ERoleAction int32

const (
	// ERoleActionSquat 战场角色动作行为-下蹲
	ERoleActionSquat ERoleAction = ERoleAction(0)
	// ERoleActionStand 站起
	ERoleActionStand ERoleAction = ERoleAction(1)
	// ERoleActionLying 匍匐
	ERoleActionLying ERoleAction = ERoleAction(2)
	// ERoleActionJump 跳跃
	ERoleActionJump ERoleAction = ERoleAction(3)
	// ERoleActionWeaponAim 瞄准
	ERoleActionWeaponAim ERoleAction = ERoleAction(4)
	// ERoleActionWeaponFire 开火
	ERoleActionWeaponFire ERoleAction = ERoleAction(5)
	// ERoleActionWeaponActive 武器激活
	ERoleActionWeaponActive ERoleAction = ERoleAction(6)
	// ERoleActionWeaponReload 武器换弹
	ERoleActionWeaponReload ERoleAction = ERoleAction(7)
	// ERoleActionThrow 投掷
	ERoleActionThrow ERoleAction = ERoleAction(8)
	// ERoleActionHeal 治疗
	ERoleActionHeal ERoleAction = ERoleAction(9)
)

type internalERoleActionInfo struct {
	value ERoleAction
	toml  string
	desc  string
}

var allERoleActionInfo = []*internalERoleActionInfo{
	&internalERoleActionInfo{
		value: ERoleActionSquat,
		toml:  "ERoleAction.Squat",
		desc:  "战场角色动作行为-下蹲",
	},
	&internalERoleActionInfo{
		value: ERoleActionStand,
		toml:  "ERoleAction.Stand",
		desc:  "站起",
	},
	&internalERoleActionInfo{
		value: ERoleActionLying,
		toml:  "ERoleAction.Lying",
		desc:  "匍匐",
	},
	&internalERoleActionInfo{
		value: ERoleActionJump,
		toml:  "ERoleAction.Jump",
		desc:  "跳跃",
	},
	&internalERoleActionInfo{
		value: ERoleActionWeaponAim,
		toml:  "ERoleAction.WeaponAim",
		desc:  "瞄准",
	},
	&internalERoleActionInfo{
		value: ERoleActionWeaponFire,
		toml:  "ERoleAction.WeaponFire",
		desc:  "开火",
	},
	&internalERoleActionInfo{
		value: ERoleActionWeaponActive,
		toml:  "ERoleAction.WeaponActive",
		desc:  "武器激活",
	},
	&internalERoleActionInfo{
		value: ERoleActionWeaponReload,
		toml:  "ERoleAction.WeaponReload",
		desc:  "武器换弹",
	},
	&internalERoleActionInfo{
		value: ERoleActionThrow,
		toml:  "ERoleAction.Throw",
		desc:  "投掷",
	},
	&internalERoleActionInfo{
		value: ERoleActionHeal,
		toml:  "ERoleAction.Heal",
		desc:  "治疗",
	},
}

var mapCodeToERoleActionInfo = map[string]*internalERoleActionInfo{
	allERoleActionInfo[int32(ERoleActionSquat)].toml:        allERoleActionInfo[int(ERoleActionSquat)],
	allERoleActionInfo[int32(ERoleActionStand)].toml:        allERoleActionInfo[int(ERoleActionStand)],
	allERoleActionInfo[int32(ERoleActionLying)].toml:        allERoleActionInfo[int(ERoleActionLying)],
	allERoleActionInfo[int32(ERoleActionJump)].toml:         allERoleActionInfo[int(ERoleActionJump)],
	allERoleActionInfo[int32(ERoleActionWeaponAim)].toml:    allERoleActionInfo[int(ERoleActionWeaponAim)],
	allERoleActionInfo[int32(ERoleActionWeaponFire)].toml:   allERoleActionInfo[int(ERoleActionWeaponFire)],
	allERoleActionInfo[int32(ERoleActionWeaponActive)].toml: allERoleActionInfo[int(ERoleActionWeaponActive)],
	allERoleActionInfo[int32(ERoleActionWeaponReload)].toml: allERoleActionInfo[int(ERoleActionWeaponReload)],
	allERoleActionInfo[int32(ERoleActionThrow)].toml:        allERoleActionInfo[int(ERoleActionThrow)],
	allERoleActionInfo[int32(ERoleActionHeal)].toml:         allERoleActionInfo[int(ERoleActionHeal)],
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *ERoleAction) UnmarshalText(data []byte) error {
	key := string(data)
	v, ok := mapCodeToERoleActionInfo[key]
	if !ok {
		return fmt.Errorf("ERoleAction.UnmarshalText failed: invalid ERoleAction[%v]", key)
	}
	*e = v.value
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e ERoleAction) MarshalText() ([]byte, error) {
	return []byte(allERoleActionInfo[e].toml), nil
}

func (e ERoleAction) String() string {
	return allERoleActionInfo[e].toml
}

// EShootMode EShootMode
type EShootMode int32

const (
	// EShootModeAutoInterval 射击模式-连续射击-时间间隔
	EShootModeAutoInterval EShootMode = EShootMode(0)
	// EShootModeAutoAction 连续射击-拉栓动作
	EShootModeAutoAction EShootMode = EShootMode(1)
	// EShootModeAutoThree 三连射
	EShootModeAutoThree EShootMode = EShootMode(2)
)

type internalEShootModeInfo struct {
	value EShootMode
	toml  string
	desc  string
}

var allEShootModeInfo = []*internalEShootModeInfo{
	&internalEShootModeInfo{
		value: EShootModeAutoInterval,
		toml:  "EShootMode.AutoInterval",
		desc:  "射击模式-连续射击-时间间隔",
	},
	&internalEShootModeInfo{
		value: EShootModeAutoAction,
		toml:  "EShootMode.AutoAction",
		desc:  "连续射击-拉栓动作",
	},
	&internalEShootModeInfo{
		value: EShootModeAutoThree,
		toml:  "EShootMode.AutoThree",
		desc:  "三连射",
	},
}

var mapCodeToEShootModeInfo = map[string]*internalEShootModeInfo{
	allEShootModeInfo[int32(EShootModeAutoInterval)].toml: allEShootModeInfo[int(EShootModeAutoInterval)],
	allEShootModeInfo[int32(EShootModeAutoAction)].toml:   allEShootModeInfo[int(EShootModeAutoAction)],
	allEShootModeInfo[int32(EShootModeAutoThree)].toml:    allEShootModeInfo[int(EShootModeAutoThree)],
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
	// EThrowableTypeFragGrenade 投掷物-破片手雷
	EThrowableTypeFragGrenade EThrowableType = EThrowableType(0)
	// EThrowableTypeMolotovCocktail 燃烧弹
	EThrowableTypeMolotovCocktail EThrowableType = EThrowableType(1)
	// EThrowableTypeSmokeGrenade 烟雾弹
	EThrowableTypeSmokeGrenade EThrowableType = EThrowableType(2)
	// EThrowableTypeStunGrenade 闪光弹
	EThrowableTypeStunGrenade EThrowableType = EThrowableType(3)
)

type internalEThrowableTypeInfo struct {
	value EThrowableType
	toml  string
	desc  string
}

var allEThrowableTypeInfo = []*internalEThrowableTypeInfo{
	&internalEThrowableTypeInfo{
		value: EThrowableTypeFragGrenade,
		toml:  "EThrowableType.FragGrenade",
		desc:  "投掷物-破片手雷",
	},
	&internalEThrowableTypeInfo{
		value: EThrowableTypeMolotovCocktail,
		toml:  "EThrowableType.MolotovCocktail",
		desc:  "燃烧弹",
	},
	&internalEThrowableTypeInfo{
		value: EThrowableTypeSmokeGrenade,
		toml:  "EThrowableType.SmokeGrenade",
		desc:  "烟雾弹",
	},
	&internalEThrowableTypeInfo{
		value: EThrowableTypeStunGrenade,
		toml:  "EThrowableType.StunGrenade",
		desc:  "闪光弹",
	},
}

var mapCodeToEThrowableTypeInfo = map[string]*internalEThrowableTypeInfo{
	allEThrowableTypeInfo[int32(EThrowableTypeFragGrenade)].toml:     allEThrowableTypeInfo[int(EThrowableTypeFragGrenade)],
	allEThrowableTypeInfo[int32(EThrowableTypeMolotovCocktail)].toml: allEThrowableTypeInfo[int(EThrowableTypeMolotovCocktail)],
	allEThrowableTypeInfo[int32(EThrowableTypeSmokeGrenade)].toml:    allEThrowableTypeInfo[int(EThrowableTypeSmokeGrenade)],
	allEThrowableTypeInfo[int32(EThrowableTypeStunGrenade)].toml:     allEThrowableTypeInfo[int(EThrowableTypeStunGrenade)],
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
