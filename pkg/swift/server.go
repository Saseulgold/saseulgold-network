package swift

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	C "hello/pkg/core/config"
)

type RetryConfig struct {
	RetryDelay time.Duration
	MaxRetries int
	MaxBackoff time.Duration
}

type Server struct {
	address       string
	connections   map[string]net.Conn
	metrics       *Metrics
	priorityQueue *PriorityQueue
	security      SecurityConfig
	compression   CompressionType
	listener      net.Listener
	retryConfig   RetryConfig
	mu            sync.RWMutex
}

func SwiftInfoLog(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func SwiftRootDir() string {
	if C.CORE_TEST_MODE {
		return filepath.Join(C.QUANTUM_ROOT_DIR, "swift.pid")
	}
	return filepath.Join(C.QUANTUM_ROOT_DIR, "swift.pid")
}

func NewServer(address string, security SecurityConfig) *Server {
	return &Server{
		address:       address,
		connections:   make(map[string]net.Conn),
		metrics:       NewMetrics(),
		priorityQueue: &PriorityQueue{},
		security:      security,
		compression:   GzipCompression,
		retryConfig: RetryConfig{
			RetryDelay: time.Second,
			MaxRetries: 3,
			MaxBackoff: 30 * time.Second,
		},
		mu: sync.RWMutex{},
	}
}

func (s *Server) Start() error {
	var listener net.Listener
	var err error

	if s.security.UseTLS {
		tlsConfig, err := newTLSConfig(s.security)
		if err != nil {
			return err
		}
		listener, err = tls.Listen("tcp", s.address, tlsConfig)
	} else {
		listener, err = net.Listen("tcp", s.address)
	}

	if err != nil {
		return err
	}

	if err = s.writePID(); err != nil {
		return err
	}

	s.listener = listener
	// PID 파일 생성

	go s.acceptConnections()
	go s.processJobQueue()
	go s.collectMetrics()

	return nil
}

func (s *Server) writePID() error {
	pidFile := SwiftRootDir()

	// PID 파일 디렉토리 생성
	if err := os.MkdirAll(filepath.Dir(pidFile), 0755); err != nil {
		return fmt.Errorf("PID 파일 디렉토리 생성 실패: %v", err)
	}

	// 현재 프로세스의 PID를 파일에 기록
	pid := os.Getpid()
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		return fmt.Errorf("PID 파일 작성 실패: %v", err)
	}

	return nil
}

func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}

		s.mu.Lock()
		remoteAddr := conn.RemoteAddr().String()
		s.connections[remoteAddr] = conn
		s.mu.Unlock()

		go s.HandleConnection(conn)
	}
}

func (s *Server) processJobQueue() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if s.priorityQueue.Len() > 0 {
				if job := s.priorityQueue.Pop(); job != nil {
					// job 처리 로직
				}
			}
		}
	}
}

func (s *Server) collectMetrics() {
	for {
		// s.metrics.Update()
		time.Sleep(time.Second)
	}
}

func (s *Server) HandleConnection(conn net.Conn) error {
	for {
		SwiftInfoLog("new connection request: %s\n", conn.RemoteAddr().String())
		header := make([]byte, 4)
		if _, err := io.ReadFull(conn, header); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("failed to read header: %v", err)
		}
		SwiftInfoLog("received header: %v\n", header)

		// get packet length
		packetLen := binary.BigEndian.Uint32(header)

		// read packet data
		packetBytes := make([]byte, packetLen)
		if _, err := io.ReadFull(conn, packetBytes); err != nil {
			SwiftInfoLog("failed to read packet: %v\n", err)
			return fmt.Errorf("failed to read packet: %v", err)
		}

		// decode packet
		var packet Packet
		if err := json.Unmarshal(packetBytes, &packet); err != nil {
			return fmt.Errorf("failed to parse packet: %v", err)
		}
		SwiftInfoLog("received packet: %v\n", packet)

		// process packet
		switch packet.Type {
		case PacketTypeTransaction:
			// process transaction
		case PacketTypeBlock:
			// process block
		case PacketTypeHeightRequest:
			// process height request
		case PacketTypeHeightResponse:
			// process height response
		case PacketTypeBlockRequest:
			// process block request
		case PacketTypeBlockResponse:
			// process block response
		case PacketTypePing:
			if err := s.Send(&Packet{
				Type:    PacketTypePong,
				Payload: json.RawMessage(`"pong"`),
			}); err != nil {
				return fmt.Errorf("failed to send pong response: %v", err)
			}
			continue
		case PacketTypePong:
			SwiftInfoLog("pong received\n")
			continue
		default:
			return fmt.Errorf("unknown packet type: %v", packet.Type)
		}
	}
}
func (s *Server) Close(addr string) error {
	conn, exists := s.connections[addr]
	if !exists || conn == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.CloseConnection(conn)
}

