package swift

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	C "hello/pkg/core/config"
)

type Server struct {
	address       string
	connections   map[string]net.Conn
	metrics       *Metrics
	priorityQueue *PriorityQueue
	security      SecurityConfig
	compression   CompressionType
	listener      net.Listener
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

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

}
