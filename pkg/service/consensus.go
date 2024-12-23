package service

import (
	"context"
	"sync"
	"time"

	"hello/pkg/core/model"
	"hello/pkg/swift"
	"hello/pkg/util"
)

type Raft struct {
	mu          sync.Mutex
	currentTerm int
	votedFor    string
	log         []*model.Block
	commitIndex int
	lastApplied int
	peers       []string
	state       string // "follower", "candidate", "leader"
	leader      string
	heartbeatCh chan struct{}
	voteCh      chan struct{}
	applyCh     chan *model.Block
	oracle      *Oracle
}

func NewRaft(peers []string, oracle *Oracle) *Raft {
	return &Raft{
		peers:       peers,
		state:       "follower",
		heartbeatCh: make(chan struct{}),
		voteCh:      make(chan struct{}),
		applyCh:     make(chan *model.Block, 10),
		oracle:      oracle,
	}
}

func (r *Raft) Start() {
	go r.run()
}

func (r *Raft) run() {
	for {
		switch r.state {
		case "follower":
			r.runFollower()
		case "candidate":
			r.runCandidate()
		case "leader":
			r.runLeader()
		}
	}
}

func (r *Raft) runFollower() {
	timer := time.NewTimer(150 * time.Millisecond)
	defer timer.Stop()

	for {
		select {
		case <-r.heartbeatCh:
			timer.Reset(150 * time.Millisecond)
		case <-timer.C:
			r.state = "candidate"
			return
		}
	}
}

func (r *Raft) runCandidate() {
	r.mu.Lock()
	r.currentTerm++
	r.votedFor = util.GetNodeID()
	r.mu.Unlock()

	timer := time.NewTimer(150 * time.Millisecond)
	defer timer.Stop()

	votes := 1 // vote for self
	for _, peer := range r.peers {
		go func(peer string) {
			if r.requestVote(peer) {
				r.mu.Lock()
				votes++
				r.mu.Unlock()
			}
		}(peer)
	}

	for {
		select {
		case <-r.heartbeatCh:
			r.state = "follower"
			return
		case <-timer.C:
			if votes > len(r.peers)/2 {
				r.state = "leader"
				return
			}
			return
		}
	}
}

func (r *Raft) runLeader() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.sendHeartbeats()
		case <-r.heartbeatCh:
			r.state = "follower"
			return
		}
	}
}

func (r *Raft) requestVote(peer string) bool {
	packet := &swift.Packet{
		Type:    swift.PacketTypeRaftRequestVote,
		Payload: nil, // Populate with necessary vote request data
	}
	err := r.oracle.swift.Send(context.Background(), packet)
	if err != nil {
		return false
	}

	// Process the response (assume success for now)
	return true
}

func (r *Raft) sendHeartbeats() {
	for _, peer := range r.peers {
		go func(peer string) {
			packet := &swift.Packet{
				Type:    swift.PacketTypeRaftHeartbeat,
				Payload: nil, // Populate with heartbeat data
			}
			err := r.oracle.swift.Send(context.Background(), packet)
			if err != nil {
				// Handle error
			}
		}(peer)
	}
}

func (r *Raft) CommitBlock(block *model.Block) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.log = append(r.log, block)
	if len(r.log)-1 > r.commitIndex {
		r.commitIndex++
		r.applyCh <- block
	}
}

func (r *Raft) ApplyBlocks() {
	for block := range r.applyCh {
		// Apply block to state machine (e.g., update chain storage)
		_ = block // Replace with actual implementation
	}
}
