// Package main provides ...
package logger

import (
	"fmt"
	"log"

	. "github.com/getsumio/getsum/internal/providers"

	color "github.com/logrusorgru/aurora"
)

const (
	LevelTrace = 0
	LevelDebug = 1
	LevelInfo  = 2
	LevelWarn  = 3
	LevelError = 4
)

const (
	trace   = "TRACE"
	debug   = "DEBUG"
	info    = "INFO"
	warn    = "WARNING"
	err     = "ERROR"
	PADDING = "\t"
)

var Level = LevelTrace
var pending string

func Debug(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params)
	}
	if Level <= LevelDebug {
		log.Printf("%s%s%s\n", color.Bold(color.Magenta(debug)), PADDING, msg)
	}
}

func Trace(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params)
	}
	if Level <= LevelTrace {
		log.Printf("%s%s%s\n", color.Bold(color.Faint(trace)), PADDING, msg)
	}
}

func Header(providers []Provider) {
	var first string
	var second string
	for _, p := range providers {
		first = fmt.Sprintf("%s%s%s%s", first, PADDING, color.Bold(color.BgMagenta(color.Yellow(p.Data().Name))), PADDING)
		second = fmt.Sprintf("%s\t%s\t%s\t", second, "Status", "Value")
	}
	fmt.Printf("%s\n", first)
	fmt.Printf("%s\n", second)

}

func Info(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params)
	}
	if Level <= LevelInfo {
		log.Printf("%s%s%s\n", color.Bold(color.Cyan(info)), PADDING, msg)
	}
}

func Status(providers []Provider) {
	var msg string
	for _, p := range providers {
		msg = fmt.Sprintf("%s%s\t%s\t", msg, p.Data().Status, p.Data().Value)
	}
	fmt.Printf("%s\r", msg)
}

func Inplace(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params...)
	}
	fmt.Printf("%s%s%s%s\r", PADDING, color.Bold(color.Cyan(info)), PADDING, msg)
}

func Warn(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params)
	}
	if Level <= LevelWarn {
		log.Printf("%s%s%s\n", color.Bold(color.Yellow(warn)), PADDING, msg)
	}
}

func Error(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params)
	}
	if Level <= LevelError {
		log.Printf("%s%s%s\n", color.Bold(color.Red(err)), PADDING, msg)
	}
}
