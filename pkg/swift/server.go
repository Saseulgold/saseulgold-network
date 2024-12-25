package swift

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	C "hello/pkg/core/config"
	. "hello/pkg/core/debug"
	"hello/pkg/core/storage"
	"hello/pkg/core/vm"
	"hello/pkg/util"
)

type PacketHandler func(ctx context.Context, packet *Packet) error

type RetryConfig struct {
	RetryDelay time.Duration
	MaxRetries int
	MaxBackoff time.Duration
}

type Server struct {
	address       string
	peers         map[string]net.Conn
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

	var peers []string
	for addr := range s.peers {
		peers = append(peers, addr)
	}

	return fmt.Sprintf("{peers: %v}", peers)
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
		peers:         make(map[string]net.Conn),
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

	isRunning := util.ServiceIsRunning(storage.DataRootDir(), "swift")

	if isRunning {
		return fmt.Errorf("swift is already running")
	}

	s.listener = listener

	go s.acceptConnections()
	go s.processJobQueue()
	go s.collectMetrics()

	err = util.ProcessStart(storage.DataRootDir(), "swift", os.Getpid())
	if err != nil {
		return err
	}

	DebugLog("server started")
	return nil
}

func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}

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
					return s.SendErrorResponse(ctx, err.Error())
				}
			} else {
				return fmt.Errorf("unknown packet type: %v", packet.Type)
			}
		}
	}
}

func (s *Server) Close(peer string) error {
	conn, exists := s.peers[peer]
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
	peers := make([]net.Conn, 0, len(s.peers))
	for _, conn := range s.peers {
		peers = append(peers, conn)
	}
	s.mu.RUnlock()

	for _, conn := range peers {
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

func (s *Server) GetPeers() []string {
	s.mu.RLock()
	peers := make([]string, 0, len(s.peers))
	for addr := range s.peers {
		peers = append(peers, addr)
	}
	s.mu.RUnlock()

	return peers
}

// Shutdown gracefully shuts down the server and cleans up resources
func (s *Server) Shutdown() error {
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return fmt.Errorf("failed to close listener: %v", err)
		}
	}

	for id, conn := range s.peers {
		if err := conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection (%s): %v", id, err)
		}
		delete(s.peers, id)
	}

	err := util.TerminateProcess(storage.DataRootDir(), "swift")
	if err != nil {
		return fmt.Errorf("failed to terminate process: %v", err)
	}

	return nil
}

func (s *Server) Connect(targetAddr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// check if already connected
	if _, exists := s.peers[targetAddr]; exists {
		return nil
	}

	// connect to target address
	conn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}

	SwiftInfoLog("peer connection success: %s\n", targetAddr)
	s.peers[targetAddr] = conn
	go s.HandleConnection(conn)

	return nil
}

// SendErrorResponse is a helper function to send an error response to the client
func (s *Server) SendErrorResponse(ctx context.Context, errMsg string) error {
	errorPayload := struct {
		Error string `json:"error"`
	}{
		Error: errMsg,
	}

	payload, err := json.Marshal(errorPayload)
	if err != nil {
		return fmt.Errorf("error message serialization failed: %v", err)
	}

	return s.Send(ctx, &Packet{
		Type:    PacketTypeErrorResponse,
		Payload: payload,
	})
}
