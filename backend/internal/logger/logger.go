package logger

import (
	"log"
	"os"
)

const (
	flags = log.Ldate | log.Ltime | log.Lshortfile
)

var (
	Debug   = log.New(os.Stdout, "[DEBUG] ", flags)
	Info    = log.New(os.Stdout, "[INFO] ", flags)
	Warning = log.New(os.Stdout, "[WARNING] ", flags)
	Error   = log.New(os.Stderr, "[ERROR] ", flags)
	Fatal   = log.New(os.Stderr, "[FATAL]", flags)
)

func init() {
	stdout := os.Stdout
	stderr := os.Stderr

	Debug.SetOutput(stdout)
	Info.SetOutput(stdout)
	Error.SetOutput(stderr)
	Fatal.SetOutput(stderr)

	log.SetOutput(Debug.Writer())
	log.SetPrefix("[DEBUG]")
	log.SetFlags(flags)
}
