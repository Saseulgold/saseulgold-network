package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"hello/pkg/swift"
)

type NodeService struct {
	server     *swift.Server
	httpServer *http.Server
	mu         sync.RWMutex
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

func NewNodeService(swiftServer *swift.Server, rpcPort string) *NodeService {
	ns := &NodeService{
		server: swiftServer,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/node/command", ns.handleCommand)

	ns.httpServer = &http.Server{
		Addr:    rpcPort,
		Handler: mux,
	}

	return ns
}

func (ns *NodeService) Start() error {
	// Start Swift server
	if err := ns.server.Start(); err != nil {
		return fmt.Errorf("failed to start swift server: %v", err)
	}

	// Start RPC server
	go func() {
		if err := ns.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("RPC server error: %v\n", err)
		}
	}()

	return nil
}

func (ns *NodeService) Stop() error {
	// Stop RPC server
	if err := ns.httpServer.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("failed to stop RPC server: %v", err)
	}

	// Stop Swift server
	if err := ns.server.Shutdown(); err != nil {
		return fmt.Errorf("failed to stop swift server: %v", err)
	}

	return nil
}

func (ns *NodeService) handleCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "method not supported", http.StatusMethodNotAllowed)
		return
	}

	var cmd CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		sendErrorResponse(w, "invalid request format", http.StatusBadRequest)
		return
	}

	ns.mu.Lock()
	defer ns.mu.Unlock()

	var response CommandResponse

	switch cmd.Action {
	case "connect":
		targetAddr, ok := cmd.Payload.(string)
		if !ok {
			sendErrorResponse(w, "invalid payload format", http.StatusBadRequest)
			return
		}
		err := ns.server.Connect(targetAddr)
		response = CommandResponse{
			Success: err == nil,
			Error:   getErrorString(err),
		}

	case "disconnect":
		targetAddr, ok := cmd.Payload.(string)
		if !ok {
			sendErrorResponse(w, "invalid payload format", http.StatusBadRequest)
			return
		}
		err := ns.server.Close(targetAddr)
		response = CommandResponse{
			Success: err == nil,
			Error:   getErrorString(err),
		}

	case "ping":
		err := ns.server.Ping(context.Background())
		response = CommandResponse{
			Success: err == nil,
			Error:   getErrorString(err),
		}

	case "get_peers":
		peers, err := ns.server.GetPeers()
		response = CommandResponse{
			Success: err == nil,
			Error:   getErrorString(err),
			Data:    peers,
		}

	default:
		sendErrorResponse(w, "unsupported command", http.StatusBadRequest)
		return
	}

	sendJSONResponse(w, response)
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := CommandResponse{
		Success: false,
		Error:   message,
	}
	w.WriteHeader(statusCode)
	sendJSONResponse(w, response)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func getErrorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
