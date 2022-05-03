package proto

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	separator = "|"

	Fin = iota
	RequestChallenge
	ResponseChallenge
	RequestResource
	ResponseResource
)

type Message struct {
	Header  int
	Payload []byte
}

func (m *Message) Marshall() []byte {
	return []byte(fmt.Sprintf("%d:%s", m.Header, m.Payload))
}

func NewMessage(b []byte) (*Message, error) {
	parts := bytes.Split(bytes.TrimSpace(b), []byte(separator))
	lenParts := len(parts)
	if lenParts < 1 || lenParts > 2 {
		return nil, fmt.Errorf("incorrect message structure")
	}

	msgType, err := strconv.Atoi(string(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("cannot parse header")
	}
	msg := Message{
		Header: msgType,
	}
	if len(parts) == 2 {
		msg.Payload = parts[1]
	}
	return &msg, nil
}
