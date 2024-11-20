package main

import (
	"fmt"
	"hello/pkg/swift"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	security := swift.SecurityConfig{
		UseTLS: false,
	}

	server := swift.NewServer(":8080", security)

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	fmt.Println("Shutting down server...")
}
