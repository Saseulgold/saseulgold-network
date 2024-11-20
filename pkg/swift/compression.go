package swift

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

type CompressionType string

const (
	NoCompression   CompressionType = "none"
	GzipCompression CompressionType = "gzip"
)

type CompressedMessage struct {
	Message
	Compressed      bool            `json:"compressed"`
	CompressionType CompressionType `json:"compression_type"`
}

func compressData(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decompressData(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}
