package main

import (
	"bufio"
	"ffCommon/log/log"
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
		input = strings.ToLower(input)

		b := []byte(input)
		input = string(b[:len(b)-2])

		if strings.HasSuffix(input, "\r\n") {
			input = input[:len(input)-2]
		}

		log.RunLogger.Printf("handleUserInput: %v", input)

		if input == "close" {
			// 标识退出状态
			applicationQuit = true
			// 通知goroutine, 进程要退出
			close(chApplicationQuit)
			break deadloop
		}
	}
}

func init() {
	go handleUserInput()
}
