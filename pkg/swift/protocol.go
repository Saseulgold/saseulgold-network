package swift

import (
	"encoding/json"
	"fmt"
)

// PacketType defines the type of communication packet
type PacketType uint8

const (
	PacketTypeUnknown        PacketType = iota
	PacketTypeTransaction               // Transaction transmission
	PacketTypeBlock                     // Block transmission
	PacketTypeHeightRequest             // Block height request
	PacketTypeHeightResponse            // Block height response
	PacketTypeBlockRequest              // Block request
	PacketTypeBlockResponse             // Block response
	PacketTypeSync                      // Sync request
	PacketTypePing                      // Ping
	PacketTypePong                      // Pong
)

// Packet is the base packet structure
type Packet struct {
	Type    PacketType      `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Validate validates the packet
func (p *Packet) Validate() error {
	// 패킷 타입 검증
	if p.Type == PacketTypeUnknown {
		return fmt.Errorf("unknown packet type")
	}

	// 페이로드 검증
	if p.Type != PacketTypeHeightRequest && len(p.Payload) == 0 {
		return fmt.Errorf("payload is empty")
	}

	// 페이로드 JSON 유효성 검사
	if len(p.Payload) > 0 {
		var js json.RawMessage
		if err := json.Unmarshal(p.Payload, &js); err != nil {
			return fmt.Errorf("invalid JSON payload: %v", err)
		}
	}
	return nil
}
