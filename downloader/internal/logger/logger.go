// filepath: /home/jelle/wd/wow/golang-terminal-downloader/downloader/internal/logger/logger.go
package logger

import (
	"io"
	"log"
	"os"
)

var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func InitLogger(logLevel string) {
	Debug = log.New(io.Discard, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	switch logLevel {
	case "debug":
		Debug.SetOutput(os.Stdout)
		Info.SetOutput(os.Stdout)
		Warning.SetOutput(os.Stdout)
		Error.SetOutput(os.Stderr)
	case "info":
		Info.SetOutput(os.Stdout)
		Warning.SetOutput(os.Stdout)
		Error.SetOutput(os.Stderr)
	case "warning":
		Warning.SetOutput(os.Stdout)
		Error.SetOutput(os.Stderr)
	case "error":
		Error.SetOutput(os.Stderr)
	default:
		Info.SetOutput(os.Stdout)
		Warning.SetOutput(os.Stdout)
		Error.SetOutput(os.Stderr)
	}
}
