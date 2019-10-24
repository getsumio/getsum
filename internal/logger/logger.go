// Package main provides ...
package logger

import (
	"fmt"
	"log"

	. "github.com/getsumio/getsum/internal/file"
	. "github.com/getsumio/getsum/internal/provider/types"

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

func SetLevel(level string) {
	switch level {
	case trace:
		Level = LevelTrace
	case debug:
		Level = LevelDebug
	case info:
		Level = LevelInfo
	case warn:
		Level = LevelWarn
	case err:
		Level = LevelError
	default:
		log.Fatal("Given log level not understood!")
	}
}

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
		second = fmt.Sprintf("%s%10s\t%6s | ", second, "Status", "Value")
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
		var c color.Value
		var val string

		switch s.Status {
		case "COMPLETED":
			c = color.Bold(color.BrightYellow("\u2713"))
			val = s.Checksum
		case "MISMATCH":
			c = color.Bold(color.Red("\u24e7"))
			val = s.Checksum
		default:
			c = color.Bold(color.Red("\u24e7"))
			val = fmt.Sprintf("%s - %s", s.Status, s.Value)

		}

		fmt.Printf("\t%s %s %s\n", c, color.Bold(color.Blue(providers[i].Data().Name)), val)
	}
}

func Status(stats []*Status) {
	var msg string
	for _, s := range stats {
		var v color.Value
		var val string
		switch s.Status {
		case "TERMINATED", "TIMEDOUT":
			v = color.Red(s.Status)
			val = s.Value
		case "ERROR", "MISMATCH":
			v = color.Red(s.Status)
		case "PREPARED":
			v = color.Cyan(s.Status)
		case "STARTED":
			v = color.Blink(s.Status)
		default:
			v = color.Green(s.Status)
			val = s.Value
		}
		msg = fmt.Sprintf("%s%10s\t%6s | ", msg, color.Bold(v), color.Bold(color.Yellow(val)))
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
