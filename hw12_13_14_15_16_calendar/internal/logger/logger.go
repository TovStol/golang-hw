package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

type Logger struct {
	level  Level
	output io.WriteCloser
	logger *log.Logger
	mu     sync.Mutex
}

func New(level string, location string) *Logger {
	output := io.WriteCloser(os.Stdout)
	if location != "" {
		if err := os.MkdirAll(filepath.Dir(location), 0o755); err == nil {
			if file, err := os.OpenFile(location, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
				output = file
			}
		}
	}

	return &Logger{
		level:  parseLevel(level),
		output: output,
		logger: log.New(output, "", 0),
	}
}

func (l *Logger) Debug(msg string) {
	l.write(DebugLevel, "DEBUG", msg)
}

func (l *Logger) Info(msg string) {
	l.write(InfoLevel, "INFO", msg)
}

func (l *Logger) Warn(msg string) {
	l.write(WarnLevel, "WARN", msg)
}

func (l *Logger) Error(msg string) {
	l.write(ErrorLevel, "ERROR", msg)
}

func (l *Logger) Close() error {
	if l.output == os.Stdout {
		return nil
	}
	return l.output.Close()
}

func (l *Logger) write(level Level, levelName string, msg string) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.logger.Printf("[%s] %s", levelName, msg)
}

func parseLevel(level string) Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return DebugLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		return InfoLevel
	}
}
