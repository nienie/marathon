package logger

import (
	"context"
	"log"
)

var (
	loggerImpl Logger
)

func init() {
	loggerImpl = newDefaultLogger()
}

//Level ...
type Level int

const (
	//Debug ...
	Debug Level = iota + 1
	//Info ...
	Info
	//Warn ...
	Warn
	//Error ...
	Error
)

//Logger ...
type Logger interface {

	//Debugf ...
	Debugf(ctx context.Context, format string, args ...interface{})

	//Infof ...
	Infof(ctx context.Context, format string, args ...interface{})

	//Warnf ...
	Warnf(ctx context.Context, format string, args ...interface{})

	//Errorf ...
	Errorf(ctx context.Context, format string, args ...interface{})

	//SetLevel ...
	SetLevel(Level)
}

type defaultLogger struct {
	level Level
}

func newDefaultLogger() *defaultLogger {
	return &defaultLogger{
		level: Debug,
	}
}

//Debugf ...
func (l *defaultLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	if l.level <= Debug {
		log.Printf("[DEBUG] "+format, args...)
	}
}

//Infof ...
func (l *defaultLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	if l.level <= Info {
		log.Printf("[INFO] "+format, args...)
	}
}

//Warnf ...
func (l *defaultLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	if l.level <= Warn {
		log.Printf("[WARN] "+format, args...)
	}
}

//Errorf ...
func (l *defaultLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	if l.level <= Error {
		log.Printf("[ERROR] "+format, args...)
	}
}

//SetLevel ...
func (l *defaultLogger) SetLevel(level Level) {
	l.level = level
}

//SetLogger ...
func SetLogger(l Logger) {
	if l != nil {
		loggerImpl = l
	}
}

//Debugf ...
func Debugf(ctx context.Context, format string, args ...interface{}) {
	loggerImpl.Debugf(ctx, format, args...)
}

//Infof ...
func Infof(ctx context.Context, format string, args ...interface{}) {
	loggerImpl.Infof(ctx, format, args...)
}

//Warnf ...
func Warnf(ctx context.Context, format string, args ...interface{}) {
	loggerImpl.Warnf(ctx, format, args...)
}

//Errorf ...
func Errorf(ctx context.Context, format string, args ...interface{}) {
	loggerImpl.Errorf(ctx, format, args...)
}

//SetLevel ...
func SetLevel(level Level) {
	loggerImpl.SetLevel(level)
}
