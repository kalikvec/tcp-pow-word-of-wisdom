package transport

import (
	"bufio"
	"bytes"
	"io"

	"github.com/kalikvec/tcp-pow-word-of-wisdom/internal/proto"
)

const (
	msgSep = '\n'
)

type MessageGate struct {
	writer io.Writer
	reader *bufio.Reader
}

func NewGate(conn io.ReadWriter) *MessageGate {
	return &MessageGate{writer: conn, reader: bufio.NewReader(conn)}
}

func (s *MessageGate) WriteMessage(m *proto.Message) error {
	b := bytes.NewBuffer(m.Marshall())
	b.WriteByte(msgSep)

	_, err := s.writer.Write(b.Bytes())
	return err
}

func (s *MessageGate) ReadMessage() (*proto.Message, error) {
	b, err := s.reader.ReadBytes(msgSep)
	if err != nil {
		return nil, err
	}
	return proto.NewMessage(b)
}
