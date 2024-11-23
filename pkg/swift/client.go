package swift

import (
	"crypto/tls"
	"errors"
	"net"
	"time"
)

type RetryConfig struct {
	MaxRetries int
	RetryDelay time.Duration
	MaxBackoff time.Duration
}

type Client struct {
	conn        net.Conn
	address     string
	nodeID      string
	retryConfig RetryConfig
	security    SecurityConfig
	connected   bool
}

func NewClient(address, nodeID string, retryConfig RetryConfig, security SecurityConfig) *Client {
	return &Client{
		address:     address,
		nodeID:      nodeID,
		retryConfig: retryConfig,
		security:    security,
	}
}

func (c *Client) Connect() error {
	var lastErr error
	currentDelay := c.retryConfig.RetryDelay

	for i := 0; i < c.retryConfig.MaxRetries; i++ {
		var conn net.Conn
		var err error

		if c.security.UseTLS {
			tlsConfig, err := newTLSConfig(c.security)
			if err != nil {
				return err
			}
			conn, err = tls.Dial("tcp", c.address, tlsConfig)
		} else {
			conn, err = net.Dial("tcp", c.address)
		}

		if err == nil {
			c.conn = conn
			c.connected = true
			return nil
		}

		lastErr = err
		time.Sleep(currentDelay)

		// Exponential backoff
		currentDelay *= 2
		if currentDelay > c.retryConfig.MaxBackoff {
			currentDelay = c.retryConfig.MaxBackoff
		}
	}

	return errors.New("Max retry count exceeded: " + lastErr.Error())
}
