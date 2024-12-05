package utils

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Define a global logger instance
	Logger *zap.Logger
)

func InitLogger() {
	// Create a log file or open an existing one
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		err := os.Mkdir("log", 0755)
		if err != nil {
			log.Fatal("Could not create log directory: ", err)
		}
	}
	file, err := os.OpenFile("log/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Could not open log file: ", err)
	}

	// Initialize zap logger (production setup with file logging)
	// Use zap.NewDevelopment() if you want a more human-readable format for development
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"log/app.log", "stderr"} // log to both file and stdout

	// Set UTC time format in zap logger configuration
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05Z07:00")

	logger, err := config.Build()
	if err != nil {
		log.Fatal("Could not initialize zap logger: ", err)
	}

	Logger = logger
	// Ensure the logger is flushed when the program exits
	defer Logger.Sync()

	// Also log to the file (for debugging purposes)
	// You can log this with a simple log package, or you can integrate this with zap too
	fileLogger := log.New(file, "APP_LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	fileLogger.Println("Logger initialized successfully")
}
