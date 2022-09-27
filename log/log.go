package log

import (
	"context"
	"fmt"
	"time"
)

type Logger interface {
	Debug(i ...interface{})
	Debugf(format string, i ...interface{})
	Info(i ...interface{})
	Infof(format string, i ...interface{})
	Warn(i ...interface{})
	Warnf(format string, i ...interface{})
	Error(i ...interface{})
	Errorf(format string, i ...interface{})
	Sync() chan struct{}
}

type Format uint8

type Level uint8

const (
	PlainFormat Format = iota
	JsonFormat

	LevelDebug Level = iota - 2
	LevelInfo
	LevelWarn
	LevelError
)

var levelNames = map[Level]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
}

type entry struct {
	level     Level
	timestamp time.Time
	data      string
}

type defaultLogger struct {
	level   Level
	format  Format
	entries chan *entry
	closed  bool
	sync    chan struct{}
}

func newDefaultLogger(ctx context.Context, level Level) Logger {
	logger := &defaultLogger{
		level:   level,
		format:  PlainFormat,
		entries: make(chan *entry, 1024),
		closed:  false,
		sync:    make(chan struct{}),
	}
	go logger.start(ctx)
	return logger
}

func (l *defaultLogger) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			l.closed = true
			close(l.entries)
			for entry := range l.entries {
				l.handle(entry)
			}
			l.sync <- struct{}{}
			return
		case entry := <-l.entries:
			l.handle(entry)
		}
	}
}

func (l *defaultLogger) handle(msg *entry) {
	fmt.Printf("%s [%s]: %s\n", msg.timestamp.Format(time.RFC3339), levelNames[msg.level], msg.data)
}

func (l *defaultLogger) log(level Level, i ...interface{}) {
	if l.closed {
		return
	}
	if level < l.level {
		return
	}
	l.entries <- &entry{
		timestamp: time.Now().UTC(),
		level:     level,
		data:      fmt.Sprintf("%v", i...),
	}
}

func (l *defaultLogger) Debug(i ...interface{}) {
	l.log(LevelDebug, i...)
}

func (l *defaultLogger) Debugf(format string, i ...interface{}) {
	l.Debug(fmt.Sprintf(format, i...))
}

func (l *defaultLogger) Info(i ...interface{}) {
	l.log(LevelInfo, i...)
}

func (l *defaultLogger) Infof(format string, i ...interface{}) {
	l.Info(fmt.Sprintf(format, i...))
}

func (l *defaultLogger) Warn(i ...interface{}) {
	l.log(LevelWarn, i...)
}

func (l *defaultLogger) Warnf(format string, i ...interface{}) {
	l.Warn(fmt.Sprintf(format, i...))
}

func (l *defaultLogger) Error(i ...interface{}) {
	l.log(LevelError, i...)
}

func (l *defaultLogger) Errorf(format string, i ...interface{}) {
	l.Error(fmt.Sprintf(format, i...))
}

func (l *defaultLogger) Sync() chan struct{} {
	return l.sync
}

func Debug(i ...interface{}) {
	loggerInstance.Debug(i...)
}

func Debugf(format string, i ...interface{}) {
	loggerInstance.Debugf(format, i...)
}

func Info(i ...interface{}) {
	loggerInstance.Info(i...)
}

func Infof(format string, i ...interface{}) {
	loggerInstance.Infof(format, i...)
}

func Warn(i ...interface{}) {
	loggerInstance.Warn(i...)
}

func Warnf(format string, i ...interface{}) {
	loggerInstance.Warnf(format, i...)
}

func Error(i ...interface{}) {
	loggerInstance.Error(i...)
}

func Errorf(format string, i ...interface{}) {
	loggerInstance.Errorf(format, i...)
}

func Sync() chan struct{} {
	return loggerInstance.Sync()
}

var loggerInstance Logger

func Init(ctx context.Context, level Level) {
	if loggerInstance == nil {
		loggerInstance = newDefaultLogger(ctx, level)
	}
}
