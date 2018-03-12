package logger

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"
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
		caller := caller()
		timestamp := time.Now().Format("2006-01-02T15:04:05.000")
		a := []interface{}{timestamp, caller}
		a = append(a, args...)
		log.Printf("[DEBUG] [%s] [%s] "+format, a...)	}
}

//Infof ...
func (l *defaultLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	if l.level <= Info {
		caller := caller()
		timestamp := time.Now().Format("2006-01-02T15:04:05.000")
		a := []interface{}{timestamp, caller}
		a = append(a, args...)
		log.Printf("[INFO] [%s] [%s] "+format, a...)
	}
}

//Warnf ...
func (l *defaultLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	if l.level <= Warn {
		caller := caller()
		timestamp := time.Now().Format("2006-01-02T15:04:05.000")
		a := []interface{}{timestamp, caller}
		a = append(a, args...)
		log.Printf("[WARN] [%s] [%s] "+format, a...)	}
}

//Errorf ...
func (l *defaultLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	if l.level <= Error {
		caller := caller()
		timestamp := time.Now().Format("2006-01-02T15:04:05.000")
		a := []interface{}{timestamp, caller}
		a = append(a, args...)
		log.Printf("[ERROR] [%s] [%s] "+format, a...)
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

func caller() string {
	_, f, l, _ := runtime.Caller(3)
	return fmt.Sprintf("%s +%d", f, l)
}