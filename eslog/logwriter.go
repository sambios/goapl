package eslog

/****** LogWriter ******/

// This is an interface for anything that should be able to write logs
type LogWriter interface {
	// Get Name of writer
	Name() string

	// This will be called to log a LogRecord message.
	LogWrite(rec *LogRecord)

	// This should clean up anything lingering about the LogWriter, as it is called before
	// the LogWriter is removed.  LogWrite should not be called after Close.
	Close()
}
