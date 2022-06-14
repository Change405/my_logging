//Package github.com/change405/my_logging expands on the standard logging package in order to add addition functionality like log levels.
package my_logging

import (
	"io"
	"log"
	"os"
)

var (
	debug    *log.Logger
	info     *log.Logger
	warning  *log.Logger
	critical *log.Logger

	log_level string // Possible values are: DEBUG, INFO, WARNING, and CRITICAL
)

func init() {
	log_level = "DEBUG"
	create_logs(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
}

// create_logs is an internal function used to initialize the log levels and set the outputs. It can be called again to change the log output locations.
func create_logs(
	debug_output io.Writer,
	info_output io.Writer,
	warning_output io.Writer,
	error_output io.Writer) {

	debug = log.New(debug_output,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	info = log.New(info_output,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	warning = log.New(warning_output,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	critical = log.New(error_output,
		"CRITICAL: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// Set_log_file is used to add a log file to the log output.
// It takes in a string that will be used as the path to the log file.
// If the log file already exists the logs will be appended to the file, if not the log file will be created.
func Set_log_file(log_file_path string) {
	file, err := os.OpenFile(log_file_path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Critical(err.Error())
	}

	multi_writer := io.MultiWriter(file, os.Stdout)
	multi_writer_error := io.MultiWriter(file, os.Stderr)

	create_logs(multi_writer, multi_writer, multi_writer, multi_writer_error)
}

// Set_log_level will set the minimum log level to be logged.
// Possible values are DEBUG, INFO, WARNING, and CRITICAL.
func Set_log_level(new_log_level string) {
	if new_log_level == "DEBUG" || new_log_level == "INFO" || new_log_level == "WARNING" || new_log_level == "CRITICAL" {
		log_level = new_log_level
	} else {
		Critical("%s is not a valid log level")
	}
}

// Debug writes a debug level log.
func Debug(log_message string) {
	if log_level == "DEBUG" {
		debug.Println(log_message)
	}
}

// Info writes a info level log.
func Info(log_message string) {
	if log_level == "DEBUG" || log_level == "INFO" {
		info.Println(log_message)
	}
}

// Warning writes a warning level log.
func Warning(log_message string) {
	if log_level == "DEBUG" || log_level == "INFO" || log_level == "WARNING" {
		warning.Println(log_message)
	}
}

// Critical writes a critical level log.
// Calling Critical will end the program.
// Critical should only be called to log an error that can not be recovered from.
func Critical(log_message string) {
	critical.Fatalln(log_message)
}
