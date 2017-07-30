package ffProto

import (
	"ffCommon/log/log"
)

// PrintModule debug print ffProto
func PrintModule() {
	log.RunLogger.Println("\nPrintModule Start ffProto:")

	log.RunLogger.Println("ffProto.messagePool:")
	for pool := range messagePool {
		log.RunLogger.Println(pool)
	}

	log.RunLogger.Println("ffProto.protoPool:")
	log.RunLogger.Println(protoPool)

	log.RunLogger.Println("ffProto.bufferPool:")
	log.RunLogger.Println(bufferPool)

	log.RunLogger.Printf("PrintModule End ffProto\n\n")
}
