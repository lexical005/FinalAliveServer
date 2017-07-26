package main

import (
	"ffCommon/log/log"
	"ffCommon/uuid"
	"ffLogic/ffDef"
)

type gameWorldUUIDGen struct {
	gens []*uuid.Generator
}

// Gen 生成 UUID
func (gwug *gameWorldUUIDGen) Gen(uuidType ffDef.UUIDType) uuid.UUID {
	return gwug.gens[uuidType].Gen()
}

func (gwug *gameWorldUUIDGen) init() {
	count := int(ffDef.UUIDTypeCount)
	gwug.gens = make([]*uuid.Generator, count, count)
	for i := 0; i < count; i++ {
		requester := uint64(appConfig.Server.ServerID)
		gen, err := uuid.NewGenerator(requester)
		if err != nil {
			log.FatalLogger.Printf("gameWorldUUIDGen.init: uuid.NewGenerator get error[%v]", err)
		}
		gwug.gens[i] = gen
	}
}