func (s *Server) CloseConnection(conn net.Conn) error {
	return conn.Close()
}

// BroadcastBlock은 블록을 모든 피어에게 전송합니다
func (s *Server) Broadcast(ctx context.Context, packet *Packet) error {
	return s.Send(packet)
}

func (s *Server) Send(packet *Packet) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.connections) == 0 {
		return fmt.Errorf("no active connections")
	}

	packetBytes, err := json.Marshal(packet)
	if err != nil {
		return fmt.Errorf("failed to serialize packet: %v", err)
	}

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(packetBytes)))

	for _, conn := range s.connections {
		if conn == nil {
			continue
		}
		if _, err := conn.Write(header); err != nil {
			return fmt.Errorf("failed to send header: %v", err)
		}
		if _, err := conn.Write(packetBytes); err != nil {
			return fmt.Errorf("failed to send packet: %v", err)
		}
	}

	return nil
}

func (s *Server) ReceiveMessage() (*Packet, error) {
	s.mu.RLock()
	if len(s.connections) == 0 {
		s.mu.RUnlock()
		return nil, errors.New("서버에 활성화된 연결이 없습니다")
	}

	// 첫 번째 활성 연결 사용
	var conn net.Conn
	for _, c := range s.connections {
		conn = c
		break
	}
	s.mu.RUnlock()

	// 4바이트 헤더 버퍼 선언
	header := make([]byte, 4)

	// Read 4-byte header (packet length)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, fmt.Errorf("failed to receive header: %v", err)
	}

	// Extract packet length
	packetLen := binary.BigEndian.Uint32(header)

	// Read packet
	packetBytes := make([]byte, packetLen)
	if _, err := io.ReadFull(conn, packetBytes); err != nil {
		return nil, fmt.Errorf("failed to receive packet: %v", err)
	}

	var packet Packet
	if err := json.Unmarshal(packetBytes, &packet); err != nil {
		return nil, fmt.Errorf("failed to parse packet: %v", err)
	}

	return &packet, nil
}

func (s *Server) GetPeers() ([]string, error) {
	s.mu.RLock()
	if len(s.connections) == 0 {
		s.mu.RUnlock()
		return nil, errors.New("서버에 활성화된 연결이 없습니다")
	}

	s.mu.RUnlock()

	// Send peer list request message
	heightReq := &Packet{
		Type:    PacketTypeHeightRequest,
		Payload: json.RawMessage(s.address),
	}
	if err := s.Send(heightReq); err != nil {
		return nil, fmt.Errorf("failed to request peer list: %v", err)
	}

	// Receive response
	response, err := s.ReceiveMessage()
	if err != nil {
		return nil, fmt.Errorf("failed to receive peer list: %v", err)
	}

	// Convert response to string slice
	var peers []string
	if err := json.Unmarshal(response.Payload, &peers); err != nil {
		return nil, fmt.Errorf("failed to parse peer list: %v", err)
	}

	return peers, nil
}

// Ping sends a ping message to check connection status
func (s *Server) Ping(ctx context.Context) error {
	// "ping" 문자열을 JSON 형식으로 변경
	pingPayload := []byte(`"ping"`)
	pingPacket := &Packet{
		Type:    PacketTypePing,
		Payload: json.RawMessage(pingPayload),
	}

	if err := s.Send(pingPacket); err != nil {
		return fmt.Errorf("failed to send ping: %v", err)
	}

	// 응답 대기
	response, err := s.ReceiveMessage()
	if err != nil {
		return fmt.Errorf("failed to receive ping response: %v", err)
	}

	if response.Type != PacketTypePong {
		return fmt.Errorf("invalid response type: %v", response.Type)
	}

	return nil
}

// Shutdown gracefully shuts down the server and cleans up resources
func (s *Server) Shutdown() error {
	// 리스너 종료
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return fmt.Errorf("failed to close listener: %v", err)
		}
	}

	// 모든 연결 종료
	for id, conn := range s.connections {
		if err := conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection (%s): %v", id, err)
		}
		delete(s.connections, id)
	}

	// PID 파일 제거
	pidFile := SwiftRootDir()
	if err := os.Remove(pidFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("PID 파일 제거 실패: %v", err)
	}

	return nil
}

func (s *Server) Connect(targetAddr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		return fmt.Errorf("연결 실패: %v", err)
	}

	SwiftInfoLog("success to connect: %s\n", targetAddr)
	s.connections[targetAddr] = conn
	go s.HandleConnection(conn)

	return nil
}
