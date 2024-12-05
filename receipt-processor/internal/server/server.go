package server

import (
	"fmt"
	"net/http"
	"os"
	"receipt-processor/internal/routes"
	"receipt-processor/internal/store"
	"receipt-processor/internal/utils"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

type Server struct {
	port int
}

func NewServer() *http.Server {
	defaultPort := 8080

	portStr := os.Getenv("PORT")

	utils.InitLogger()

	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		port = defaultPort
		utils.Logger.Warn("Invalid or missing PORT environment variable, using default port", zap.Int("port", port))
	} else {
		utils.Logger.Info("Port configured from environment variable", zap.Int("port", port))
	}

	NewServer := &Server{
		port: port,
	}

	store.InitializeStores()
	utils.Logger.Info("Stores initialized")

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      routes.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	utils.Logger.Info("Server started", zap.String("address", server.Addr))

	return server
}
