package kernel

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TCPClient struct {
	*TCPBase
	connection   net.Conn
	timer        *time.Timer
	commandQueue []*TCPCommand
	listen       bool
	readsData    string
	writesData   []string
	read         []net.Conn
	write        []net.Conn
}

func NewTCPClient() *TCPClient {
	client := &TCPClient{
		TCPBase:      NewTCPBase(),
		timer:        time.NewTimer(0),
		listen:       true,
		commandQueue: make([]*TCPCommand, 0),
	}
	client.AddListener("response", client.response)
	return client
}

func (c *TCPClient) Connect(addr string, port int) bool {
	if c.IsConnected() {
		c.Disconnect()
	}

	if addr != "" {
		c.Addr = addr
	}
	if port != 0 {
		c.Port = port
	}

	uri := fmt.Sprintf("%s:%d", c.Addr, c.Port)
	conn, err := net.Dial("tcp", uri)
	if err != nil {
		fmt.Printf("[TCPClient] Connection to %s failed: %v\n", uri, err)
		return false
	}

	c.connection = conn
	return true
}

func (c *TCPClient) Disconnect() {
	if c.connection != nil {
		c.connection.Close()
		c.connection = nil
	}
}

func (c *TCPClient) IsConnected() bool {
	return c.connection != nil
}

func (c *TCPClient) AddWriteData(data string) {
	const chunkSize = 4096
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		c.writesData = append(c.writesData, data[i:end])
	}
}

func (c *TCPClient) RemoveWriteData() {
	c.writesData = nil
}

func (c *TCPClient) AddReadData(data string) {
	c.readsData += data
}

func (c *TCPClient) RemoveReadData() {
	c.readsData = ""
}

func (c *TCPClient) Init() {
	c.listen = true
	c.commandQueue = make([]*TCPCommand, 0)
	c.readsData = ""
	c.writesData = make([]string, 0)
}

func (c *TCPClient) End(disconnect bool) {
	c.listen = false
	c.available = false
	c.read = nil
	c.write = nil
	c.RemoveReadData()
	c.RemoveWriteData()
	if disconnect {
		c.Disconnect()
	}
}

func (c *TCPClient) Send(command *TCPCommand, timeout time.Duration) interface{} {
	c.Init()
	var response interface{}

	start := time.Now()
	encoded, _ := c.Encode(command)
	c.AddWriteData(string(encoded))

	for c.listen && time.Since(start) < timeout {
		c.SelectOperation()
		c.ReadOperation()

		for cmd := c.PopCommand(); cmd != nil; cmd = c.PopCommand() {
			response = c.Run(cmd)
			c.End(false)
		}

		c.WriteOperation()
	}

	if c.listen {
		c.End(false)
	}

	return response
}

func (c *TCPClient) Response(command *TCPCommand) interface{} {
	return command.Data()
}

func (c *TCPClient) AddCommand(command *TCPCommand) {
	c.commandQueue = append(c.commandQueue, command)
}

func (c *TCPClient) PopCommand() *TCPCommand {
	if len(c.commandQueue) == 0 {
		return nil
	}
	command := c.commandQueue[0]
	c.commandQueue = c.commandQueue[1:]
	return command
}

func (c *TCPClient) SelectOperation() {
	c.available = false
	if c.connection == nil {
		return
	}

	c.read = []net.Conn{c.connection}
	c.write = []net.Conn{c.connection}
}

func (c *TCPClient) WriteOperation() {
	if !c.available {
		return
	}

	for _, conn := range c.write {
		if len(c.writesData) == 0 {
			return
		}

		if conn == nil {
			fmt.Println("[TCPClient] There is no connection.")
			c.End()
			return
		}

		data := c.writesData[0]
		c.writesData = c.writesData[1:]

		n, err := conn.Write([]byte(data))
		if err != nil || n == 0 {
			fmt.Println("[TCPClient] Unable to write data.")
			c.End()
			return
		}
	}
}

func (c *TCPClient) ReadOperation() {
	if !c.available {
		return
	}

	buf := make([]byte, 4096)
	for _, conn := range c.read {
		n, err := conn.Read(buf)
		if err == io.EOF || n == 0 {
			fmt.Println("[TCPClient] Unable to read.")
			c.End()
			break
		}

		c.AddReadData(string(buf[:n]))
		dataLen := len(c.readsData)
		cmdLen := c.GetLength([]byte(c.readsData), dataLen)

		if cmdLen > c.ReadMaxLength {
			fmt.Println("[TCPClient] Invalid command.")
			c.End()
			return
		}

		if cmdLen == 0 || dataLen < cmdLen {
			// insufficient data
			return
		}

		command, _ := c.Decode([]byte(c.readsData[:cmdLen]))
		c.readsData = c.readsData[cmdLen:]
		c.AddCommand(command)
		c.End()
	}
}
