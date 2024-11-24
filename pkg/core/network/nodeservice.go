package network

import (
	"context"
	"encoding/json"
	"sync"

	. "hello/pkg/core/model"
	. "hello/pkg/core/storage"
	swift "hello/pkg/swift"
)

// NodeService manages communication between nodes
type NodeService struct {
	chainStorage  *ChainStorage
	statusStorage *StatusFile
	mempool       *MempoolStorage
	swift         *swift.Server

	mu sync.RWMutex
}

// NewNodeService creates a new NodeService instance
func NewNodeService(
	chainStorage *ChainStorage,
	statusStorage *StatusFile,
	mempool *MempoolStorage,
	swift *swift.Server,
) *NodeService {
	return &NodeService{
		chainStorage:  chainStorage,
		statusStorage: statusStorage,
		mempool:       mempool,
		swift:         swift,
	}
}

// BroadcastTransaction broadcasts a new transaction to other nodes
func (s *NodeService) BroadcastTransaction(ctx context.Context, tx *SignedTransaction) error {
	txPacket := &swift.Packet{
		Type:    swift.PacketTypeTransaction,
		Payload: json.RawMessage(tx.Ser()),
	}
	return s.swift.Broadcast(ctx, txPacket)
}
