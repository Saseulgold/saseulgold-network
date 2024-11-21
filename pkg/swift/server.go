package swift

import (
	"crypto/tls"
	"net"
	"time"
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
	s.listener = listener

	go s.acceptConnections()
	// go s.processJobQueue()
	// go s.collectMetrics()

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
	for {
		if job := s.priorityQueue.Pop(); job != nil {
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
