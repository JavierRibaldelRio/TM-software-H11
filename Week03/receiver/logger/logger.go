package logger

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	logDir = "logs/logs"
)

var (
	logFile *os.File
	logger  *log.Logger
)

// Create the configuration of the file
func init() {
	os.MkdirAll(logDir, 0755)

	filename := time.Now().Format("log-2006-01-02_15-04-05.txt")
	fullPath := filepath.Join(logDir, filename)

	var err error
	logFile, err = os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("The file could not be created: %v", err)
	}
	logger = log.New(logFile, "", log.LstdFlags)
}

// Writes info to the file
func Info(msg string) {
	WriteToLog("[INFO] " + msg)
}

// Writes an error to the file
func Error(msg string) {
	WriteToLog("[ERROR] " + msg)
}

// WriteToLog writes a message to the log including the date and time.
func WriteToLog(msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logger.Printf("[%s] %s", timestamp, msg)
}

// Action writes a message to the log with a custom [ACTION] tag.
func Action(msg string) {
	WriteToLog("[ACTION] " + msg)
}

// Not used
func Close() {
	logFile.Close()
}
