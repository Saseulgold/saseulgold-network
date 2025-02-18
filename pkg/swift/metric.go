package swift

import (
	"sync/atomic"
	"time"

	"hello/pkg/util"
)

type Metrics struct {
	totalMessages     uint64
	failedMessages    uint64
	avgProcessingTime uint64
	activeConnections uint32
	lastHeartbeat     time.Time
	bytesReceived     uint64
	bytesSent         uint64

	broadcastCount   uint64
	broadcastBytes   uint64
	broadcastTime    uint64
	broadcastAvgTime uint64
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
	atomic.AddUint64(&m.avgProcessingTime, uint64(duration))
}

func (m *Metrics) UpdateActiveConnections(delta uint32) {
	atomic.AddUint32(&m.activeConnections, delta)
}

func (m *Metrics) AddBytesReceived(bytes uint64) {
	atomic.AddUint64(&m.bytesReceived, bytes)
}

func (m *Metrics) AddBytesSent(bytes uint64) {
	atomic.AddUint64(&m.bytesSent, bytes)
}

func (m *Metrics) BroadcastStart() int64 {
	atomic.AddUint64(&m.broadcastCount, 1)
	return util.Utime()
}

func (m *Metrics) BroadcastEnd(start int64, packetBytes []byte) {
	duration := util.Utime() - start
	atomic.AddUint64(&m.broadcastBytes, uint64(len(packetBytes)))
	atomic.AddUint64(&m.broadcastTime, uint64(duration))

	avgTime := uint64(m.broadcastTime) / atomic.LoadUint64(&m.broadcastCount)
	atomic.StoreUint64(&m.broadcastAvgTime, avgTime)
}

func (m *Metrics) GetStats() map[string]interface{} {
	return map[string]interface{}{
		// "total_messages":      atomic.LoadUint64(&m.totalMessages),
		//"failed_messages":     atomic.LoadUint64(&m.failedMessages),
		"bytes_received": atomic.LoadUint64(&m.bytesReceived),
		"bytes_sent":     atomic.LoadUint64(&m.bytesSent),

		"broadcast_count":    atomic.LoadUint64(&m.broadcastCount),
		"broadcast_bytes":    atomic.LoadUint64(&m.broadcastBytes),
		"broadcast_time":     atomic.LoadUint64(&m.broadcastTime),
		"broadcast_avg_time": atomic.LoadUint64(&m.broadcastAvgTime),

		// "last_heartbeat": m.lastHeartbeat,
	}
}
