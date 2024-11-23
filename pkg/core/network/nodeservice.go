package network

import (
	"context"
	"sync"
	"time"

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
	return s.swift.BroadcastTransaction(ctx, tx)
}

// SyncBlocks synchronizes blocks with other nodes
func (s *NodeService) SyncBlocks(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentHeight := s.statusStorage.BundleHeight()

	// Check latest block heights from other nodes
	peerHeights, err := s.swift.GetPeerHeights(ctx)
	if err != nil {
		return err
	}

	// Sync blocks from peers with higher blocks
	for peer, height := range peerHeights {
		if height > currentHeight {
			blocks, err := s.swift.FetchBlocks(ctx, peer, currentHeight+1, height)
			if err != nil {
				continue
			}

			// Validate and store received blocks
			for _, block := range blocks {
				if err := s.chainStorage.AddBlock(block); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// BroadcastBlock broadcasts a newly created block to other nodes
func (s *NodeService) BroadcastBlock(ctx context.Context, block *Block) error {
	return s.swift.BroadcastBlock(ctx, block)
}

// Start starts the NodeService
func (s *NodeService) Start(ctx context.Context) error {
	// Periodically perform block synchronization
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := s.SyncBlocks(ctx); err != nil {
					// Error logging
				}
			}
		}
	}()

	return nil
}
