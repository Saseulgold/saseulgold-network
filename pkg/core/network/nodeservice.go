package network

import (
	"fmt"
	"sync"

	"hello/pkg/swift"
)

type NodeService struct {
	server *swift.Server
	mu     sync.RWMutex
}

type CommandRequest struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

type CommandResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewNodeService(swiftServer *swift.Server) *NodeService {
	ns := &NodeService{
		server: swiftServer,
	}

	return ns
}

func (ns *NodeService) Start() error {
	// Start Swift server
	if err := ns.server.Start(); err != nil {
		return fmt.Errorf("failed to start swift server: %v", err)
	}

	return nil
}

func (ns *NodeService) Stop() error {
	// Stop Swift server
	if err := ns.server.Shutdown(); err != nil {
		return fmt.Errorf("failed to stop swift server: %v", err)
	}

	return nil
}
