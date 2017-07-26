package main

import (
	"ffCommon/log/log"
	"ffLogic/ffDef"
	"testing"
	"time"
)

var uuidGen *gameWorldUUIDGen

func Test_uuidGen1(t *testing.T) {
	uuidGen = &gameWorldUUIDGen{}
	uuidGen.init()

	for i := 0; i < 10; i++ {
		select {
		case <-time.After(time.Millisecond * 160):
			log.RunLogger.Printf("%X\n", uuidGen.Gen(ffDef.UUIDTimer))
		case <-time.After(time.Millisecond * 160):
			log.RunLogger.Printf("%X\n", uuidGen.Gen(ffDef.UUIDEquipment))
		}
	}
}
