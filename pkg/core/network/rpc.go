package network

import (
	"encoding/binary"
	"encoding/json"
	"hello/pkg/swift"
	"io"
	"log"
	"net"
)

type RPCRequest struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

type RPCResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func CallRPC(targetNode string, packet swift.Packet) (swift.Packet, error) {
	// Connect to the server
	conn, err := net.Dial("tcp", targetNode)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Serialize the packet
	packetData, err := json.Marshal(packet)
	if err != nil {
		log.Fatalf("Failed to serialize packet: %v", err)
	}

	// Prepare header (packet length)
	packetLen := uint32(len(packetData))
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, packetLen)

	// Send header + packet data
	if _, err := conn.Write(header); err != nil {
		log.Fatalf("Failed to send header: %v", err)
	}
	if _, err := conn.Write(packetData); err != nil {
		log.Fatalf("Failed to send packet: %v", err)
	}

	// Read response header
	respHeader := make([]byte, 4)
	if _, err := io.ReadFull(conn, respHeader); err != nil {
		log.Fatalf("Failed to read response header: %v", err)
	}

	// Get response packet length
	respPacketLen := binary.BigEndian.Uint32(respHeader)

	// Read response packet data
	respData := make([]byte, respPacketLen)
	if _, err := io.ReadFull(conn, respData); err != nil {
		log.Fatalf("Failed to read response packet: %v", err)
	}

	// Deserialize the response
	var response swift.Packet
	if err := json.Unmarshal(respData, &response); err != nil {
		log.Fatalf("Invalid response format: %v", err)
	}

	// Log and return the response
	log.Printf("RPC Response: %+v", response)
	return response, nil
}
