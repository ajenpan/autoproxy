package log

import (
	"fmt"
)

type Logger interface {
	Fields(fields map[string]interface{}) Logger
	Log(level Level, v ...interface{})
	Logf(level Level, format string, v ...interface{})

	// Info(args ...interface{})
	// Infof(template string, args ...interface{})
	// Trace(args ...interface{})
	// Tracef(template string, args ...interface{})
	// Debug(args ...interface{})
	// Debugf(template string, args ...interface{})
	// Warn(args ...interface{})
	// Warnf(template string, args ...interface{})
	// Error(args ...interface{})
	// Errorf(template string, args ...interface{})
	// Fatal(args ...interface{})
	// Fatalf(template string, args ...interface{})
}

type Level int8

const (
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel Level = iota - 2
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// InfoLevel is the default logging priority.
	// General operational entries about what's going on inside the application.
	InfoLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	ErrorLevel
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. highest level of severity.
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	}
	return ""
}

// Enabled returns true if the given level is at or above this level.
func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}

// GetLevel converts a level string into a logger Level value.
// returns an error if the input string does not match known values.
func GetLevel(levelStr string) (Level, error) {
	switch levelStr {
	case TraceLevel.String():
		return TraceLevel, nil
	case DebugLevel.String():
		return DebugLevel, nil
	case InfoLevel.String():
		return InfoLevel, nil
	case WarnLevel.String():
		return WarnLevel, nil
	case ErrorLevel.String():
		return ErrorLevel, nil
	case FatalLevel.String():
		return FatalLevel, nil
	}
	return InfoLevel, fmt.Errorf("unknown Level String: '%s', defaulting to InfoLevel", levelStr)
}
