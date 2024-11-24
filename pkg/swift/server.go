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
	nodeID        string
	connections   map[string]net.Conn
	metrics       *Metrics
	priorityQueue *PriorityQueue
	security      SecurityConfig
	compression   CompressionType
	listener      net.Listener
	retryConfig   RetryConfig
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
		nodeID:        fmt.Sprintf("%016x", time.Now().UnixNano()),
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

		clientID := conn.RemoteAddr().String()
		s.connections[clientID] = conn

		go s.handleConnection(conn)
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

func (s *Server) handleConnection(conn net.Conn) error {
	defer conn.Close()

	for {
		// 패킷 길이를 담은 헤더(4바이트) 읽기
		header := make([]byte, 4)
		if _, err := io.ReadFull(conn, header); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("헤더 읽기 실패: %v", err)
		}

		// 패킷 길이 추출
		packetLen := binary.BigEndian.Uint32(header)

		// 패킷 데이터 읽기
		packetBytes := make([]byte, packetLen)
		if _, err := io.ReadFull(conn, packetBytes); err != nil {
			return fmt.Errorf("패킷 읽기 실패: %v", err)
		}

		// 패킷 디코딩
		var packet Packet
		if err := json.Unmarshal(packetBytes, &packet); err != nil {
			return fmt.Errorf("패킷 파싱 실패: %v", err)
		}

		// 패킷 타입에 따른 처리
		switch packet.Type {
		case PacketTypeTransaction:
			// 트랜잭션 처리 로직
		case PacketTypeBlock:
			// 블록 처리 로직
		case PacketTypeHeightRequest:
			// 블록 높이 요청 처리
		case PacketTypeHeightResponse:
			// 블록 높이 응답 처리
		case PacketTypeBlockRequest:
			// 블록 요청 처리
		case PacketTypeBlockResponse:
			// 블록 응답 처리
		default:
			return fmt.Errorf("알 수 없는 패킷 타입: %v", packet.Type)
		}
	}
}

func (s *Server) closeConnection(conn net.Conn) error {
	return conn.Close()
}

// BroadcastBlock은 블록을 모든 피어에게 전송합니다
func (s *Server) Broadcast(ctx context.Context, packet *Packet) error {
	return s.Send(packet.Type, packet.Payload)
}

func (s *Server) Send(packetType PacketType, payload interface{}) error {
	conn := s.connections[s.nodeID]
	packet := &Packet{
		Type: packetType,
	}

	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to serialize payload: %v", err)
		}
		packet.Payload = payloadBytes
	}

	packetBytes, err := json.Marshal(packet)
	if err != nil {
		return fmt.Errorf("failed to serialize packet: %v", err)
	}

	// Create header with packet length
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(packetBytes)))

	// Send header and packet sequentially
	for _, conn := range s.connections {
		if _, err := conn.Write(header); err != nil {
			return fmt.Errorf("헤더 전송 실패: %v", err)
		}
	}

	if _, err := conn.Write(packetBytes); err != nil {
		return fmt.Errorf("failed to send packet: %v", err)
	}

	return nil
}

func (s *Server) ReceiveMessage() (*Packet, error) {
	conn := s.connections[s.nodeID]

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
	conn := s.connections[s.nodeID]

	if conn == nil {
		return nil, errors.New("server is not connected")
	}

	// Send peer list request message
	heightReq := &Packet{
		Type:    PacketTypeHeightRequest,
		Payload: json.RawMessage(s.nodeID),
	}
	if err := s.Send(PacketTypeHeightRequest, heightReq); err != nil {
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
