package log

import (
	"time"

	"github.com/Sirupsen/logrus"
)

// Config defines how we want our logger to behave
type Config struct {
	Level     string
	Formatter string
}

// Default values
const (
	DefaultLevel     = "info"
	DefaultFormatter = "text"
)

// DefaultConfig provides a a set of config that is
// suitable for use in most cases
var DefaultConfig = Config{
	Level:     DefaultLevel,
	Formatter: DefaultFormatter,
}

var levels = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"fatal": logrus.FatalLevel,
}

var formatters = map[string]logrus.Formatter{
	"json": &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	},
	"text": &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	},
}

// NewDefault returns a logger with the
// Default configuration
func NewDefault() *logrus.Logger {
	return New(DefaultConfig)
}

// New returns a new logrus.Logger
func New(conf Config) *logrus.Logger {
	l := logrus.New()

	if lvl, ok := levels[conf.Level]; ok {
		l.Level = lvl
	} else {
		l.Level = levels[DefaultLevel]
	}

	if fmter, ok := formatters[conf.Formatter]; ok {
		l.Formatter = fmter
	} else {
		l.Formatter = formatters[DefaultFormatter]
	}

	return l
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	logrus.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	logrus.Error(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}
