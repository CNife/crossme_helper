package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(initText)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			panic(scanner.Err())
		}
		line := scanner.Text()

		line = strings.ToLower(strings.TrimSpace(line))
		words := split(line)
		if len(words) > 0 {
			err := exec(words[0], words[1:])
			if err == quitWithoutError {
				return
			}
			if err != nil {
				_, err = fmt.Fprintln(os.Stderr, err)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func split(s string) []string {
	var result []string
	var buffer bytes.Buffer
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			if buffer.Len() > 0 {
				result = append(result, buffer.String())
				buffer.Reset()
			}
		} else {
			buffer.WriteByte(s[i])
		}
	}
	if buffer.Len() > 0 {
		result = append(result, buffer.String())
	}
	return result
}

const initText = `交互式 CrossMe 游戏助手
输入 help 获取帮助，输入 quit 退出
输入命令时可以只输入前缀`
