package logging

import (
	"log"
	"os"
)

type Logger struct {
	info     *log.Logger
	warning  *log.Logger
	err      *log.Logger
	critical *log.Logger
}

func New() *Logger {
	return &Logger{
		info:     log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warning:  log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		err:      log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		critical: log.New(os.Stderr, "CRITICAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) GetErrorLogger() *log.Logger {
	return l.err
}

func (l *Logger) Info(msg string) {
	l.info.Printf(msg)
}

func (l *Logger) Warning(msg string) {
	l.warning.Printf(msg)
}

func (l *Logger) Error(msg string) {
	l.err.Printf(msg)
}

func (l *Logger) Critical(msg string) {
	l.critical.Printf(msg)
}
