package my_logging

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"
)

func noErrorTestRun(t *testing.T, logLevel string) {
	message := fmt.Sprintf("Testing %s. Only %s or higher level logs should show.\n", logLevel, logLevel)
	t.Log(message)
	SetLogLevel(logLevel)
	Debugf("Testing %s", logLevel)
	Infof("Testing %s", logLevel)
	Warningf("Testing %s", logLevel)
}

func withNilErrorTestRun(t *testing.T, nilError error) {
	if nilError != nil {
		fmt.Println("TESTING ERROR: Real error passed to withNilErrorTestRun")
		os.Exit(1)
	}
	t.Log("Testing nil error. No logs should show")
	SetLogLevel("DEBUG")
	DebugfError(nilError, "Testing %s", logLevel)
	InfofError(nilError, "Testing %s", logLevel)
	WarningfError(nilError, "Testing %s", logLevel)
	CriticalfError(nilError, "Testing %s", logLevel)
}

func withRealErrorTestRun(t *testing.T, realError error, logLevel string) {
	if realError == nil {
		fmt.Println("TESTING ERROR: Nil error passed to withRealErrorTestRun")
		os.Exit(1)
	}
	message := fmt.Sprintf("Testing real error at %s level. Only %s or higher level logs should show.\n", logLevel, logLevel)
	t.Log(message)
	DebugfError(realError, "Testing %s", logLevel)
	InfofError(realError, "Testing %s", logLevel)
	WarningfError(realError, "Testing %s", logLevel)
}

func TestStandardLogging(t *testing.T) {
	fmt.Println("TestStandardLogging")
	logFile := "./test_log.txt"
	message := fmt.Sprintf("Setting log file to %s", logFile)
	t.Log(message)
	SetLogFile(logFile)

	noErrorTestRun(t, "DEBUG")
	noErrorTestRun(t, "INFO")
	noErrorTestRun(t, "WARNING")
	noErrorTestRun(t, "CRITICAL")

	//Open log file and remove it for the next test
	fileBytes, err := os.ReadFile(logFile)
	if err != nil {
		fmt.Println("TESTING ERROR: Cant open log file")
		os.Exit(1)
	}
	fileData := string(fileBytes)
	err = os.Remove(logFile)
	if err != nil {
		fmt.Println("TESTING ERROR: Couldn't delete log file after testing")
		os.Exit(1)
	}

	logRegex := `^DEBUG: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: Testing DEBUG\nINFO: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: Testing DEBUG\nWARNING: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: Testing DEBUG\nINFO: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: Testing INFO\nWARNING: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: Testing INFO\nWARNING: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: Testing WARNING`

	logMatch, err := regexp.MatchString(logRegex, fileData)
	if err != nil {
		fmt.Println("TESTING ERROR: Regex failed when comparing to log output")
		os.Exit(1)
	}
	if logMatch {
		fmt.Println()
		return
	} else {
		t.Error("Log output in log file didn't match expected output for TestStandardLogging")
	}

}

func TestNilError(t *testing.T) {
	fmt.Println("TestNilError")
	logFile := "./test_log.txt"
	message := fmt.Sprintf("Setting log file to %s", logFile)
	t.Log(message)
	SetLogFile(logFile)

	// Create nil error
	returnNil := func() error {
		return nil
	}
	nilError := returnNil()

	withNilErrorTestRun(t, nilError)

	//Open log file and remove it for the next test
	fileBytes, err := os.ReadFile(logFile)
	if err != nil {
		fmt.Println("TESTING ERROR: Cant open log file")
		os.Exit(1)
	}
	err = os.Remove(logFile)
	if err != nil {
		fmt.Println("TESTING ERROR: Couldn't delete log file after testing")
		os.Exit(1)
	}

	if len(fileBytes) == 0 {
		fmt.Println()
		return
	} else {
		t.Error("Logs file has output when none was expected for TestNilErrorTestRun")
	}
}

func TestRealError(t *testing.T) {
	fmt.Println("TestRealError")
	logFile := "./test_log.txt"
	message := fmt.Sprintf("Setting log file to %s", logFile)
	t.Log(message)
	SetLogFile(logFile)

	realError := errors.New("testing valid error. You should see me in both the console and file")
	withRealErrorTestRun(t, realError, "DEBUG")
	withRealErrorTestRun(t, realError, "INFO")
	withRealErrorTestRun(t, realError, "WARNING")

	//Open log file and remove it for the next test
	fileBytes, err := os.ReadFile(logFile)
	if err != nil {
		fmt.Println("TESTING ERROR: Cant open log file")
		os.Exit(1)
	}
	fileData := string(fileBytes)
	err = os.Remove(logFile)
	if err != nil {
		fmt.Println("TESTING ERROR: Couldn't delete log file after testing")
		os.Exit(1)
	}

	logRegex := `^DEBUG: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: Testing DEBUG\nDEBUG: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: testing valid error. You should see me in both the console and file\nINFO: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: Testing DEBUG\nINFO: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: testing valid error. You should see me in both the console and file\nWARNING: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: Testing DEBUG\nWARNING: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: testing valid error. You should see me in both the console and file\nDEBUG: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: Testing INFO\nDEBUG: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: testing valid error. You should see me in both the console and file\nINFO: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: Testing INFO\nINFO: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: testing valid error. You should see me in both the console and file\nWARNING: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: Testing INFO\nWARNING: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: testing valid error. You should see me in both the console and file\nDEBUG: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: Testing WARNING\nDEBUG: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: testing valid error. You should see me in both the console and file\nINFO: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: Testing WARNING\nINFO: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: testing valid error. You should see me in both the console and file\nWARNING: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: Testing WARNING\nWARNING: [0-9/]{10} [0-9:]{8} my_logger.go:[0-9]+: testing valid error. You should see me in both the console and file`

	logMatch, err := regexp.MatchString(logRegex, fileData)
	if err != nil {
		fmt.Println("TESTING ERROR: Regex failed when comparing to log output")
		os.Exit(1)
	}
	if logMatch {
		fmt.Println()
		return
	} else {
		t.Error("Log output in log file didn't match expected output for TestRealError")
	}
}
