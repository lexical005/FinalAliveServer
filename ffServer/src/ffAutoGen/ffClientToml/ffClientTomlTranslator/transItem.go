package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"path/filepath"
	"sort"

	proto "github.com/golang/protobuf/proto"
)

func transItem() {
	message := &Item{}

	// ItemTemplate
	ItemTemplateKeys := make([]int, 0, len(tomlItem.ItemTemplate)) // 必须使用64位机器
	//ItemTemplateKeys := make([]int, 0, len(tomlItem.ItemTemplate)) // 必须使用64位机器
	//ItemTemplateKeys := make([]string, 0, len(tomlItem.ItemTemplate)) // 必须使用64位机器
	for key := range tomlItem.ItemTemplate {
		ItemTemplateKeys = append(ItemTemplateKeys, int(key))
		//ItemTemplateKeys = append(ItemTemplateKeys, int(key))
		//ItemTemplateKeys = append(ItemTemplateKeys, string(key))
	}
	sort.Ints(ItemTemplateKeys)
	//sort.Ints(ItemTemplateKeys)
	//sort.Strings(ItemTemplateKeys)

	message.ItemTemplateKey = make([]int32, len(tomlItem.ItemTemplate))
	message.ItemTemplateValue = make([]*Item_StItemTemplate, len(tomlItem.ItemTemplate))
	for k, key := range ItemTemplateKeys {
		i := int32(key)
		//i := int32(key)
		//i := int32(key)
		v := tomlItem.ItemTemplate[i]

		message.ItemTemplateKey[k] = i
		message.ItemTemplateValue[k] = &Item_StItemTemplate{
			Name:     v.Name,
			Desc:     v.Desc,
			SceneKey: v.SceneKey,
			Icon:     v.Icon,
		}

		message.ItemTemplateValue[k].ItemType = int32(v.ItemType)
	}

	// GunWeapon
	GunWeaponKeys := make([]int, 0, len(tomlItem.GunWeapon)) // 必须使用64位机器
	//GunWeaponKeys := make([]int, 0, len(tomlItem.GunWeapon)) // 必须使用64位机器
	//GunWeaponKeys := make([]string, 0, len(tomlItem.GunWeapon)) // 必须使用64位机器
	for key := range tomlItem.GunWeapon {
		GunWeaponKeys = append(GunWeaponKeys, int(key))
		//GunWeaponKeys = append(GunWeaponKeys, int(key))
		//GunWeaponKeys = append(GunWeaponKeys, string(key))
	}
	sort.Ints(GunWeaponKeys)
	//sort.Ints(GunWeaponKeys)
	//sort.Strings(GunWeaponKeys)

	message.GunWeaponKey = make([]int32, len(tomlItem.GunWeapon))
	message.GunWeaponValue = make([]*Item_StGunWeapon, len(tomlItem.GunWeapon))
	for k, key := range GunWeaponKeys {
		i := int32(key)
		//i := int32(key)
		//i := int32(key)
		v := tomlItem.GunWeapon[i]

		message.GunWeaponKey[k] = i
		message.GunWeaponValue[k] = &Item_StGunWeapon{
			AttrsValue: v.AttrsValue,
		}

		message.GunWeaponValue[k].GunWeaponType = int32(v.GunWeaponType)
		message.GunWeaponValue[k].ShootMode = make([]int32, len(v.ShootMode), len(v.ShootMode))
		for xx, yy := range v.ShootMode {
			message.GunWeaponValue[k].ShootMode[xx] = int32(yy)
		}
		message.GunWeaponValue[k].AmmunitionType = int32(v.AmmunitionType)
		message.GunWeaponValue[k].AttachmentTypes = make([]int32, len(v.AttachmentTypes), len(v.AttachmentTypes))
		for xx, yy := range v.AttachmentTypes {
			message.GunWeaponValue[k].AttachmentTypes[xx] = int32(yy)
		}
		message.GunWeaponValue[k].AttrsKey = make([]int32, len(v.AttrsKey), len(v.AttrsKey))
		for xx, yy := range v.AttrsKey {
			message.GunWeaponValue[k].AttrsKey[xx] = int32(yy)
		}
	}

	// Ammunition
	AmmunitionKeys := make([]int, 0, len(tomlItem.Ammunition)) // 必须使用64位机器
	//AmmunitionKeys := make([]int, 0, len(tomlItem.Ammunition)) // 必须使用64位机器
	//AmmunitionKeys := make([]string, 0, len(tomlItem.Ammunition)) // 必须使用64位机器
	for key := range tomlItem.Ammunition {
		AmmunitionKeys = append(AmmunitionKeys, int(key))
		//AmmunitionKeys = append(AmmunitionKeys, int(key))
		//AmmunitionKeys = append(AmmunitionKeys, string(key))
	}
	sort.Ints(AmmunitionKeys)
	//sort.Ints(AmmunitionKeys)
	//sort.Strings(AmmunitionKeys)

	message.AmmunitionKey = make([]int32, len(tomlItem.Ammunition))
	message.AmmunitionValue = make([]*Item_StAmmunition, len(tomlItem.Ammunition))
	for k, key := range AmmunitionKeys {
		i := int32(key)
		//i := int32(key)
		//i := int32(key)
		v := tomlItem.Ammunition[i]

		message.AmmunitionKey[k] = i
		message.AmmunitionValue[k] = &Item_StAmmunition{
			AmmunitionStack: v.AmmunitionStack,
		}

		message.AmmunitionValue[k].AmmunitionType = int32(v.AmmunitionType)
	}

	// Attachment
	AttachmentKeys := make([]int, 0, len(tomlItem.Attachment)) // 必须使用64位机器
	//AttachmentKeys := make([]int, 0, len(tomlItem.Attachment)) // 必须使用64位机器
	//AttachmentKeys := make([]string, 0, len(tomlItem.Attachment)) // 必须使用64位机器
	for key := range tomlItem.Attachment {
		AttachmentKeys = append(AttachmentKeys, int(key))
		//AttachmentKeys = append(AttachmentKeys, int(key))
		//AttachmentKeys = append(AttachmentKeys, string(key))
	}
	sort.Ints(AttachmentKeys)
	//sort.Ints(AttachmentKeys)
	//sort.Strings(AttachmentKeys)

	message.AttachmentKey = make([]int32, len(tomlItem.Attachment))
	message.AttachmentValue = make([]*Item_StAttachment, len(tomlItem.Attachment))
	for k, key := range AttachmentKeys {
		i := int32(key)
		//i := int32(key)
		//i := int32(key)
		v := tomlItem.Attachment[i]

		message.AttachmentKey[k] = i
		message.AttachmentValue[k] = &Item_StAttachment{
			GunWeapons: v.GunWeapons,
			ShutSound:  v.ShutSound,
			ShutFire:   v.ShutFire,
			AttrsValue: v.AttrsValue,
			ClipKey:    v.ClipKey,
			ClipValue:  v.ClipValue,
		}

		message.AttachmentValue[k].AttachmentType = int32(v.AttachmentType)
		message.AttachmentValue[k].AttrsKey = make([]int32, len(v.AttrsKey), len(v.AttrsKey))
		for xx, yy := range v.AttrsKey {
			message.AttachmentValue[k].AttrsKey[xx] = int32(yy)
		}
	}

	// MelleeWeapon
	MelleeWeaponKeys := make([]int, 0, len(tomlItem.MelleeWeapon)) // 必须使用64位机器
	//MelleeWeaponKeys := make([]int, 0, len(tomlItem.MelleeWeapon)) // 必须使用64位机器
	//MelleeWeaponKeys := make([]string, 0, len(tomlItem.MelleeWeapon)) // 必须使用64位机器
	for key := range tomlItem.MelleeWeapon {
		MelleeWeaponKeys = append(MelleeWeaponKeys, int(key))
		//MelleeWeaponKeys = append(MelleeWeaponKeys, int(key))
		//MelleeWeaponKeys = append(MelleeWeaponKeys, string(key))
	}
	sort.Ints(MelleeWeaponKeys)
	//sort.Ints(MelleeWeaponKeys)
	//sort.Strings(MelleeWeaponKeys)

	message.MelleeWeaponKey = make([]int32, len(tomlItem.MelleeWeapon))
	message.MelleeWeaponValue = make([]*Item_StMelleeWeapon, len(tomlItem.MelleeWeapon))
	for k, key := range MelleeWeaponKeys {
		i := int32(key)
		//i := int32(key)
		//i := int32(key)
		v := tomlItem.MelleeWeapon[i]

		message.MelleeWeaponKey[k] = i
		message.MelleeWeaponValue[k] = &Item_StMelleeWeapon{
			AttrsValue: v.AttrsValue,
		}

		message.MelleeWeaponValue[k].MelleeWeaponType = int32(v.MelleeWeaponType)
		message.MelleeWeaponValue[k].AttrsKey = make([]int32, len(v.AttrsKey), len(v.AttrsKey))
		for xx, yy := range v.AttrsKey {
			message.MelleeWeaponValue[k].AttrsKey[xx] = int32(yy)
		}
	}

	// Equipment
	EquipmentKeys := make([]int, 0, len(tomlItem.Equipment)) // 必须使用64位机器
	//EquipmentKeys := make([]int, 0, len(tomlItem.Equipment)) // 必须使用64位机器
	//EquipmentKeys := make([]string, 0, len(tomlItem.Equipment)) // 必须使用64位机器
	for key := range tomlItem.Equipment {
		EquipmentKeys = append(EquipmentKeys, int(key))
		//EquipmentKeys = append(EquipmentKeys, int(key))
		//EquipmentKeys = append(EquipmentKeys, string(key))
	}
	sort.Ints(EquipmentKeys)
	//sort.Ints(EquipmentKeys)
	//sort.Strings(EquipmentKeys)

	message.EquipmentKey = make([]int32, len(tomlItem.Equipment))
	message.EquipmentValue = make([]*Item_StEquipment, len(tomlItem.Equipment))
	for k, key := range EquipmentKeys {
		i := int32(key)
		//i := int32(key)
		//i := int32(key)
		v := tomlItem.Equipment[i]

		message.EquipmentKey[k] = i
		message.EquipmentValue[k] = &Item_StEquipment{
			AttrsValue: v.AttrsValue,
		}

		message.EquipmentValue[k].EquipmentType = int32(v.EquipmentType)
		message.EquipmentValue[k].AttrsKey = make([]int32, len(v.AttrsKey), len(v.AttrsKey))
		for xx, yy := range v.AttrsKey {
			message.EquipmentValue[k].AttrsKey[xx] = int32(yy)
		}
	}

	// Consumable
	ConsumableKeys := make([]int, 0, len(tomlItem.Consumable)) // 必须使用64位机器
	//ConsumableKeys := make([]int, 0, len(tomlItem.Consumable)) // 必须使用64位机器
	//ConsumableKeys := make([]string, 0, len(tomlItem.Consumable)) // 必须使用64位机器
	for key := range tomlItem.Consumable {
		ConsumableKeys = append(ConsumableKeys, int(key))
		//ConsumableKeys = append(ConsumableKeys, int(key))
		//ConsumableKeys = append(ConsumableKeys, string(key))
	}
	sort.Ints(ConsumableKeys)
	//sort.Ints(ConsumableKeys)
	//sort.Strings(ConsumableKeys)

	message.ConsumableKey = make([]int32, len(tomlItem.Consumable))
	message.ConsumableValue = make([]*Item_StConsumable, len(tomlItem.Consumable))
	for k, key := range ConsumableKeys {
		i := int32(key)
		//i := int32(key)
		//i := int32(key)
		v := tomlItem.Consumable[i]

		message.ConsumableKey[k] = i
		message.ConsumableValue[k] = &Item_StConsumable{
			UseTime:             v.UseTime,
			UseHpLimit:          v.UseHpLimit,
			UseRecover:          v.UseRecover,
			UseRecoverUpLimit:   v.UseRecoverUpLimit,
			KeepTime:            v.KeepTime,
			KeepRecoverInterval: v.KeepRecoverInterval,
			KeepRecover:         v.KeepRecover,
		}

		message.ConsumableValue[k].ConsumableType = int32(v.ConsumableType)
	}

	// Throwable
	ThrowableKeys := make([]int, 0, len(tomlItem.Throwable)) // 必须使用64位机器
	//ThrowableKeys := make([]int, 0, len(tomlItem.Throwable)) // 必须使用64位机器
	//ThrowableKeys := make([]string, 0, len(tomlItem.Throwable)) // 必须使用64位机器
	for key := range tomlItem.Throwable {
		ThrowableKeys = append(ThrowableKeys, int(key))
		//ThrowableKeys = append(ThrowableKeys, int(key))
		//ThrowableKeys = append(ThrowableKeys, string(key))
	}
	sort.Ints(ThrowableKeys)
	//sort.Ints(ThrowableKeys)
	//sort.Strings(ThrowableKeys)

	message.ThrowableKey = make([]int32, len(tomlItem.Throwable))
	message.ThrowableValue = make([]*Item_StThrowable, len(tomlItem.Throwable))
	for k, key := range ThrowableKeys {
		i := int32(key)
		//i := int32(key)
		//i := int32(key)
		v := tomlItem.Throwable[i]

		message.ThrowableKey[k] = i
		message.ThrowableValue[k] = &Item_StThrowable{
			RadiusClose: v.RadiusClose,
		}

		message.ThrowableValue[k].ThrowableType = int32(v.ThrowableType)
	}

	pbBuf := proto.NewBuffer(make([]byte, 0, 1024*10))
	if err := pbBuf.Marshal(message); err != nil {
		log.RunLogger.Printf("transItem err[%v]", err)
		return
	}

	util.WriteFile(filepath.Join("ProtoBuf", "Client", "bytes", tomlItem.Name()+".bytes"), pbBuf.Bytes())
}

func init() {
	allTrans = append(allTrans, transItem)
}
