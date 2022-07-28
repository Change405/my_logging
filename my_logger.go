//Package my_logging expands on the standard logging package in order to add addition functionality like log levels.
package my_logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	debug    *log.Logger
	info     *log.Logger
	warning  *log.Logger
	critical *log.Logger

	logLevel string // Possible values are: DEBUG, INFO, WARNING, and CRITICAL
)

func init() {
	logLevel = "DEBUG"
	createLogs(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
}

// createLogs is an internal function used to initialize the log levels and set the outputs. It can be called again to change the log output locations.
func createLogs(
	debugOutput io.Writer,
	infoOutput io.Writer,
	warningOutput io.Writer,
	errorOutput io.Writer) {

	debug = log.New(debugOutput,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	info = log.New(infoOutput,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	warning = log.New(warningOutput,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	critical = log.New(errorOutput,
		"CRITICAL: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// SetLogFile is used to add a log file to the log output.
// It takes in a string that will be used as the path to the log file.
// If the log file already exists the logs will be appended to the file, if not the log file will be created.
func SetLogFile(logFilePath string) {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Criticalf(err.Error())
	}

	multiWriter := io.MultiWriter(file, os.Stdout)
	multiWriterError := io.MultiWriter(file, os.Stderr)

	createLogs(multiWriter, multiWriter, multiWriter, multiWriterError)
}

// SetLogLevel will set the minimum log level to be logged.
// Possible values are DEBUG, INFO, WARNING, and CRITICAL.
func SetLogLevel(newLogLevel string) {
	if newLogLevel == "DEBUG" || newLogLevel == "INFO" || newLogLevel == "WARNING" || newLogLevel == "CRITICAL" {
		logLevel = newLogLevel
	} else {
		// TODO: Fix this so that it can take in an argument
		Criticalf("%s is not a valid log level", newLogLevel)
	}
}

// Debugf writes a debug level log. It works as both a logger and a wrapper for the fmt.Sprintf function.
// It takes in a log message and optional arguments. These can be used to format the log message.
// See fmt.Sprintf documentation for more detail. Using the formatting is optional.
func Debugf(logMessage string, arguments ...any) {
	if logLevel == "DEBUG" {
		message := fmt.Sprintf(logMessage, arguments...)
		debug.Println(message)
	}
}

// DebugfError writes a debug level log if the supplied error contains an error.
// If the supplied error is empty, no log message is created.
// It acts as an extension of Debugf and accepts the same formatted string and arguments.
func DebugfError(err error, logMessage string, arguments ...any) {
	if logLevel == "DEBUG" {
		if err != nil {
			Debugf(logMessage, arguments...)
			Debugf(err.Error())
		}
	}
}

// Infof writes an info level log. It works as both a logger and a wrapper for the fmt.Sprintf function.
// It takes in a log message and optional arguments. These can be used to format the log message.
// See fmt.Sprintf documentation for more detail. Using the formatting is optional.
func Infof(logMessage string, arguments ...any) {
	if logLevel == "DEBUG" || logLevel == "INFO" {
		message := fmt.Sprintf(logMessage, arguments...)
		info.Println(message)
	}
}

// InfofError writes an info level log if the supplied error contains an error.
// If the supplied error is empty, no log message is created.
// It acts as an extension of Infof and accepts the same formatted string and arguments.
func InfofError(err error, logMessage string, arguments ...any) {
	if logLevel == "DEBUG" || logLevel == "INFO" {
		if err != nil {
			Infof(logMessage, arguments...)
			Infof(err.Error())
		}
	}
}

// Warningf writes a warning level log. It works as both a logger and a wrapper for the fmt.Sprintf function.
// It takes in a log message and optional arguments. These can be used to format the log message.
// See fmt.Sprintf documentation for more detail. Using the formatting is optional.
func Warningf(logMessage string, arguments ...any) {
	if logLevel == "DEBUG" || logLevel == "INFO" || logLevel == "WARNING" {
		message := fmt.Sprintf(logMessage, arguments...)
		warning.Println(message)
	}
}

// WarningfError writes a warning level log if the supplied error contains an error.
// If the supplied error is empty, no log message is created.
// It acts as an extension of Warningf and accepts the same formatted string and arguments.
func WarningfError(err error, logMessage string, arguments ...any) {
	if logLevel == "DEBUG" || logLevel == "INFO" || logLevel == "WARNING" {
		if err != nil {
			Warningf(logMessage, arguments...)
			Warningf(err.Error())
		}
	}
}

// Criticalf writes a critical level log. Calling Criticalf will end the program.
// Criticalf should only be called to log an error that can not be recovered from.
// It works as both a logger and a wrapper for the fmt.Sprintf function.
// It takes in a log message and optional arguments. These can be used to format the log message.
// See fmt.Sprintf documentation for more detail. Using the formatting is optional.
func Criticalf(logMessage string, arguments ...any) {
	message := fmt.Sprintf(logMessage, arguments...)
	critical.Fatalln(message)
}

// CriticalfError writes a critical level log if the supplied error contains an error.
// If the supplied error is empty, no log message is created.
// It acts as an extension of Criticalf and accepts the same formatted string and arguments.
func CriticalfError(err error, logMessage string, arguments ...any) {
	if err != nil {
		message := fmt.Sprintf(logMessage, arguments...)
		message = message + "\n" + err.Error()
		Criticalf(message)
	}
}
