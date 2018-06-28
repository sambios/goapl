package eslog

import "testing"

func TestLog4go(t *testing.T) {
	logWriter := DefaultConsoleLogWriter()

	logger := NewLogger()
	logger.AddFilter("Test", ERROR, logWriter)
	//logger.AddFilter("Test", ERROR, NewFileLogWriter("./test.log", true))
	logger.Error("Test", "ERR.12345679")
	logger.Printf("Test", CRIT, "CRIT.%d,%d", 1, 2)
	logger.Fatal("Test", "FATAL:%d,%d", 3, 4)
	logger.Printf("Test", DEBUG, "DEBUG:%d,%d", 5, 6)
	logger.Printf("Test", TRACE, "TRACE:%d,%d", 5, 6)

	// Set Log Level
	logger.Module("Test").SetLevel(logWriter.Name(), TRACE)
	logger.Printf("Test", ERROR, "ERR.12345679")
	logger.Printf("Test", CRIT, "CRIT.%d,%d", 1, 2)
	logger.Printf("Test", FATAL, "FATAL:%d,%d", 3, 4)
	logger.Printf("Test", DEBUG, "DEBUG:%d,%d", 5, 6)
	logger.Printf("Test", TRACE, "TRACE:%d,%d", 5, 6)

	logger.Close()
}
