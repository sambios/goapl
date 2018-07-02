package eslog

import (
	"testing"
	"time"
)

func NoTestLog4go(t *testing.T) {
	logWriter := DefaultConsoleLogWriter()

	logger := NewLogger()
	logger.AddWriter(logWriter)
	//logger.AddFilter("Test", ERROR, NewFileLogWriter("./test.log", true))
	logger.Error("Test", "ERR.12345679")
	//logger.Printf("Test", CRIT, "CRIT.%d,%d", 1, 2)
	logger.Fatal("Test", "FATAL:%d,%d", 3, 4)
	logger.Debug("Test", "DEBUG:%d,%d", 5, 6)
	logger.Trace("Test", "TRACE:%d,%d", 5, 6)

	// Set Log Level
	logger.SetLevel("Test", TRACE)
	logger.Error("Test", "ERR.12345679")

	logger.Fatal("Test", "FATAL:%d,%d", 3, 4)
	logger.Debug("Test", "DEBUG:%d,%d", 5, 6)
	logger.Trace("Test", "TRACE:%d,%d", 5, 6)

	logger.Close()
}

func TestTelnetLogWriter(t *testing.T) {

	logwriter := NewTelnetLogWriter(3000)

	log := NewLogger()
	log.AddWriter(logwriter)

	log.AddModule("Test", TRACE)
	log.RawPrintf("RawText:%s\n", "1111111111111")

	log.Fatal("Test", "FATAL:%d,%d", 3, 4)
	log.Debug("Test", "DEBUG:%d,%d", 5, 6)
	log.Trace("Test", "TRACE:%d,%d", 5, 6)

	log.RawPrintf("RawText:%s", "222222 ")
	log.RawPrintf("RawText:%s", "333333 ")
	log.RawPrintf("RawText:%s", "444444\n")

	time.Sleep(10000 * time.Second)

	log.Close()

}
