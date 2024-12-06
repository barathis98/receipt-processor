package utils

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
)

func InitLogger() {

	if os.Getenv("ENV") == "test" {
		Logger = zap.NewNop()
		return
	}
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

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"log/app.log", "stderr"}

	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05Z07:00")

	logger, err := config.Build()
	if err != nil {
		log.Fatal("Could not initialize zap logger: ", err)
	}

	Logger = logger

	defer Logger.Sync()

	fileLogger := log.New(file, "APP_LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	fileLogger.Println("Logger initialized successfully")
}
