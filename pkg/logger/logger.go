package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	*log.Logger
	level LogLevel
}

func NewLogger(level LogLevel) *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", 0),
		level:  level,
	}
}

func (l *Logger) log(level LogLevel, message string, args ...interface{}) {
	if level < l.level {
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	// Extract just the filename from the full path
	if idx := strings.LastIndex(file, "/"); idx >= 0 {
		file = file[idx+1:]
	}

	// Format the log message
	levelStr := [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[level]
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	caller := fmt.Sprintf("%s:%d", file, line)
	logMessage := fmt.Sprintf(message, args...)

	// Log the message
	l.Logger.Printf("[%s] %s %-5s %s: %s", timestamp, caller, levelStr, message, logMessage)

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.log(DEBUG, message, args...)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.log(INFO, message, args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(WARN, message, args...)
}

func (l *Logger) Error(message string, args ...interface{}) {
	l.log(ERROR, message, args...)
}

func (l *Logger) Fatal(message string, args ...interface{}) {
	l.log(FATAL, message, args...)
}

// SetLevel allows changing the logging level at runtime
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}