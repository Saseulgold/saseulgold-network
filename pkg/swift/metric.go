package swift

import (
	"sync/atomic"
	"time"
)

type Metrics struct {
	totalMessages     uint64
	failedMessages    uint64
	avgProcessingTime int64
	activeConnections int32
	lastHeartbeat     time.Time
	bytesReceived     uint64
	bytesSent         uint64
}

func NewMetrics() *Metrics {
	return &Metrics{
		lastHeartbeat: time.Now(),
	}
}

func (m *Metrics) IncrementMessages() {
	atomic.AddUint64(&m.totalMessages, 1)
}

func (m *Metrics) IncrementFailedMessages() {
	atomic.AddUint64(&m.failedMessages, 1)
}

func (m *Metrics) AddProcessingTime(duration time.Duration) {
	atomic.AddInt64(&m.avgProcessingTime, int64(duration))
}

func (m *Metrics) UpdateActiveConnections(delta int32) {
	atomic.AddInt32(&m.activeConnections, delta)
}

func (m *Metrics) AddBytesReceived(bytes uint64) {
	atomic.AddUint64(&m.bytesReceived, bytes)
}

func (m *Metrics) AddBytesSent(bytes uint64) {
	atomic.AddUint64(&m.bytesSent, bytes)
}

func (m *Metrics) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"total_messages":      atomic.LoadUint64(&m.totalMessages),
		"failed_messages":     atomic.LoadUint64(&m.failedMessages),
		"avg_processing_time": atomic.LoadInt64(&m.avgProcessingTime),
		"active_connections":  atomic.LoadInt32(&m.activeConnections),
		"bytes_received":      atomic.LoadUint64(&m.bytesReceived),
		"bytes_sent":          atomic.LoadUint64(&m.bytesSent),
		"last_heartbeat":      m.lastHeartbeat,
	}
}
