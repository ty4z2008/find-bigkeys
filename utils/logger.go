package utils

import (
	"fmt"
	"log"
	"os"
	// "time"
)

// A Logger represents an active logging object. Multiple loggers can be used
// simultaneously even if they are using the same same writers.
type Logger struct {
	infoLog     *log.Logger
	warningLog  *log.Logger
	errorLog    *log.Logger
	fatalLog    *log.Logger
	initialized bool
}

const (
	tagInfo    = "INFO"
	tagWarning = "WARN"
	tagError   = "ERROR"
	tagFatal   = "FATAL"
)

type severity int

const (
	flags          = log.Ldate | log.Lmicroseconds | log.Lshortfile
	sInfo severity = iota
	sWarning
	sError
	sFatal
)

func Info(v ...interface{}) {
	output(tagInfo, sInfo, fmt.Sprintln(v...))
}

func output(tag string, s severity, text string) {
	prefix := fmt.Sprintf("[%s] ", tag)
	switch s {
	case sInfo:
		logger := log.New(os.Stderr, prefix, flags)
		logger.Print(text)
	case sWarning:
	case sError:
	case sFatal:
	default:
		panic(fmt.Sprintln("unrecognized severity:", s))
	}

}
