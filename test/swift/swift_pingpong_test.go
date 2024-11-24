package main

import (
	"testing"

	. "hello/pkg/swift"
)

func TestServerPingPong(t *testing.T) {
	// Configure two servers
	t.Log("Configuring servers...")
	server1 := NewServer("localhost:9001", SecurityConfig{UseTLS: false})
	server2 := NewServer("localhost:9002", SecurityConfig{UseTLS: false})

	// Start servers
	t.Log("Starting server 1...")
	if err := server1.Start(); err != nil {
		t.Fatalf("Failed to start server 1: %v", err)
	}
	t.Log("Starting server 2...")
	if err := server2.Start(); err != nil {
		t.Fatalf("Failed to start server 2: %v", err)
	}

	// Cleanup after test
	defer func() {
		t.Log("Shutting down servers...")
		server1.Shutdown()
		server2.Shutdown()
	}()

	if err := server1.Connect("localhost:9002"); err != nil {
		t.Fatalf("Failed to connect servers: %v", err)
	}

	t.Log("Ping-pong test completed successfully")
}
