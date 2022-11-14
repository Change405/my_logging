package my_logging

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

type testBuffer struct {
	Buffer []byte
}

func (t *testBuffer) Write(p []byte) (n int, err error) {
	t.Buffer = append(t.Buffer, p...)
	return 0, nil
}

func removeDateAndTime(logMessage string) string {
	tempMessage := strings.Split(logMessage, " ")
	tempMessage = append(tempMessage[:1], tempMessage[3:]...)
	logMessage = strings.Join(tempMessage, " ")
	return logMessage
}

func checkLogOutput(t *testing.T, buffer []byte, logLevel string, logMessage string) {
	expectedLogMessage := logLevel + ": " + logMessage + "\n"

	if len(buffer) == 0 {
		t.Fatalf("%s level log was expected but not received", logLevel)
	}
	rawMessage := string(buffer)
	finalMessage := removeDateAndTime(rawMessage)
	if finalMessage != expectedLogMessage {
		t.Fatalf("Expected output: '%s' Received output: '%s'", expectedLogMessage, finalMessage)
	}
}

func checkNoLogOutput(t *testing.T, buffer []byte, logLevel string) {
	if len(buffer) != 0 {
		t.Fatalf("%s log was received when it wasn't expected", logLevel)
	}
}

func checkEachLogLevel(t *testing.T, logLevel string) {
	var debugBuffer testBuffer
	var infoBuffer testBuffer
	var warningBuffer testBuffer
	var criticalBuffer testBuffer

	var log Logger
	log.initTest(&debugBuffer, &infoBuffer, &warningBuffer, &criticalBuffer)
	log.SetLogLevel(logLevel)

	testLogMessage := "This is a %s of the %s"
	arguments := []string{"test", "logging."}
	expectedLogMessage := fmt.Sprintf(testLogMessage, arguments[0], arguments[1])

	log.Debugf(testLogMessage, arguments[0], arguments[1])
	log.Infof(testLogMessage, arguments[0], arguments[1])
	log.Warningf(testLogMessage, arguments[0], arguments[1])
	log.Criticalf(testLogMessage, arguments[0], arguments[1])

	// Check DEBUG
	if logLevel == "DEBUG" {
		checkLogOutput(t, debugBuffer.Buffer, "DEBUG", expectedLogMessage)
	} else {
		checkNoLogOutput(t, debugBuffer.Buffer, "DEBUG")
	}

	// Check INFO
	if logLevel == "INFO" || logLevel == "DEBUG" {
		checkLogOutput(t, infoBuffer.Buffer, "INFO", expectedLogMessage)
	} else {
		checkNoLogOutput(t, infoBuffer.Buffer, "INFO")
	}

	// Check WARNING
	if logLevel == "WARNING" || logLevel == "INFO" || logLevel == "DEBUG" {
		checkLogOutput(t, warningBuffer.Buffer, "WARNING", expectedLogMessage)
	} else {
		checkNoLogOutput(t, warningBuffer.Buffer, "WARNING")
	}

	// Check CRITICAL
	checkLogOutput(t, criticalBuffer.Buffer, "CRITICAL", expectedLogMessage)
}

func TestStandardLoggingOutput(t *testing.T) {

	logLevels := []string{"DEBUG", "INFO", "WARNING", "CRITICAL"}

	for _, v := range logLevels {
		checkEachLogLevel(t, v)
	}
}

func checkEachLogLevelFileOutput(t *testing.T, logLevel string) {
	var logBuffer testBuffer
	var log Logger
	log.initTest(&logBuffer, &logBuffer, &logBuffer, &logBuffer)

	logFile := "./test_log.txt"
	// Remove old log file
	os.Remove(logFile)
	log.SetLogFile(logFile)

	log.SetLogLevel(logLevel)

	testLogMessage := "This is a %s of the %s"
	arguments := []string{"test", "logging."}
	expectedLogMessage := fmt.Sprintf(testLogMessage, arguments[0], arguments[1])

	log.Debugf(testLogMessage, arguments[0], arguments[1])
	log.Infof(testLogMessage, arguments[0], arguments[1])
	log.Warningf(testLogMessage, arguments[0], arguments[1])
	log.Criticalf(testLogMessage, arguments[0], arguments[1])

	fileData, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Cant open log file")
	}

	fileDataString := string(fileData)
	logLines := strings.Split(fileDataString, "\n")
	var formattedLogLines []string
	for _, i := range logLines {
		if len(i) != 0 {
			formattedLine := removeDateAndTime(i)
			formattedLogLines = append(formattedLogLines, formattedLine)
		}
	}

	expectedDebugLog := "DEBUG: " + expectedLogMessage
	expectedInfoLog := "INFO: " + expectedLogMessage
	expectedWarningLog := "WARNING: " + expectedLogMessage
	expectedCriticalLog := "CRITICAL: " + expectedLogMessage

	if logLevel == "DEBUG" {
		expectedOutput := []string{expectedDebugLog, expectedInfoLog, expectedWarningLog, expectedCriticalLog}
		for i, v := range expectedOutput {
			if v != formattedLogLines[i] {
				t.Fatalf("Expected output: '%s' Received Output: '%s'", v, formattedLogLines[i])
			}
		}
	} else if logLevel == "INFO" {
		expectedOutput := []string{expectedInfoLog, expectedWarningLog, expectedCriticalLog}
		for i, v := range expectedOutput {
			if v != formattedLogLines[i] {
				t.Fatalf("Expected output: '%s' Received Output: '%s'", v, formattedLogLines[i])
			}
		}
	} else if logLevel == "WARNING" {
		expectedOutput := []string{expectedWarningLog, expectedCriticalLog}
		for i, v := range expectedOutput {
			if v != formattedLogLines[i] {
				t.Fatalf("Expected output: '%s' Received Output: '%s'", v, formattedLogLines[i])
			}
		}
	} else if logLevel == "CRITICAL" {
		expectedOutput := []string{expectedCriticalLog}
		for i, v := range expectedOutput {
			if v != formattedLogLines[i] {
				t.Fatalf("Expected output: '%s' Received Output: '%s'", v, formattedLogLines[i])
			}
		}
	}
	os.Remove(logFile)
}

