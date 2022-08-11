package logger

import (
	"fmt"
	"os"
	"time"
)

func getTime() string {
	return time.Now().Format("2006:01:02 15:04:05")
}

type logger struct {
	from string
}

func Log(from string) *logger {
	return &logger{from: from}
}

func (log *logger) Fatal(args ...interface{}) {
	fmt.Println(getTime(), log.from, args)
	os.Exit(1)
}

func (log *logger) Error(args ...interface{}) {
	fmt.Println(getTime(), log.from, args)
}

func (log *logger) Info(args ...interface{}) {
	fmt.Println(getTime(), log.from, args)
}
