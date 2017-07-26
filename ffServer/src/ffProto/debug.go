package ffProto

import (
	"ffCommon/log/log"
)

// PrintModule debug print ffProto
func PrintModule() {
	log.RunLogger.Println("\nPrintModule Start ffProto:")

	log.RunLogger.Println("ffProto.protoPool:")
	for pool := range protoPool {
		log.RunLogger.Println(pool)
	}

	log.RunLogger.Println("ffProto.msgPool:")
	for pool := range msgPool {
		log.RunLogger.Println(pool)
	}

	log.RunLogger.Printf("PrintModule End ffProto\n\n")
}
