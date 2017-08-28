package ffClientToml

import (
	"ffCommon/util"

	"ffAutoGen/ffEnum"

	"fmt"

	"github.com/lexical005/toml"
)

// Item excel Item
type Item struct {
	ItemTemplate map[int32]*ItemTemplate
	GunWeapon    map[int32]*GunWeapon
	Ammunition   map[int32]*Ammunition
	Attachment   map[int32]*Attachment
	MelleeWeapon map[int32]*MelleeWeapon
	Equipment    map[int32]*Equipment
	Consumable   map[int32]*Consumable
	Throwable    map[int32]*Throwable
}

func (i *Item) String() string {
	result := ""
	result += "ItemTemplate"
	for k, v := range i.ItemTemplate {
		result += fmt.Sprintf("%v:%v\n", k, v)
	}

	result += "GunWeapon"
	for k, v := range i.GunWeapon {
		result += fmt.Sprintf("%v:%v\n", k, v)
	}

	result += "Ammunition"
	for k, v := range i.Ammunition {
		result += fmt.Sprintf("%v:%v\n", k, v)
	}

	result += "Attachment"
	for k, v := range i.Attachment {
		result += fmt.Sprintf("%v:%v\n", k, v)
	}

	result += "MelleeWeapon"
	for k, v := range i.MelleeWeapon {
		result += fmt.Sprintf("%v:%v\n", k, v)
	}

	result += "Equipment"
	for k, v := range i.Equipment {
		result += fmt.Sprintf("%v:%v\n", k, v)
	}

	result += "Consumable"
	for k, v := range i.Consumable {
		result += fmt.Sprintf("%v:%v\n", k, v)
	}

	result += "Throwable"
	for k, v := range i.Throwable {
		result += fmt.Sprintf("%v:%v\n", k, v)
	}

	return result
}

// Name the toml config's name
func (i *Item) Name() string {
	return "Item"
}

// ItemTemplate sheet ItemTemplate of excel Item
type ItemTemplate struct {
	Name     string
	Desc     string
	AssetID  int32
	Icon     string
	ItemType ffEnum.EItemType
}

func (i *ItemTemplate) String() string {
	result := "["
	result += fmt.Sprintf("Name:%v,", i.Name)
	result += fmt.Sprintf("Desc:%v,", i.Desc)
	result += fmt.Sprintf("AssetID:%v,", i.AssetID)
	result += fmt.Sprintf("Icon:%v,", i.Icon)
	result += fmt.Sprintf("ItemType:%v,", i.ItemType)
	result += "]"
	return result
}

// GunWeapon sheet GunWeapon of excel Item
type GunWeapon struct {
	GunWeaponType   ffEnum.EGunWeaponType
	ShootMode       []ffEnum.EShootMode
	AmmunitionType  ffEnum.EAmmunitionType
	AttachmentTypes []ffEnum.EAttachmentType
	Attrs           map[ffEnum.EAttr]int32
	AttrsKey        []ffEnum.EAttr
	AttrsValue      []int32
}

func (g *GunWeapon) String() string {
	result := "["
	result += fmt.Sprintf("GunWeaponType:%v,", g.GunWeaponType)
	result += fmt.Sprintf("ShootMode:%v,", g.ShootMode)
	result += fmt.Sprintf("AmmunitionType:%v,", g.AmmunitionType)
	result += fmt.Sprintf("AttachmentTypes:%v,", g.AttachmentTypes)
	result += fmt.Sprintf("Attrs:%v,", g.Attrs)
	result += fmt.Sprintf("AttrsKey:%v,", g.AttrsKey)
	result += fmt.Sprintf("AttrsValue:%v,", g.AttrsValue)
	result += "]"
	return result
}

// Ammunition sheet Ammunition of excel Item
type Ammunition struct {
	AmmunitionType  ffEnum.EAmmunitionType
	AmmunitionStack int32
}

func (a *Ammunition) String() string {
	result := "["
	result += fmt.Sprintf("AmmunitionType:%v,", a.AmmunitionType)
	result += fmt.Sprintf("AmmunitionStack:%v,", a.AmmunitionStack)
	result += "]"
	return result
}

// Attachment sheet Attachment of excel Item
type Attachment struct {
	AttachmentType ffEnum.EAttachmentType
	GunWeapons     []int32
	ShutSound      int32
	ShutFire       int32
	Attrs          map[ffEnum.EAttr]int32
	AttrsKey       []ffEnum.EAttr
	AttrsValue     []int32
	Clip           map[int32]int32
	ClipKey        []int32
	ClipValue      []int32
}

