package logger

import (
	log "github.com/sirupsen/logrus"
	"github.com/softwareplace/http-utils/utils"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"time"
)

var (
	logDateFormat = "2006-01-02"      // logDateFormat defines the date format (YYYY-MM-DD) used for naming log files and tracking the current date.
	logFile       *os.File            // Global variable to keep the file open
	logDirPath    string              // Global variable to store the log directory path
	logFilePath   string              // Global variable to store the full log file path
	currentDate   string              // Global variable to track the current date
	appName       string         = "" // appName holds the name of the application for logging and identification purposes.
)

func init() {
	logDirPath = "./.log"
	if envLogPath, exists := os.LookupEnv("LOG_DIR"); exists && envLogPath != "" {
		logDirPath = envLogPath
	}

	if envAppName, exists := os.LookupEnv("LOG_APP_NAME"); exists && envAppName != "" {
		appName = envAppName
	}

	logDirPath = utils.UserHomePathFix(logDirPath)
}

type Config struct {
	logDirPath string // Global variable to store the log directory path
	// LogDateFormat defines the date format (e.g., "2006-01-02") for naming log files.
	// Default value is "2006-01-02".
	LogDateFormat string

	// LumberjackLogger represents the configuration for the log rotation mechanism
	// using lumberjack.Logger.
	// Default DefaultLumberjackLogger()
	LumberjackLogger *lumberjack.Logger
}

func LumberjackLoggerSetup(config Config) {
	if config.LogDateFormat != "" {
		logDateFormat = config.LogDateFormat
	}
	if config.LumberjackLogger != nil {
		LogSetup(logDirPath, config.LumberjackLogger)
		return
	}

	if config.logDirPath != "" {
		logDirPath = utils.UserHomePathFix(config.logDirPath)
	}

	LogSetup(logDirPath, DefaultLumberjackLogger())
}

func LogSetup(dirPath string, logger *lumberjack.Logger) {
	logDirPath = dirPath
	createLogFile(logger)
	go monitorLogFile(logger)
}

// createLogFile creates or reopens the log file
func createLogFile(logger *lumberjack.Logger) {
	if err := os.MkdirAll(logDirPath, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Generate the log file name with the current date
	currentDate = time.Now().Format(logDateFormat) // yyyy-MM-dd format
	logFileName := currentDate + ".log"

	if appName != "" {
		logFileName = appName + "-" + logFileName
	}
	logFilePath = filepath.Join(logDirPath, logFileName)

	// Open or create the log file
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	// Create a multi-writer that writes to both the terminal and the file
	multiWriter := io.MultiWriter(os.Stdout, DefaultLumberjackLogger())
	// Set the output of the logger to the multi-writer
	log.SetOutput(multiWriter)
}

func DefaultLumberjackLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100,  // Max file size in MB
		MaxBackups: 3,    // Max number of old log files to retain
		MaxAge:     28,   // Max number of days to retain old log files
		Compress:   true, // Compress old log files
	}
}

// monitorLogFile periodically checks if the log file exists and recreates it if necessary
func monitorLogFile(logger *lumberjack.Logger) {
	for {
		time.Sleep(1 * time.Second) // Check every second

		// Check if the current date has changed
		newDate := time.Now().Format(logDateFormat)
		if newDate != currentDate {
			rotateLogFile(newDate, logger)
		}

		// Check if the log file still exists
		if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
			log.Println("Log file was deleted. Recreating...")
			createLogFile(logger) // Recreate the log file
		}
	}
}

// rotateLogFile closes the current log file and creates a new one with the updated date
func rotateLogFile(newDate string, logger *lumberjack.Logger) {
	// Close the current log file
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			log.Printf("Failed to close log file: %v", err)
		}
	}

	// Update the current date
	currentDate = newDate

	// Create a new log file with the updated date
	createLogFile(logger)
}
