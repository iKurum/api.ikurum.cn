package logs

import (
	"fmt"
	"log"
	"os"
)

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
