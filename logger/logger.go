package logger

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"github.com/softwareplace/goserve/utils"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"time"
)

var (
	logFileNameDateFormat = "2006-01-02"     // logDateFormat defines the date format (YYYY-MM-DD) used for naming log files and tracking the current date.
	logFile               *os.File           // Global variable to keep the file open
	logDirPath            string             // Global variable to store the log directory path
	logFilePath           string             // Global variable to store the full log file path
	currentDate           string             // Global variable to track the current date
	appName               = ""               // appName holds the name of the application for logging and identification purposes.
	LogReportCaller       = false            // LogReportCaller determines whether the logger includes caller information (file and line number) in log messages.
	LumberjackLogger      *lumberjack.Logger // LumberjackLogger is initialized as the default rotating logger for log file management via DefaultLumberjackLogger.
	Formatter             log.Formatter      // Formatter defines the formatter used for structuring log messages.
)

func init() {
	logFileNameDateFormat = utils.GetEnvOrDefault("LOG_FILE_NAME_DATE_FORMAT", "2006-01-02")
	logDirPath = utils.GetEnvOrDefault("LOG_DIR", "./.log/")
	logDirPath = utils.UserHomePathFix(logDirPath)
	appName = utils.GetEnvOrDefault("LOG_APP_NAME", "")
	LogReportCaller = utils.GetBoolEnvOrDefault("LOG_REPORT_CALLER", false)
	LumberjackLogger = DefaultLumberjackLogger()
	Formatter = NestedLogFormatter()
}

// LogSetup initializes the logging system by creating a log file and starting a monitoring routine for file rotation.
// Default:
//   - LumberjackLogger used as default. You can also change it
func LogSetup() {
	createLogFile()
	go monitorLogFile()
}

// createLogFile creates or reopens the log file
func createLogFile() {
	if err := os.MkdirAll(logDirPath, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Generate the log file name with the current date
	currentDate = time.Now().Format(logFileNameDateFormat)
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
	LumberjackLogger.Filename = logFilePath

	multiWriter := io.MultiWriter(os.Stdout, LumberjackLogger)
	// Set the output of the logger to the multi-writer
	log.SetOutput(multiWriter)
	log.SetReportCaller(LogReportCaller)
	log.SetFormatter(Formatter)
}

func NestedLogFormatter() *nested.Formatter {
	return &nested.Formatter{
		HideKeys:        true,
		NoColors:        true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		FieldsOrder:     []string{"component", "category"},
	}
}

// DefaultLumberjackLogger returns a configured Lumberjack logger instance.
// The logger is used for managing file-based log rotation.
// Configuration:
//   - Filename: Path to the log file for writing logs.
//   - MaxSize: Maximum size of a log file in MB before rotation occurs.
//   - MaxBackups: Maximum number of old log files to retain.
//   - MaxAge: Maximum number of days old log files are retained.
//   - Compress: Specifies whether to compress old log files.
func DefaultLumberjackLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		MaxSize:    100,  // Max file size in MB
		MaxBackups: 3,    // Max number of old log files to retain
		MaxAge:     28,   // Max number of days to retain old log files
		Compress:   true, // Compress old log files
	}
}

// monitorLogFile periodically checks if the log file exists and recreates it if necessary
func monitorLogFile() {
	for {
		time.Sleep(1 * time.Second) // Check every second

		// Check if the current date has changed
		newDate := time.Now().Format(logFileNameDateFormat)
		if newDate != currentDate {
			rotateLogFile(newDate)
		}

		// Check if the log file still exists
		if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
			log.Println("Log file was deleted. Recreating...")
			createLogFile() // Recreate the log file
		}
	}
}

// rotateLogFile closes the current log file and creates a new one with the updated date
func rotateLogFile(newDate string) {
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
	createLogFile()
}
