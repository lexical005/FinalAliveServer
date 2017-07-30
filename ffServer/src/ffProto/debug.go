package ffProto

import (
	"ffCommon/log/log"
)

// PrintModule debug print ffProto
func PrintModule() {
	log.RunLogger.Println("\nPrintModule Start ffProto:")

	log.RunLogger.Println("ffProto.messagePool:")
	for _, pool := range messagePool {
		log.RunLogger.Println(pool)
	}

	log.RunLogger.Println("ffProto.protoPool:")
	log.RunLogger.Println(protoPool)

	log.RunLogger.Println("ffProto.bufferPool:")
	for _, pool := range bufferPool {
		log.RunLogger.Println(pool)
	}

	log.RunLogger.Printf("PrintModule End ffProto\n\n")
}
