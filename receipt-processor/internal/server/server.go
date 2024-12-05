package server

import (
	"fmt"
	"net/http"
	"os"
	"receipt-processor/internal/routes"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
}

func NewServer() *http.Server {
	defaultPort := 8080

	portStr := os.Getenv("PORT")

	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		port = defaultPort
	}

	NewServer := &Server{
		port: port,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      routes.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
