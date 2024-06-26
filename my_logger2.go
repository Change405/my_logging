package my_logging

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type Logger struct {
	debug    *log.Logger
	info     *log.Logger
	warning  *log.Logger
	critical *log.Logger
	logLevel string // Possible values are: DEBUG, INFO, WARNING, and CRITICAL
	logFile  string
	prefix   string
	testing  bool
}

func (l *Logger) init() {
	l.logLevel = "DEBUG"
	l.createLogs(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
}

// Used for testing
func (l *Logger) initTest(debugWriter io.Writer, infoWriter io.Writer, warningWriter io.Writer, criticalWriter io.Writer) {
	l.logLevel = "DEBUG"
	l.testing = true
	l.createLogs(debugWriter, infoWriter, warningWriter, criticalWriter)
}

// createLogs is an internal method used to initialize the log levels and set the outputs. It can be called again to change the log output locations.
func (l *Logger) createLogs(
	debugOutput io.Writer,
	infoOutput io.Writer,
	warningOutput io.Writer,
	errorOutput io.Writer) {

	l.debug = log.New(debugOutput,
		"DEBUG: ",
		log.Ldate|log.Ltime)

	l.info = log.New(infoOutput,
		"INFO: ",
		log.Ldate|log.Ltime)

	l.warning = log.New(warningOutput,
		"WARNING: ",
		log.Ldate|log.Ltime)

	l.critical = log.New(errorOutput,
		"CRITICAL: ",
		log.Ldate|log.Ltime)
}

// SetLogFile is used to add a log file to the log output.
// It takes in a string that will be used as the path to the log file.
// If the log file already exists the logs will be appended to the file, if not the log file will be created.
func (l *Logger) SetLogFile(logFilePath string) {
	l.logFile = logFilePath
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		l.Criticalf(err.Error())
	}

	multiWriter := io.MultiWriter(file, os.Stdout)
	multiWriterError := io.MultiWriter(file, os.Stderr)

	l.createLogs(multiWriter, multiWriter, multiWriter, multiWriterError)
}

// SetLogLevel will set the minimum log level to be logged.
// Possible values are DEBUG, INFO, WARNING, and CRITICAL.
func (l *Logger) SetLogLevel(newLogLevel string) {
	if newLogLevel == "DEBUG" || newLogLevel == "INFO" || newLogLevel == "WARNING" || newLogLevel == "CRITICAL" {
		l.logLevel = newLogLevel
	} else {
		l.Criticalf("%s is not a valid log level\n", newLogLevel)
	}
}

// SetPrefixf adds a prefix for all log messages. It additionally  acts as a wrapper for the fmt.Sprintf function.
// It takes in a prefix string and an optional amount of arguments which can be used to format the prefix.
// See fmt.Sprintf documentation for more detail. Using the formatting is optional.
func (l *Logger) SetPrefixf(prefix string, arguments ...any) {
	newPrefix := fmt.Sprintf(prefix, arguments...)
	newPrefix = newPrefix + " "
	l.prefix = newPrefix
}

// sendLog is an internal method that is used by all logging methods to send the log message.
func (l *Logger) sendLog(logLevel string, logMessage string, arguments ...any) {
	message := fmt.Sprintf(logMessage, arguments...)
	message = l.prefix + message
	switch logLevel {
	case "DEBUG":
		if l.logLevel == "DEBUG" {
			l.debug.Println(message)
		}
	case "INFO":
		if l.logLevel == "DEBUG" || l.logLevel == "INFO" {
			l.info.Println(message)
		}
	case "WARNING":
		if l.logLevel == "DEBUG" || l.logLevel == "INFO" || l.logLevel == "WARNING" {
			l.warning.Println(message)
		}
	case "CRITICAL":
		l.critical.Println(message)
		if !l.testing {
			os.Exit(1)
		}
	}
}