func TestFileLoggingOutput(t *testing.T) {
	logLevels := []string{"DEBUG", "INFO", "WARNING", "CRITICAL"}

	for _, v := range logLevels {
		checkEachLogLevelFileOutput(t, v)
	}

}

func checkEachLogLevelError(t *testing.T, logLevel string) {
	var debugBuffer testBuffer
	var infoBuffer testBuffer
	var warningBuffer testBuffer
	var criticalBuffer testBuffer

	var err = errors.New("Test error")

	var log Logger
	log.initTest(&debugBuffer, &infoBuffer, &warningBuffer, &criticalBuffer)
	log.SetLogLevel(logLevel)

	testLogMessage := "This is a %s of the %s"
	arguments := []string{"test", "logging."}
	part1LogMessage := fmt.Sprintf(testLogMessage, arguments[0], arguments[1])
	part2LogMessage := "Error check failed. Dev message: '%s' Error message: '%s'"
	expectedLogMessage := fmt.Sprintf(part2LogMessage, part1LogMessage, err.Error())

	log.DebugfError(err, testLogMessage, arguments[0], arguments[1])
	log.InfofError(err, testLogMessage, arguments[0], arguments[1])
	log.WarningfError(err, testLogMessage, arguments[0], arguments[1])
	log.CriticalfError(err, testLogMessage, arguments[0], arguments[1])

	// Check DEBUG
	if logLevel == "DEBUG" {
		checkLogOutput(t, debugBuffer.Buffer, "DEBUG", expectedLogMessage)
	} else {
		checkNoLogOutput(t, debugBuffer.Buffer, "DEBUG")
	}

	// Check INFO
	if logLevel == "INFO" || logLevel == "DEBUG" {
		checkLogOutput(t, infoBuffer.Buffer, "INFO", expectedLogMessage)
	} else {
		checkNoLogOutput(t, infoBuffer.Buffer, "INFO")
	}

	// Check WARNING
	if logLevel == "WARNING" || logLevel == "INFO" || logLevel == "DEBUG" {
		checkLogOutput(t, warningBuffer.Buffer, "WARNING", expectedLogMessage)
	} else {
		checkNoLogOutput(t, warningBuffer.Buffer, "WARNING")
	}

	// Check CRITICAL
	checkLogOutput(t, criticalBuffer.Buffer, "CRITICAL", expectedLogMessage)
}

func TestErrorLogging(t *testing.T) {

	// Check Debug
	checkEachLogLevelError(t, "DEBUG")

	// Check Info
	checkEachLogLevelError(t, "INFO")

	// Check Warning
	checkEachLogLevelError(t, "WARNING")

	// Check Critical
	checkEachLogLevelError(t, "CRITICAL")
}

func TestNilErrorLogging(t *testing.T) {
	var logBuffer testBuffer
	var log Logger
	log.initTest(&logBuffer, &logBuffer, &logBuffer, &logBuffer)

	testLogMessage := "This is a %s of the %s"
	arguments := []string{"test", "logging."}

	log.DebugfError(nil, testLogMessage, arguments[0], arguments[1])
	log.InfofError(nil, testLogMessage, arguments[0], arguments[1])
	log.WarningfError(nil, testLogMessage, arguments[0], arguments[1])
	log.CriticalfError(nil, testLogMessage, arguments[0], arguments[1])

	if len(logBuffer.Buffer) > 0 {
		t.Fatalf("Expected no output. Ouput received: '%s'", string(logBuffer.Buffer))
	}
}

// TODO: test error trigger
// TODO: test that prefixs are added
