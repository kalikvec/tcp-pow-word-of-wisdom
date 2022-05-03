package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/kalikvec/tcp-pow-word-of-wisdom/internal/env"
	"github.com/kalikvec/tcp-pow-word-of-wisdom/internal/pow"
	"github.com/kalikvec/tcp-pow-word-of-wisdom/internal/proto"
	"github.com/kalikvec/tcp-pow-word-of-wisdom/internal/transport"
)

type client struct {
	cfg *env.Config
}

func NewClient(cfg *env.Config) *client {
	return &client{cfg: cfg}
}

// Run - start tcp client
func (c *client) Run(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", c.cfg.ServerHost, c.cfg.ServerPort)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	fmt.Println("connected to", address)
	defer conn.Close()

	// client will send new request every 3 seconds endlessly
	for {
		message, err := c.handleConnection(ctx, conn)
		if err != nil {
			return err
		}
		fmt.Println("quote: ", message)
		time.Sleep(3 * time.Second)
	}
}

// step-by-step sequential flow for requesting resource - quote
// 1. requesting challenge
// 2. calc hashcash to perform Proof of Work
// 3. write message with found hashcash proof back to server
// 4. get result quote from server
func (c *client) handleConnection(ctx context.Context, conn io.ReadWriter) (string, error) {
	gate := transport.NewGate(conn)

	// 1. requesting challenge
	err := gate.WriteMessage(&proto.Message{
		Header: proto.RequestChallenge,
	})
	if err != nil {
		return "", fmt.Errorf("write RequestChallenge message error: %w", err)
	}

	msg, err := gate.ReadMessage()
	if err != nil {
		return "", fmt.Errorf("read ResponseChallenge message error: %w", err)
	}

	var h pow.HashCash
	err = json.Unmarshal(msg.Payload, &h)
	if err != nil {
		return "", fmt.Errorf("unmarshall hashcash error: %w", err)
	}

	// 2. calc hashcash to perform Proof of Work
	h, err = pow.CalcHashCash(h, c.cfg.HashcashMaxIterations)
	if err != nil {
		return "", fmt.Errorf("calc hashcash error: %w", err)
	}

	b, err := json.Marshal(h)
	if err != nil {
		return "", fmt.Errorf("marshall hashcash error: %w", err)
	}

	// 3. write message with found hashcash proof back to server
	err = gate.WriteMessage(&proto.Message{
		Header:  proto.RequestResource,
		Payload: b,
	})
	if err != nil {
		return "", fmt.Errorf("write message error: %w", err)
	}

	// 4. get result quote from server
	msg, err = gate.ReadMessage()
	if err != nil {
		return "", fmt.Errorf("read ResponseResource message error: %w", err)
	}
	return string(msg.Payload), nil
}
