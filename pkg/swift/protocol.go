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

	PacketTypeReplicateBlockRequest  // Block transmission
	PacketTypeReplicateBlockResponse // Block transmission

	PacketTypeErrorResponse // Error

	// Consensus
	PacketTypeRaftRequestVote  // Raft request vote
	PacketTypeRaftResponseVote // Raft response vote
	PacketTypeRaftHeartbeat    // Raft heartbeat

	PacketTypeRegisterReplicaRequest  // Register replica
	PacketTypeRegisterReplicaResponse // Register replica

	PacketTypeMetricsRequest  // Metrics request
	PacketTypeMetricsResponse // Metrics response

	PacketTypeSearchRequest  // Search request
	PacketTypeSearchResponse // Search response

	PacketTypeLastHeightRequest
	PacketTypeLastHeightResponse

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
