package log

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	debugLvl   = "debug"
	warningLvl = "warn "
	errorLvl   = "error"
	infoLvl    = "info "
)

func Debug(format string, args ...interface{}) {
	fmt.Println(wrapper(format, debugLvl, args...))
}

func Warn(format string, args ...interface{}) {
	fmt.Println(wrapper(format, warningLvl, args...))
}

func Error(format string, args ...interface{}) {
	fmt.Println(wrapper(format, errorLvl, args...))
}

func Info(format string, args ...interface{}) {
	fmt.Println(wrapper(format, infoLvl, args...))
}
func Fatal(format string, args ...interface{}) {
	fmt.Println(wrapper(format, infoLvl, args...))
	log.Fatal()
	os.Exit(0)
}

func wrapper(txt string, lvl string, args ...interface{}) string {
	if len(args) > 0 {
		txt = fmt.Sprintf(txt, args...)
	}
	return fmt.Sprintf("[%s %s] %s", lvl, timeForLog(), txt)
}

func timeForLog() string  {
	return time.Now().Format("2006.01.02 15:04:05")
}
