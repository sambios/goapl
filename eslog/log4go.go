package eslog

import (
	"fmt"
	"runtime"
	"time"
	"path"
)

const (
	L4G_VERSION = "eslog-v0.1"
	L4G_MAJOR   = 0
	L4G_MINOR   = 1
	L4G_BUILD   = 1
)

type level_t int

const (
	ALL = iota
	FATAL
	CRIT
	ERROR
	WARN
	DEBUG
	INFO
	TRACE
)

var (
	levelStrings = [...]string{"ALL", "FATAL", "CRIT", "ERROR", "WARN", "DEBUG", "INFO", "TRACE"}
)

func (l level_t) String() string {
	if l < 0 || int(l) > len(levelStrings) {
		return "UNKNOWN"
	}

	return levelStrings[int(l)]
}

/****** Variables ******/
var (
	// LogBufferLength specifies how many log messages a particular eslog
	// logger can buffer at a time before writing them.
	LogBufferLength = 32
)

// A LogRecord contains all of the pertinent information for each message
type LogRecord struct {
	level    level_t   // The log level
	created  time.Time // The time at which the log message was created (nanoseconds)
	source   string    // The message source
	message  string    // The log message
	category string    // The category
}

/****** Logger ******/
type LogWriterFilter struct {
	level level_t
	LogWriter
}

/****** Module Logger ******/

type ModuleLogger struct {
	name    string
	filters map[string]*LogWriterFilter
}

func (m *ModuleLogger) AddWriter(lvl level_t, writer LogWriter) {

	filter := LogWriterFilter{
		lvl,
		writer,
	}

	m.filters[writer.Name()] = &filter
}

func (this *ModuleLogger) Close() {
	for _, filter := range this.filters {
		filter.Close()
	}
}

// Set print level
func (m *ModuleLogger) SetLevel(which string, level level_t) {
	if filter, ok := m.filters[which]; ok {
		filter.level = level
	}
}

func (m *ModuleLogger) Printf(lvl level_t, src, format string, args ...interface{}) {
	skip := true
	for _, filter := range m.filters {
		if lvl > filter.level {
			continue
		}

		skip = false
	}

	if skip {
		return
	}

	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	rec := &LogRecord{
		level:    level_t(lvl),
		created:  time.Now(),
		source:   src,
		message:  msg,
		category: m.name,
	}

	// Write log
	for _, filter := range m.filters {
		if lvl > filter.level {
			continue
		}

		filter.LogWrite(rec)
	}
}

//
// Logger object
//

type Logger struct {
	modules map[string]*ModuleLogger
}

func NewLogger() *Logger {
	return &Logger{modules: make(map[string]*ModuleLogger)}
}

func (this *Logger) Close() {
	for _, v := range this.modules {
		v.Close()
	}
}

func (this *Logger) Module(name string) *ModuleLogger {
	if m, ok := this.modules[name]; ok {
		return m
	}

	return nil
}

// Add Log Filter
func (this *Logger) AddFilter(moduleName string, lvl level_t, writer LogWriter) *ModuleLogger {

	if m, ok := this.modules[moduleName]; ok {
		m.AddWriter(lvl, writer)
		return m
	}

	m := ModuleLogger{
		name:    moduleName,
		filters: make(map[string]*LogWriterFilter),
	}

	m.AddWriter(lvl, writer)
	this.modules[moduleName] = &m

	return &m
}

func (this *Logger) Printf(name string, level level_t, format string, args ...interface{}) {

	if _, ok := this.modules[name]; !ok {
		return
	}

	moduleLogger, _ := this.modules[name]

	_, filename, line, ok := runtime.Caller(1)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", path.Base(filename), line)
	}

	moduleLogger.Printf(level, src, format, args...)
}

//
// Utils
//

func (this *Logger) Fatal(name string, format string, args ...interface{}) {
	if _, ok := this.modules[name]; !ok {
		return
	}

	moduleLogger, _ := this.modules[name]

	_, filename, line, ok := runtime.Caller(1)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", path.Base(filename), line)
	}

	moduleLogger.Printf(FATAL, src, format, args...)
}

func (this *Logger) Error(name string, format string, args ...interface{}) {
	if _, ok := this.modules[name]; !ok {
		return
	}

	moduleLogger, _ := this.modules[name]

	_, filename, line, ok := runtime.Caller(1)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", path.Base(filename), line)
	}

	moduleLogger.Printf(ERROR, src, format, args...)
}

func (this *Logger) Warn(name string, format string, args ...interface{}) {
	if _, ok := this.modules[name]; !ok {
		return
	}

	moduleLogger, _ := this.modules[name]

	_, filename, line, ok := runtime.Caller(1)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", path.Base(filename), line)
	}

	moduleLogger.Printf(WARN, src, format, args...)
}

func (this *Logger) Info(name string, format string, args ...interface{}) {
	if _, ok := this.modules[name]; !ok {
		return
	}

	moduleLogger, _ := this.modules[name]

	_, filename, line, ok := runtime.Caller(1)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", path.Base(filename), line)
	}

	moduleLogger.Printf(INFO, src, format, args...)
}

func (this *Logger) Trace(name string, format string, args ...interface{}) {
	if _, ok := this.modules[name]; !ok {
		return
	}

	moduleLogger, _ := this.modules[name]

	_, filename, line, ok := runtime.Caller(1)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", path.Base(filename), line)
	}

	moduleLogger.Printf(TRACE, src, format, args...)
}
