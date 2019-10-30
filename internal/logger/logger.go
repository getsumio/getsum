// Package main provides ...
package logger

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

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

func LogRequest(r *http.Request) {
	x, err := httputil.DumpRequest(r, true)
	if err != nil {
		Warn("Can not dump given request! %s", err.Error())
		return
	}
	Info("A request received: %q", x)

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

func Header(providers *Providers) {
	printHeader(providers, 0)
}
func printHeader(providers *Providers, start int) {
	var first string
	var second string
	var count int
	for i := start; i < providers.Length; i++ {
		if count >= 5 {
			break
		}
		count++
		first = fmt.Sprintf("%s%s%-23s", first, " ", color.Bold(color.Italic(color.BrightCyan(color.Underline((*providers.All[i]).Data().Name)))))
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
func Logsum(providers *Providers) {
	fmt.Println("\n\n")
	for i, s := range providers.Statuses {
		var c color.Value
		var val string

		switch s.Type {
		case status.COMPLETED, status.VALIDATED:
			c = color.Bold(color.BrightYellow("\u2713"))
			val = s.Checksum
		case status.MISMATCH:
			c = color.Bold(color.Red("\u24e7"))
			val = s.Checksum
		default:
			c = color.Bold(color.Red("\u24e7"))
			val = fmt.Sprintf("%s - %s", s.Type.Name(), s.Value)

		}

		fmt.Printf("\t%s %s %s\n", c, color.Bold(color.Blue((*providers.All[i]).Data().Name)), val)
	}
}

var currentColumn int

func Status(providers *Providers) {
	stats := providers.Status()
	printStatus(stats, providers, currentColumn)
	var anyRunner bool
	for j := currentColumn; j < currentColumn+5; j++ {
		if j >= providers.Length {
			break
		}
		if stats[j].Type < status.COMPLETED || (providers.HasValidation && stats[j].Type < status.SUSPENDED) || stats[j].Type == status.SUSPENDED {
			anyRunner = true
			break
		}
	}
	if !anyRunner {
		currentColumn += 5
		if currentColumn >= providers.Length {
			return
		}
		fmt.Println("\n")
		printHeader(providers, currentColumn)
		Status(providers)
	}
}

func printStatus(stats []*status.Status, providers *Providers, start int) {
	var msg string
	for i := start; i < providers.Length; i++ {
		if i >= start+5 {
			diff := providers.Length - i
			msg = fmt.Sprintf("%s ...%d more running", msg, diff)
			break
		}
		var v color.Value
		var val string
		switch stats[i].Type {
		case status.TERMINATED, status.TIMEDOUT, status.MISMATCH:
			v = color.Red(stats[i].Type.Name())
			val = stats[i].Value
		case status.ERROR:
			v = color.Red(stats[i].Type.Name())
		case status.PREPARED:
			v = color.Cyan(stats[i].Type.Name())
		case status.STARTED:
			v = color.Blink(stats[i].Type.Name())
		default:
			v = color.Green(stats[i].Type.Name())
			val = stats[i].Value
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
