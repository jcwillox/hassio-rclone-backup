package main

import (
	"github.com/jcwillox/emerald"
	"os"
	"time"
)

func Logln(tag string, color string, a ...interface{}) {
	emerald.Print(
		emerald.White, time.Now().Format("[2006-01-02] [15:04:05]"), emerald.Reset,
		" [", color, tag, emerald.Reset, "]: ", color,
	)
	emerald.Println(a...)
	emerald.Print(emerald.Reset)
}

func Debugln(a ...interface{}) {
	if config.LogLevel == "debug" {
		Logln("DEBUG", emerald.Cyan, a...)
	}
}

func Infoln(a ...interface{}) {
	Logln("INFO", emerald.Green, a...)
}

func Errorln(a ...interface{}) {
	Logln("ERROR", emerald.Red, a...)
}

func Fatalln(a ...interface{}) {
	Logln("FATAL", emerald.Bold+emerald.Red, a...)
	os.Exit(1)
}
