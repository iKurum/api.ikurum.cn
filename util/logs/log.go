// Package logs 调用打印消息
package logs

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// logs 初始化 传入SetPrefix值
//
// 默认 [IKURUM]~
func Init(title ...string) {
	s := strings.Join(title, "")
	if s == "" {
		s = "[IKURUM]~"
	}
	log.SetPrefix(s)
	log.SetFlags(2)
	log.SetOutput(os.Stdout)
}

func Info(v ...interface{}) {
	log.Printf("\033[1;30;42m%v\033[0m %s\n", " INF: ", fmt.Sprint(v...))
}

func Warning(v ...interface{}) {
	log.Printf("\033[1;30;43m%v\033[0m %s\n", " WAR: ", fmt.Sprint(v...))
}

func Error(v ...interface{}) {
	log.Printf("\033[1;37;41m%v\033[0m %s\n", " ERR: ", fmt.Sprint(v...))
}

func Exit(v ...interface{}) {
	log.Printf("\033[1;37;45m%v\033[0m %s\n", " EXT: ", fmt.Sprint(v...))
	os.Exit(15)
}
