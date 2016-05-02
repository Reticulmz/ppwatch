package main

import (
	"fmt"
	"github.com/mgutz/ansi"
	log "github.com/Sirupsen/logrus"
)

type LogFormatter struct {
	DisableColor bool

	PanicColor string
	FatalColor string
	ErrorColor string
	WarnColor  string
	InfoColor  string
	DebugColor string
}

func NewLogFormatter(disablecolor bool) *LogFormatter {
	return &LogFormatter {
		DisableColor: disablecolor,
		PanicColor: ansi.ColorCode("red"),
		FatalColor: ansi.ColorCode("red"),
		ErrorColor: ansi.ColorCode("red"),
		WarnColor:  ansi.ColorCode("yellow"),
		InfoColor:  ansi.ColorCode("blue"),
		DebugColor: ansi.ColorCode("magenta"),
	}
}

func (f *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	var color string
	var level string

	resetcolor := ansi.ColorCode("reset")

	switch entry.Level {
	case log.PanicLevel:
		color = f.PanicColor
		level = "PANIC"
	case log.FatalLevel:
		color = f.FatalColor
		level = "FATAL"
	case log.ErrorLevel:
		color = f.ErrorColor
		level = "ERROR"
	case log.WarnLevel:
		color = f.WarnColor
		level = "WARN"
	case log.InfoLevel:
		color = f.InfoColor
		level = "INFO"
	case log.DebugLevel:
		color = f.DebugColor
		level = "DEBUG"
	default:
		color = ""
		level = "?"
	}

	if f.DisableColor {
		color = ""
		resetcolor = ""
	}

	time := entry.Time.Format("15:04:05.000")
	out := fmt.Sprintf("%s%s %s ▶%s %s\n", color, time, level, resetcolor, entry.Message)

	return []byte(out), nil
}
