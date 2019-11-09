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
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelQuite
)

const (
	trace   = "TRACE"
	debug   = "DEBUG"
	info    = "INFO"
	warn    = "WARNING"
	err     = "ERROR"
	quite   = "QUITE"
	PADDING = "\t"
)

var Level = LevelTrace
var pending string

//sets Level
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
	case quite:
		Level = LevelQuite
	default:
		log.Fatal("Given log level not understood!")
	}
}

//dumps http.Request
func LogRequest(r *http.Request) {
	x, err := httputil.DumpRequest(r, true)
	if err != nil {
		Warn("Can not dump given request! %s", err.Error())
		return
	}
	Info("A request received: %q", x)

}

//log in debug mode
func Debug(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params...)
	}
	if Level <= LevelDebug {
		log.Printf("%s%s%s\n", color.Bold(color.Magenta(debug)), PADDING, msg)
	}
}

//log in trace mode
func Trace(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params...)
	}
	if Level <= LevelTrace {
		log.Printf("%s%s%s\n", color.Bold(color.Faint(trace)), PADDING, msg)
	}
}

//prints headers of running processes
func Header(providers *Providers) {
	if Level == LevelQuite {
		return
	}
	printHeader(providers, 0)
}

//prints headers from given column index
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

//info level
func Info(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params...)
	}
	if Level <= LevelInfo {
		log.Printf("%s%s%s\n", color.Bold(color.Cyan(info)), PADDING, msg)
	}
}

//prints checksums or error values
//for each runners
func Logsum(providers *Providers) {
	if Level != LevelQuite {
		fmt.Println("\n\n")
	}
	for i, s := range providers.Statuses {
		if Level == LevelQuite {
			fmt.Println("%s\t%s", s.Checksum, providers.Filename)
			return
		}
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

	Info("Calculated file: %s", providers.Filename)
}

var currentColumn int

//prints status of each providers
//if log level INFO+ prints horizontally
//if horizontal 5 column per row
//if all processes finished in current row
//prints new one
func Status(providers *Providers) {
	if Level == LevelQuite {
		return
	}
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

//see Status comment
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

//prints given log at the same line by removing previous comment
//not used currently
func Inplace(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params...)
	}
	fmt.Printf("%s%s%s%s\r", PADDING, color.Bold(color.Cyan(info)), PADDING, msg)
}

//log in warn level
func Warn(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params...)
	}
	if Level <= LevelWarn {
		log.Printf("%s%s%s\n", color.Bold(color.Yellow(warn)), PADDING, msg)
	}
}

//log in error level
func Error(msg string, params ...interface{}) {
	if params != nil {
		msg = fmt.Sprintf(msg, params...)
	}
	if Level <= LevelError {
		log.Printf("%s%s%s\n", color.Bold(color.Red(err)), PADDING, msg)
	}
}
