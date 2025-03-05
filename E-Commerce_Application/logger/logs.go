package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Global Logger instance
var Log *logrus.Logger

func InitLogger() {
	// Open log file
	logFile, err := os.OpenFile("C:\\Users\\mahes\\ECSA\\E-Commerce-Scalable-Platform\\E-Commerce_Application\\logs\\app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}

	// Create new logger
	Log = logrus.New()
	Log.SetOutput(logFile)                    // Write logs to file
	Log.SetFormatter(&logrus.JSONFormatter{}) // Use JSON format
	Log.SetLevel(logrus.InfoLevel)            // Set log level
}
