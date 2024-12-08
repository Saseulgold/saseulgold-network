package swift

import (
	"time"
)

type MessageType string

const (
	TypeTransaction MessageType = "transaction"
	TypeRequest     MessageType = "request"
	TypeHeartbeat   MessageType = "heartbeat"
)

// Message is the basic message struct used for communication between nodes
type Message struct {
	Type      MessageType `json:"type"`
	Timestamp int64       `json:"timestamp"`
	Payload   interface{} `json:"payload"`
	From      string      `json:"from"`
	To        string      `json:"to"`
}

func NewMessage(msgType MessageType, payload interface{}, from, to string) *Message {
	return &Message{
		Type:      msgType,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
		From:      from,
		To:        to,
	}
}
