package swift

import (
	"encoding/json"
)

// PacketType defines the type of communication packet
type PacketType uint8

const (
	PacketTypeUnknown                        PacketType = iota
	PacketTypeBroadcastTransactionRequest               // Transaction transmission
	PacketTypeBroadcastTransactionResponse              // Transaction transmission
	PacketTypeSendTransactionRequest                    // Transaction transmission
	PacketTypeSendTransactionResponse                   // Transaction transmission
	PacketTypeListMempoolTransactionRequest             // Transaction transmission
	PacketTypeListMempoolTransactionResponse            // Transaction transmission
	PacketTypeSync                                      // Sync request
	PacketTypePing                                      // Ping
	PacketTypePong                                      // Ping
	PacketTypePeerRequest                               // Peer
	PacketTypePeerResponse                              // Peer
	PacketTypeHandshakeCMDRequest                       // Handshake
	PacketTypeHandshakeCMDResponse                      // Handshake

	PacketTypeBroadcastBlockRequest  // Block transmission
	PacketTypeBroadcastBlockResponse // Block transmission
	PacketTypeRawRequest             // Raw request
	PacketTypeRawResponse            // Raw response

	PacketTypeErrorResponse // Error

	// FOR DEV
	PacketTypeGetStatusBundleRequest  // Get status bundle
	PacketTypeGetStatusBundleResponse // Get status bundle

)

// Packet is the base packet structure
type Packet struct {
	Type    PacketType      `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// CallRawRequest sends a raw RPC request to the node