func (a *Attachment) String() string {
	result := "["
	result += fmt.Sprintf("AttachmentType:%v,", a.AttachmentType)
	result += fmt.Sprintf("GunWeapons:%v,", a.GunWeapons)
	result += fmt.Sprintf("ShutSound:%v,", a.ShutSound)
	result += fmt.Sprintf("ShutFire:%v,", a.ShutFire)
	result += fmt.Sprintf("Attrs:%v,", a.Attrs)
	result += fmt.Sprintf("AttrsKey:%v,", a.AttrsKey)
	result += fmt.Sprintf("AttrsValue:%v,", a.AttrsValue)
	result += fmt.Sprintf("Clip:%v,", a.Clip)
	result += fmt.Sprintf("ClipKey:%v,", a.ClipKey)
	result += fmt.Sprintf("ClipValue:%v,", a.ClipValue)
	result += "]"
	return result
}

// MelleeWeapon sheet MelleeWeapon of excel Item
type MelleeWeapon struct {
	MelleeWeaponType ffEnum.EMelleeWeaponType
	Attrs            map[ffEnum.EAttr]int32
	AttrsKey         []ffEnum.EAttr
	AttrsValue       []int32
}

func (m *MelleeWeapon) String() string {
	result := "["
	result += fmt.Sprintf("MelleeWeaponType:%v,", m.MelleeWeaponType)
	result += fmt.Sprintf("Attrs:%v,", m.Attrs)
	result += fmt.Sprintf("AttrsKey:%v,", m.AttrsKey)
	result += fmt.Sprintf("AttrsValue:%v,", m.AttrsValue)
	result += "]"
	return result
}

// Equipment sheet Equipment of excel Item
type Equipment struct {
	EquipmentType ffEnum.EEquipmentType
	Attrs         map[ffEnum.EAttr]int32
	AttrsKey      []ffEnum.EAttr
	AttrsValue    []int32
}

func (e *Equipment) String() string {
	result := "["
	result += fmt.Sprintf("EquipmentType:%v,", e.EquipmentType)
	result += fmt.Sprintf("Attrs:%v,", e.Attrs)
	result += fmt.Sprintf("AttrsKey:%v,", e.AttrsKey)
	result += fmt.Sprintf("AttrsValue:%v,", e.AttrsValue)
	result += "]"
	return result
}

// Consumable sheet Consumable of excel Item
type Consumable struct {
	ConsumableType      ffEnum.EConsumableType
	UseTime             int32
	UseHpLimit          int32
	UseRecover          int32
	UseRecoverUpLimit   int32
	KeepTime            int32
	KeepRecoverInterval int32
	KeepRecover         int32
}

func (c *Consumable) String() string {
	result := "["
	result += fmt.Sprintf("ConsumableType:%v,", c.ConsumableType)
	result += fmt.Sprintf("UseTime:%v,", c.UseTime)
	result += fmt.Sprintf("UseHpLimit:%v,", c.UseHpLimit)
	result += fmt.Sprintf("UseRecover:%v,", c.UseRecover)
	result += fmt.Sprintf("UseRecoverUpLimit:%v,", c.UseRecoverUpLimit)
	result += fmt.Sprintf("KeepTime:%v,", c.KeepTime)
	result += fmt.Sprintf("KeepRecoverInterval:%v,", c.KeepRecoverInterval)
	result += fmt.Sprintf("KeepRecover:%v,", c.KeepRecover)
	result += "]"
	return result
}

// Throwable sheet Throwable of excel Item
type Throwable struct {
	ThrowableType ffEnum.EThrowableType
	RadiusClose   int32
}

func (t *Throwable) String() string {
	result := "["
	result += fmt.Sprintf("ThrowableType:%v,", t.ThrowableType)
	result += fmt.Sprintf("RadiusClose:%v,", t.RadiusClose)
	result += "]"
	return result
}

// ReadItem read excel Item
func ReadItem() (i *Item, err error) {
	// 读取文件内容
	fileContent, err := util.ReadFile("toml/client/Item.toml")
	if err != nil {
		return
	}

	// 解析
	i = &Item{}
	err = toml.Unmarshal(fileContent, i)
	if err != nil {
		return
	}

	for _, one := range i.GunWeapon {
		one.Attrs = make(map[ffEnum.EAttr]int32, len(one.AttrsKey))
		for index, v := range one.AttrsKey {
			one.Attrs[v] = one.AttrsValue[index]
		}
	}

	for _, one := range i.Attachment {
		one.Attrs = make(map[ffEnum.EAttr]int32, len(one.AttrsKey))
		for index, v := range one.AttrsKey {
			one.Attrs[v] = one.AttrsValue[index]
		}
	}

	for _, one := range i.Attachment {
		one.Clip = make(map[int32]int32, len(one.ClipKey))
		for index, v := range one.ClipKey {
			one.Clip[v] = one.ClipValue[index]
		}
	}

	for _, one := range i.MelleeWeapon {
		one.Attrs = make(map[ffEnum.EAttr]int32, len(one.AttrsKey))
		for index, v := range one.AttrsKey {
			one.Attrs[v] = one.AttrsValue[index]
		}
	}

	for _, one := range i.Equipment {
		one.Attrs = make(map[ffEnum.EAttr]int32, len(one.AttrsKey))
		for index, v := range one.AttrsKey {
			one.Attrs[v] = one.AttrsValue[index]
		}
	}

	return
}
