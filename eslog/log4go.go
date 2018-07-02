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
	logType  int       // 0 Rec, 1:rawText
	level    level_t   // The log level
	created  time.Time // The time at which the log message was created (nanoseconds)
	source   string    // The message source
	message  string    // The log message
	category string    // The category
}

type LogModuleInfo struct{
	name string
	level level_t
}

//
// Logger object
//

type Logger struct {
	modules map[string]*LogModuleInfo
	logWriters map[string]LogWriter
	telnetWriter *TelnetLogWriter
}


func (this *Logger) AddWriter(writer LogWriter) {
	this.logWriters[writer.Name()] = writer
	if writer.Name() == "TelnetLogWriter" {
		this.telnetWriter = writer.(*TelnetLogWriter)
		this.telnetWriter.RegCommand("mlist", this, dbgModuleList, "Show all modules.")
	}
}

func (this *Logger) AddModule(name string, lvl level_t) {
	this.modules[name] = &LogModuleInfo{
		name:name,
	    level:lvl,
	}
}

// Set print level
func (this *Logger) SetLevel(which string, level level_t) {
	if module, ok := this.modules[which]; ok {
		module.level = level
	}
}


func NewLogger() *Logger {
	return &Logger{
		modules: make(map[string]*LogModuleInfo),
		logWriters:make(map[string]LogWriter),
    }
}

func (this *Logger) Close() {
	for _, v := range this.logWriters {
		v.Close()
	}
}

//
//Internal functions
//

func (this *Logger) RawPrintf(format string, args ...interface{}) {

	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	rec := &LogRecord{
		logType:  1,
		message:  msg,
	}

	// Write log
	for _, filter := range this.logWriters {
		filter.LogWrite(rec)
	}
}

func (this *Logger) levelPrintf(calldep int, moduleName string, level level_t, format string, args ...interface{}) {

	//Check module
	if _, ok := this.modules[moduleName]; !ok {
		return
	}

	// Check Level
	m := this.modules[moduleName]
	if level > m.level {
		return
	}

	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	_, filename, line, ok := runtime.Caller(calldep)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", path.Base(filename), line)
	}

	rec := &LogRecord{
		level:    level_t(level),
		created:  time.Now(),
		source:   src,
		message:  msg,
		category: moduleName,
	}

	// Write log
	for _, writer := range this.logWriters {
		writer.LogWrite(rec)
	}
}

//
// Utils
//

func (this *Logger) Fatal(name string, format string, args ...interface{}) {
	this.levelPrintf(2, name, FATAL, format, args...)
}

func (this *Logger) Debug(name string, format string, args ...interface{}) {
	this.levelPrintf(2, name, DEBUG, format, args...)
}

func (this *Logger) Error(name string, format string, args ...interface{}) {
	this.levelPrintf(2, name, ERROR, format, args...)
}

func (this *Logger) Warn(name string, format string, args ...interface{}) {
	this.levelPrintf(2, name, WARN, format, args...)
}

func (this *Logger) Info(name string, format string, args ...interface{}) {
	this.levelPrintf(2, name, INFO, format, args...)
}

func (this *Logger) Trace(name string, format string, args ...interface{}) {
	this.levelPrintf(2, name, TRACE, format, args...)
}


//
// Commands
//

func dbgModuleList(args ...interface{}) {
	c := args[0].(*Logger)

	for name, m := range c.modules {
		outText:= fmt.Sprintf("%s:level=%s\n", name, m.level.String())
		c.RawPrintf(outText)
	}
}