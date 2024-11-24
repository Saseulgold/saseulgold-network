package main

import (
	. "hello/pkg/core/network"
	. "hello/pkg/swift"
	"testing"
)

func TestPingPong(t *testing.T) {
	server := NewServer("localhost:9001", SecurityConfig{UseTLS: false})

	node := NewNodeService(server)
	if err := node.Start(); err != nil {
		t.Fatalf("Failed to start node service: %v", err)
	}
}
