package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Level int

const (
	INFO Level = iota
	WARN
	ERROR
)

var levelNames = []string{
	"INFO",
	"WARN",
	"ERROR",
}

// Logger is a simple thread-safe logger
type Logger struct {
	mu     sync.Mutex
	out    io.Writer
	levels map[Level]bool
}

// New creates a logger writing to stdout by default
func New() *Logger {
	return &Logger{
		out: os.Stdout,
		levels: map[Level]bool{
			INFO:  true,
			WARN:  true,
			ERROR: true,
		},
	}
}

// SetOutput allows writing logs to a file or other io.Writer
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

// EnableLevel allows enabling/disabling specific log levels
func (l *Logger) EnableLevel(level Level, enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.levels[level] = enabled
}

// log prints a message with timestamp and level
func (l *Logger) log(level Level, format string, a ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.levels[level] {
		return
	}

	ts := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(l.out, "[%s] [%s] %s\n", ts, levelNames[level], msg)
}

// Info logs an info-level message
func (l *Logger) Info(format string, a ...any) {
	l.log(INFO, format, a...)
}

// Warn logs a warning
func (l *Logger) Warn(format string, a ...any) {
	l.log(WARN, format, a...)
}

// Error logs an error
func (l *Logger) Error(format string, a ...any) {
	l.log(ERROR, format, a...)
}
