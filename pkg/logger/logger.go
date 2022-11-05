package logger

import (
	"io"
	"strings"
)

type levelType uint8

const (
	levelOff levelType = iota
	levelError
	levelWarning
	levelInfo
	levelDebug
)

type LogInterface interface {
	Error(string)
	Warning(string)
	Info(string)
	Debug(string)
}

type impl struct {
	writer io.Writer
	level  levelType
}

func New(level string, writer io.Writer) LogInterface {
	var uintLevel levelType
	switch strings.ToLower(level) {
	case "error":
		uintLevel = levelError
	case "warning":
		uintLevel = levelWarning
	case "info":
		uintLevel = levelInfo
	case "debug":
		uintLevel = levelDebug
	default:
		uintLevel = levelOff
	}
	return &impl{level: uintLevel, writer: writer}
}

func (l *impl) Error(msg string) {
	l.log(levelError, msg)
}

func (l *impl) Warning(msg string) {
	l.log(levelWarning, msg)
}

func (l *impl) Info(msg string) {
	l.log(levelInfo, msg)
}

func (l *impl) Debug(msg string) {
	l.log(levelDebug, msg)
}

func (l *impl) log(level levelType, msg string) {
	if l.level >= level {
		_, _ = l.writer.Write([]byte(msg + "\n"))
	}
}
