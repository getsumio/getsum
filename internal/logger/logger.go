// Package main provides ...
package logger

import (
	"fmt"
	"log"

	. "github.com/getsumio/getsum/internal/provider/types"
	"github.com/getsumio/getsum/internal/status"

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
		msg = fmt.Sprintf(msg, params...)
	}
	if Level <= LevelDebug {
		log.Printf("%s%s%s\n", color.Bold(color.Magenta(debug)), PADDING, msg)
	}
}

func Trace(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params...)
	}
	if Level <= LevelTrace {
		log.Printf("%s%s%s\n", color.Bold(color.Faint(trace)), PADDING, msg)
	}
}

func Header(providers []Provider) {
	var first string
	var second string
	for i, p := range providers {
		if i > 6 {
			break
		}
		first = fmt.Sprintf("%s%s%s%s", first, PADDING, color.Bold(color.Italic(color.BrightCyan(color.Underline(p.Data().Name)))), PADDING)
		second = fmt.Sprintf("%s%10s\t%6s | ", second, "Status", "Value")
	}
	fmt.Printf("%s\n", first)
	fmt.Printf("%s\n", second)

}

func Info(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params...)
	}
	if Level <= LevelInfo {
		log.Printf("%s%s%s\n", color.Bold(color.Cyan(info)), PADDING, msg)
	}
}
func Logsum(providers []Provider, stats []*status.Status) {
	fmt.Println("\n\n")
	for i, s := range stats {
		var c color.Value
		var val string

		switch s.Type {
		case status.COMPLETED:
			c = color.Bold(color.BrightYellow("\u2713"))
			val = s.Checksum
		case status.MISMATCH:
			c = color.Bold(color.Red("\u24e7"))
			val = s.Checksum
		default:
			c = color.Bold(color.Red("\u24e7"))
			val = fmt.Sprintf("%s - %s", s.Type.Name(), s.Value)

		}

		fmt.Printf("\t%s %s %s\n", c, color.Bold(color.Blue(providers[i].Data().Name)), val)
	}
}

func Status(stats []*status.Status) {
	var msg string
	for i, s := range stats {
		if i > 6 {
			diff := len(stats) - 5
			msg = fmt.Sprintf("%s ...%d more running", msg, diff)
			break
		}
		var v color.Value
		var val string
		switch s.Type {
		case status.TERMINATED, status.TIMEDOUT, status.MISMATCH:
			v = color.Red(s.Type.Name())
			val = s.Value
		case status.ERROR:
			v = color.Red(s.Type.Name())
		case status.PREPARED:
			v = color.Cyan(s.Type.Name())
		case status.STARTED:
			v = color.Blink(s.Type.Name())
		default:
			v = color.Green(s.Type.Name())
			val = s.Value
		}
		msg = fmt.Sprintf("%s%10s\t%6s | ", msg, color.Bold(v), color.Bold(color.Yellow(val)))
	}
	if Level < LevelInfo {
		fmt.Printf("%s\n", msg)
	} else {
		fmt.Printf("%s\r", msg)
	}
}

func Inplace(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params...)
	}
	fmt.Printf("%s%s%s%s\r", PADDING, color.Bold(color.Cyan(info)), PADDING, msg)
}

func Warn(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params...)
	}
	if Level <= LevelWarn {
		log.Printf("%s%s%s\n", color.Bold(color.Yellow(warn)), PADDING, msg)
	}
}

func Error(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params...)
	}
	if Level <= LevelError {
		log.Printf("%s%s%s\n", color.Bold(color.Red(err)), PADDING, msg)
	}
}
