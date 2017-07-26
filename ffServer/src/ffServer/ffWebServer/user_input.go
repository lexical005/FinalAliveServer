package main

import (
	"ffCommon/log/log"

	"bufio"
	"fmt"
	"os"
	"strings"
)

func handleUserInput() {

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

		if input == "VIVO" {
			dictDatas := map[string]string{
				"channel":     "vivo",
				"orderAmount": "0.01",
				"orderDesc":   "orderDesc",
				"orderTitle":  "orderTitle",
				"storeOrder":  "1",
			}
			fmt.Println(vivo.onSetupIAPvivo("test", dictDatas))
		} else if strings.HasPrefix(input, "HOTRES") {
		} else if input == "CONFIG" {
		} else if input == "CLOSE" {
			mysql.close()
			os.Exit(0)
		}
	}
}

func init() {
	go handleUserInput()
}
