package kernel

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

type TCPBase struct {
	ReadMaxLength int
	IdBytesSize   int
	TypeBytesSize int
	DataBytesSize int
	PrefixSize    int
	Addr          string
	Port          int
	listeners     map[string]func(*TCPCommand) interface{}
	available     bool
}

func NewTCPBase() *TCPBase {
	return &TCPBase{
		ReadMaxLength: 67108864,
		IdBytesSize:   1,
		TypeBytesSize: 1,
		DataBytesSize: 4,
		PrefixSize:    6, // IdBytesSize + TypeBytesSize + DataBytesSize
		Addr:          "127.0.0.1",
		Port:          9933,
		listeners:     make(map[string]func(*TCPCommand) interface{}),
		available:     false,
	}
}

func (t *TCPBase) AddListener(cmdType string, handler func(*TCPCommand) interface{}) {
	t.listeners[cmdType] = handler
}

func (t *TCPBase) Run(command *TCPCommand) interface{} {
	if command.Type() == "" {
		return nil
	}

	handler, exists := t.listeners[command.Type()]
	if !exists {
		return nil
	}

	return handler(command)
}

func (t *TCPBase) GetLength(encodedCommand []byte, dataLength int) int {
	if dataLength < 6 {
		return 0
	}

	idBytes := int(binary.BigEndian.Uint64(encodedCommand[:t.IdBytesSize]))
	typeBytes := int(binary.BigEndian.Uint64(encodedCommand[t.IdBytesSize : t.IdBytesSize+t.TypeBytesSize]))
	dataBytes := int(binary.BigEndian.Uint64(encodedCommand[t.IdBytesSize+t.TypeBytesSize : t.PrefixSize]))

	return t.PrefixSize + idBytes + typeBytes + dataBytes
}

func (t *TCPBase) Encode(command *TCPCommand) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(command.Data()); err != nil {
		return nil, err
	}
	data := buf.Bytes()

	idLen := len(command.Id())
	typeLen := len(command.Type())
	dataLen := len(data)

	result := make([]byte, t.PrefixSize+idLen+typeLen+dataLen)

	binary.BigEndian.PutUint64(result[:t.IdBytesSize], uint64(idLen))
	binary.BigEndian.PutUint64(result[t.IdBytesSize:t.IdBytesSize+t.TypeBytesSize], uint64(typeLen))
	binary.BigEndian.PutUint64(result[t.IdBytesSize+t.TypeBytesSize:t.PrefixSize], uint64(dataLen))

	copy(result[t.PrefixSize:], []byte(command.Id()))
	copy(result[t.PrefixSize+idLen:], []byte(command.Type()))
	copy(result[t.PrefixSize+idLen+typeLen:], data)

	return result, nil
}

func (t *TCPBase) Decode(encodedCommand []byte) (*TCPCommand, error) {
	idBytes := int(binary.BigEndian.Uint64(encodedCommand[:t.IdBytesSize]))
	typeBytes := int(binary.BigEndian.Uint64(encodedCommand[t.IdBytesSize : t.IdBytesSize+t.TypeBytesSize]))
	dataBytes := int(binary.BigEndian.Uint64(encodedCommand[t.IdBytesSize+t.TypeBytesSize : t.PrefixSize]))

	id := string(encodedCommand[t.PrefixSize : t.PrefixSize+idBytes])
	cmdType := string(encodedCommand[t.PrefixSize+idBytes : t.PrefixSize+idBytes+typeBytes])

	var data interface{}
	buf := bytes.NewBuffer(encodedCommand[t.PrefixSize+idBytes+typeBytes : t.PrefixSize+idBytes+typeBytes+dataBytes])
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}

	command := NewTCPCommand()
	command.SetId(id)
	command.SetType(cmdType)
	command.SetData(data)

	return command, nil
}