// checkError is an internal method used to check if an error is present before sending a log.
func (l *Logger) checkError(err error, logLevel string, logMessage string, arguments ...any) {
	if err != nil {
		developerMessage := fmt.Sprintf(logMessage, arguments...)
		errorMessage := err.Error()
		message := fmt.Sprintf("Error check failed. Dev message: '%s' Error message: '%s'", developerMessage, errorMessage)
		l.sendLog(logLevel, message)
	}
}

// Debugf writes a debug level log. It works as both a logger and a wrapper for the fmt.Sprintf function.
// It takes in a log message and optional arguments. These can be used to format the log message.
// See fmt.Sprintf documentation for more detail. Using the formatting is optional.
func (l *Logger) Debugf(logMessage string, arguments ...any) {
	l.sendLog("DEBUG", logMessage, arguments...)
}

// DebugfError writes a debug level log if the supplied error contains an error.
// If the supplied error is empty, no log message is created.
// It acts as an extension of Debugf and accepts the same formatted string and arguments.
func (l *Logger) DebugfError(err error, logMessage string, arguments ...any) {
	l.checkError(err, "DEBUG", logMessage, arguments...)
}

// Infof writes an info level log. It works as both a logger and a wrapper for the fmt.Sprintf function.
// It takes in a log message and optional arguments. These can be used to format the log message.
// See fmt.Sprintf documentation for more detail. Using the formatting is optional.
func (l *Logger) Infof(logMessage string, arguments ...any) {
	l.sendLog("INFO", logMessage, arguments...)
}

// InfofError writes an info level log if the supplied error contains an error.
// If the supplied error is empty, no log message is created.
// It acts as an extension of Infof and accepts the same formatted string and arguments.
func (l *Logger) InfofError(err error, logMessage string, arguments ...any) {
	l.checkError(err, "INFO", logMessage, arguments...)
}

// Warningf writes a warning level log. It works as both a logger and a wrapper for the fmt.Sprintf function.
// It takes in a log message and optional arguments. These can be used to format the log message.
// See fmt.Sprintf documentation for more detail. Using the formatting is optional.
func (l *Logger) Warningf(logMessage string, arguments ...any) {
	l.sendLog("WARNING", logMessage, arguments...)
}

// WarningfError writes a warning level log if the supplied error contains an error.
// If the supplied error is empty, no log message is created.
// It acts as an extension of Warningf and accepts the same formatted string and arguments.
func (l *Logger) WarningfError(err error, logMessage string, arguments ...any) {
	l.checkError(err, "WARNING", logMessage, arguments...)
}

// Criticalf writes a critical level log. Calling Criticalf will end the program.
// Criticalf should only be called to log an error that can not be recovered from.
// It works as both a logger and a wrapper for the fmt.Sprintf function.
// It takes in a log message and optional arguments. These can be used to format the log message.
// See fmt.Sprintf documentation for more detail. Using the formatting is optional.
func (l *Logger) Criticalf(logMessage string, arguments ...any) {
	l.sendLog("CRITICAL", logMessage, arguments...)
}

// CriticalfError writes a critical level log if the supplied error contains an error.
// If the supplied error is empty, no log message is created.
// It acts as an extension of Criticalf and accepts the same formatted string and arguments.
func (l *Logger) CriticalfError(err error, logMessage string, arguments ...any) {
	l.checkError(err, "CRITICAL", logMessage, arguments...)
}

type Test struct {
	Buffer []byte
}

func (t *Test) Write(p []byte) (n int, err error) {
	t.Buffer = append(t.Buffer, p...)
	return 0, nil
}

// CreateLogger creates a new logger and returns a pointer to it.
// It is recommended to use this method to create a new logger.
func CreateLogger() *Logger {
	var log Logger
	log.init()
	return &log
}

func main() {
	var log_test Logger
	//var log_output Test

	print := fmt.Println

	log_test.init()
	//log_test.SetPrefixf("TEST PREFIX")
	log_test.Debugf("test")
	print(log_test.logLevel)
	//print(string(log_output.Buffer))

	err := errors.New("test error")
	log_test.DebugfError(err, "Shit went down")
	//print(string(log_output.Buffer))
}
