package main

import (
	"fmt"
	"receipt-processor/internal/server"
)

func main() {

	server := server.NewServer()
	port := server.Addr
	fmt.Printf("Server is running on port %s\n", port)

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
