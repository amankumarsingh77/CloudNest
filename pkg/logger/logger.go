package logger

import (
	"log"
	"os"
)

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
	Debug *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		Info:  log.New(os.Stdout, "\033[32mINFO: \033[0m", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(os.Stdout, "\033[31mERROR: \033[0m", log.Ldate|log.Ltime|log.Lshortfile),
		Debug: log.New(os.Stdout, "\033[33mDEBUG: \033[0m", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) LogInfo(msg string) {
	l.Info.Println(msg)
}

func (l *Logger) LogError(msg string) {
	l.Error.Println(msg)
}

func (l *Logger) LogDebug(msg string) {
	l.Debug.Println(msg)
}
