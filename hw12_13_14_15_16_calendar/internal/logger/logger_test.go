package logger

import (
	"bytes"
	"io"
	"log"
	"strings"
	"testing"
)

type nopWriteCloser struct{ io.Writer }

func (nopWriteCloser) Close() error { return nil }

func newTestLogger(level string) (*Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	l := &Logger{
		level:  parseLevel(level),
		output: nopWriteCloser{buf},
		logger: log.New(buf, "", 0),
	}
	return l, buf
}

func TestLogger_InfoLevel(t *testing.T) {
	l, buf := newTestLogger("info")

	l.Debug("hidden")
	l.Info("visible")
	l.Warn("also visible")
	l.Error("also visible")

	if strings.Contains(buf.String(), "hidden") {
		t.Error("Debug message should be filtered at info level")
	}
	if !strings.Contains(buf.String(), "visible") {
		t.Error("Info message should be written at info level")
	}
}

func TestLogger_DebugLevel(t *testing.T) {
	l, buf := newTestLogger("debug")

	l.Debug("debug msg")

	if !strings.Contains(buf.String(), "debug msg") {
		t.Error("Debug message should be written at debug level")
	}
}

func TestLogger_ErrorLevel(t *testing.T) {
	l, buf := newTestLogger("error")

	l.Info("filtered")
	l.Warn("filtered")
	l.Error("shown")

	if strings.Contains(buf.String(), "filtered") {
		t.Error("Info/Warn should be filtered at error level")
	}
	if !strings.Contains(buf.String(), "shown") {
		t.Error("Error message should be written at error level")
	}
}

func TestLogger_LevelPrefix(t *testing.T) {
	tests := []struct {
		level  string
		method func(*Logger, string)
		want   string
	}{
		{"debug", (*Logger).Debug, "[DEBUG]"},
		{"debug", (*Logger).Info, "[INFO]"},
		{"debug", (*Logger).Warn, "[WARN]"},
		{"debug", (*Logger).Error, "[ERROR]"},
	}
	for _, tc := range tests {
		l, buf := newTestLogger(tc.level)
		tc.method(l, "msg")
		if !strings.Contains(buf.String(), tc.want) {
			t.Errorf("expected prefix %q in output %q", tc.want, buf.String())
		}
	}
}

func TestLogger_ParseLevel(t *testing.T) {
	cases := []struct {
		input string
		want  Level
	}{
		{"debug", DebugLevel},
		{"DEBUG", DebugLevel},
		{"info", InfoLevel},
		{"warn", WarnLevel},
		{"warning", WarnLevel},
		{"error", ErrorLevel},
		{"unknown", InfoLevel},
		{"", InfoLevel},
	}
	for _, tc := range cases {
		got := parseLevel(tc.input)
		if got != tc.want {
			t.Errorf("parseLevel(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestLogger_CloseStdout(t *testing.T) {
	l := New("info", "")
	if err := l.Close(); err != nil {
		t.Errorf("Close() on stdout logger returned error: %v", err)
	}
}
