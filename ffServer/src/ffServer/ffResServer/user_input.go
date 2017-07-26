package main

import (
	"ffCommon/log/log"

	"bufio"
	"os"
	"strings"
)

func handleUserInput() {

deadloop:
	for {
		inputReader := bufio.NewReader(os.Stdin)
		input, err := inputReader.ReadString('\n')
		if err != nil {
			log.RunLogger.Println(err)
			continue
		}
		input = strings.ToUpper(input)

		b := []byte(input)
		input = string(b[:len(b)-2])

		if strings.HasSuffix(input, "\r\n") {
			input = input[:len(input)-2]
		}

		if strings.HasPrefix(input, "REMOTERES") {
			log.RunLogger.Println("user input REMOTERES ==> call genRemoteResMap")
			genRemoteResMap()
		} else if strings.HasPrefix(input, "HOTRES") {
			tmp := strings.Split(input, " ")
			if len(tmp) > 1 {
				var channelName = tmp[1]
				for name, channelInfo := range globalChannelInfo {
					if strings.ToUpper(name) == channelName {
						log.RunLogger.Printf("user input [%s] ==> call genHotResMap\n\n", input)
						genHotResMap(channelInfo, true)
						continue deadloop
					}
				}
				log.RunLogger.Printf("invalid channelName[%s]\n\n", channelName)
			}
			log.RunLogger.Println("valid hotres input format:\nhotres channelName")
		} else if input == "CONFIG" {

		} else if input == "LL" {
			log.RunLogger.Println(len(ll.sem), cap(ll.sem)-len(ll.sem))
		}
	}
}

func init() {
	go handleUserInput()
}
