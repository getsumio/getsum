// Package main provides ...
package logger

import (
	"fmt"
	"log"

	. "github.com/getsumio/getsum/internal/algorithm/supplier"
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
		first = fmt.Sprintf("%s%s%s%s", first, PADDING, color.Bold(color.Italic(color.BrightCyan(color.Underline(p.Data().Name)))), PADDING)
		second = fmt.Sprintf("%s%10s\t%5s |\t", second, "Status", "Value")
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
func Logsum(providers []Provider, stats []*Status) {
	fmt.Println("\n\n")
	for i, s := range stats {
		fmt.Printf("\t%s %s %s", color.Bold(color.BrightYellow("\u2713")), color.Bold(color.Blue(providers[i].Data().Name)), s.Checksum)
	}
}

func Status(stats []*Status) {
	var msg string
	for _, s := range stats {
		var v color.Value
		switch s.Status {
		case "TERMINATED", "TIMEDOUT", "ERROR":
			v = color.Red(s.Status)
		case "PREPARED":
			v = color.Cyan(s.Status)
		case "STARTED":
			v = color.Blink(s.Status)
		default:
			v = color.Green(s.Status)
		}
		msg = fmt.Sprintf("%s%10s\t%5s |\t", msg, color.Bold(v), color.Bold(color.Yellow(s.Value)))
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
