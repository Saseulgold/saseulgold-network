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
	. "hello/pkg/core/debug"
	"hello/pkg/core/vm"
)

type PacketHandler func(ctx context.Context, packet *Packet) error

type RetryConfig struct {
	RetryDelay time.Duration
	MaxRetries int
	MaxBackoff time.Duration
}

type Server struct {
	address       string
	connections   map[string]net.Conn
	handlers      map[PacketType]func(ctx context.Context, packet *Packet) error
	metrics       *Metrics
	priorityQueue *PriorityQueue
	security      SecurityConfig
	compression   CompressionType
	listener      net.Listener
	retryConfig   RetryConfig
	mu            sync.RWMutex
}

func (s *Server) RegisterHandler(packetType PacketType, handler PacketHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[packetType] = handler
}

func (s *Server) FormatPeerList() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var connections []string
	for addr := range s.connections {
		connections = append(connections, addr)
	}

	return fmt.Sprintf("{connections: %v}", connections)
}

func SwiftInfoLog(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
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
		mu:       sync.RWMutex{},
		handlers: make(map[PacketType]func(ctx context.Context, packet *Packet) error),
	}
}

func (s *Server) Start() error {
	var listener net.Listener
	var err error

	machine := vm.GetMachineInstance()
	machine.GetInterpreter().Reset()

	if s.security.UseTLS {
		tlsConfig, err := newTLSConfig(s.security)
		if err != nil {
			return err
		}
		listener, err = tls.Listen("tcp", s.address, tlsConfig)
		DebugLog("tls listener created")
	} else {
		listener, err = net.Listen("tcp", s.address)
		DebugLog("tcp listener created")
	}

	if err != nil {
		return err
	}

	if err = s.writePID(); err != nil {
		return err
	}

	s.listener = listener

	go s.acceptConnections()
	go s.processJobQueue()
	go s.collectMetrics()

	DebugLog("server started")
	return nil
}

func (s *Server) writePID() error {
	pidFile := SwiftRootDir()

	// Create PID file directory
	if err := os.MkdirAll(filepath.Dir(pidFile), 0755); err != nil {
		return fmt.Errorf("PID file directory creation failed: %v", err)
	}

	// Write the current process's PID to the file
	pid := os.Getpid()
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		return fmt.Errorf("failed to write PID file: %v", err)
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
					// job processing logic
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
	// connection cleanup
	defer func() {
		s.mu.Lock()
		delete(s.connections, conn.RemoteAddr().String())
		s.mu.Unlock()
		conn.Close()
	}()

	// read timeout;
	if err := conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		return fmt.Errorf("failed to set read deadline: %v", err)
	}

	// context에 연결 정보 추가
	ctx := context.WithValue(context.Background(), "connection", conn)

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
		case PacketTypePing:
			if err := s.Send(ctx, &Packet{
				Type:    PacketTypePong,
				Payload: json.RawMessage(`"pong"`),
			}); err != nil {
				return fmt.Errorf("failed to send pong response: %v", err)
			}
			continue
		case PacketTypePeerRequest:
			peerList := s.FormatPeerList()
			SwiftInfoLog("peer list: %s\n", peerList)

			// server peer list to client
			peerListJSON, err := json.Marshal(peerList)
			if err != nil {
				return fmt.Errorf("peer list 직렬화 실패: %v", err)
			}

			if err := s.Send(ctx, &Packet{
				Type:    PacketTypePeerResponse,
				Payload: peerListJSON,
			}); err != nil {
				return fmt.Errorf("failed to send peer response: %v", err)
			}
			continue
		default:
			if handler, ok := s.handlers[packet.Type]; ok {
				err := handler(ctx, &packet)
				if err != nil {
					SwiftInfoLog("packet handler error: %v", err)
					return fmt.Errorf("packet handler error: %v", err)
				}
			} else {
				return fmt.Errorf("unknown packet type: %v", packet.Type)
			}
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
	s.mu.RLock()
	connections := make([]net.Conn, 0, len(s.connections))
	for _, conn := range s.connections {
		connections = append(connections, conn)
	}
	s.mu.RUnlock()

	if len(connections) == 0 {
		return fmt.Errorf("no active connections")
	}

	for _, conn := range connections {
		connCtx := context.WithValue(ctx, "connection", conn)
		if err := s.Send(connCtx, packet); err != nil {
			return fmt.Errorf("broadcast to peer(%s) failed: %v", conn.RemoteAddr().String(), err)
		}
	}

	return nil
}

func (s *Server) Send(ctx context.Context, packet *Packet) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Extract connection info from context
	connInfo, ok := ctx.Value("connection").(net.Conn)
	if !ok || connInfo == nil {
		return fmt.Errorf("connection info not found in context")
	}

	packetBytes, err := json.Marshal(packet)
	if err != nil {
		return fmt.Errorf("failed to serialize packet: %v", err)
	}

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(packetBytes)))

	if _, err := connInfo.Write(header); err != nil {
		return fmt.Errorf("failed to send header: %v", err)
	}
	if _, err := connInfo.Write(packetBytes); err != nil {
		return fmt.Errorf("failed to send packet: %v", err)
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
	peers := make([]string, 0, len(s.connections))
	for addr := range s.connections {
		peers = append(peers, addr)
	}
	s.mu.RUnlock()

	if len(peers) == 0 {
		return nil, fmt.Errorf("no active connections")
	}

	// 응답이 온 피어들을 저장할 슬라이스
	activePeers := make([]string, 0)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 각 피어에 대해 ping 전송
	for _, peer := range peers {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()

			// ping 패킷 전송
			pingPacket := &Packet{
				Type:    PacketTypePing,
				Payload: json.RawMessage(`"ping"`),
			}

			if err := s.Send(context.Background(), pingPacket); err == nil {
				mu.Lock()
				activePeers = append(activePeers, addr)
				mu.Unlock()
			}
		}(peer)
	}

	wg.Wait()
	return activePeers, nil
}

// Shutdown gracefully shuts down the server and cleans up resources
func (s *Server) Shutdown() error {
	// 스너 종료
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
