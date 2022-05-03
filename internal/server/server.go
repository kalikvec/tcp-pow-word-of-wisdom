package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/kalikvec/tcp-pow-word-of-wisdom/internal/env"
	"github.com/kalikvec/tcp-pow-word-of-wisdom/internal/pow"
	"github.com/kalikvec/tcp-pow-word-of-wisdom/internal/proto"
	"github.com/kalikvec/tcp-pow-word-of-wisdom/internal/transport"
)

var quotes = []string{
	"Self-Improvement and success go hand in hand. " +
		"Taking the steps to make yourself a better and more well-rounded individual will prove to be a wise decision.",

	"The wise person feels the pain of one arrow. The unwise feels the pain of two",

	"When looking for wise words, the best ones often come from our elders.",

	"Some of us think holding on makes us strong, but sometimes it is letting go.",

	"Don't waste your time with explanations, people only hear what they want to hear.",

	"To make difficult decisions wisely," +
		" it helps to have a systematic process for assessing each choice and its consequences " +
		"- the potential impact on each aspect of your life",

	"If we manage ego wisely, we get the upside it delivers followed by strong returns.",
}

var ErrFin = errors.New("client requests to close connection")

type server struct {
	cfg   *env.Config
	cache Cache
}

func NewServer(cfg *env.Config, cache Cache) *server {
	return &server{cfg: cfg, cache: cache}
}

// Run - main function, launches server to listen on given address and handle new connections
func (s *server) Run(ctx context.Context, address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Println("listening", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("accept connection error: %w", err)
		}

		go s.handleConnection(ctx, conn)
	}
}

func (s *server) handleConnection(ctx context.Context, conn net.Conn) {
	fmt.Println("new client:", conn.RemoteAddr())
	defer conn.Close()

	gate := transport.NewGate(conn)
	for {
		msgIn, err := gate.ReadMessage()
		if err != nil {
			fmt.Println("read ResponseChallenge message error: ", err)
			return
		}

		msg, err := s.handleMsg(ctx, msgIn, conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("handleMsg error:", err)
			return
		}

		err = gate.WriteMessage(&msg)
		if err != nil {
			fmt.Println("WriteMessage error:", err)
		}

	}
}

func (s *server) handleMsg(ctx context.Context, in *proto.Message, clientAddress string) (proto.Message, error) {
	switch in.Header {
	case proto.Fin:
		return proto.Message{}, ErrFin
	case proto.RequestChallenge:
		fmt.Printf("client %s requests challenge\n", clientAddress)
		// add new created rand value to cache to check it later on RequestResource stage
		// with duration in seconds
		randValue := rand.Intn(100000)
		err := s.cache.Add(randValue, s.cfg.ChallengeTimeout)
		if err != nil {
			return proto.Message{}, fmt.Errorf("add rand to cache error: %w", err)
		}

		challenge := pow.HashCash{
			Version:    1,
			ZerosCount: s.cfg.HashcashZerosCount,
			Date:       time.Now().Unix(),
			Resource:   clientAddress,
			Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", randValue))),
			Counter:    0,
		}
		payload, err := json.Marshal(challenge)
		if err != nil {
			return proto.Message{}, fmt.Errorf("err marshal challenge: %v", err)
		}
		msg := proto.Message{
			Header:  proto.ResponseChallenge,
			Payload: payload,
		}
		return msg, nil
	case proto.RequestResource:
		fmt.Printf("client %s requests resource with payload %v\n", clientAddress, in.Payload)
		// parse client's solution
		var h pow.HashCash
		err := json.Unmarshal(in.Payload, &h)
		if err != nil {
			return proto.Message{}, fmt.Errorf("unmarshal hashcash error: %w", err)
		}

		// validate hashcash params
		if h.Resource != clientAddress {
			return proto.Message{}, fmt.Errorf("invalid hashcash resource")
		}

		randValBytes, err := base64.StdEncoding.DecodeString(h.Rand)
		if err != nil {
			return proto.Message{}, fmt.Errorf("decode rand error: %w", err)
		}
		randVal, err := strconv.Atoi(string(randValBytes))
		if err != nil {
			return proto.Message{}, fmt.Errorf("decode rand error: %w", err)
		}

		// check whether rand value exists in cache, that means that user's already requested challenge
		exists, err := s.cache.Get(randVal)
		if err != nil {
			return proto.Message{}, fmt.Errorf("get rand from cache error: %w", err)
		}
		if !exists {
			return proto.Message{}, fmt.Errorf("client has not requested challenge")
		}

		// sent solution should not be outdated
		if time.Now().Unix()-h.Date > s.cfg.ChallengeTimeout {
			return proto.Message{}, fmt.Errorf("challenge expired")
		}

		if !pow.IsValidHashCash(h) {
			return proto.Message{}, fmt.Errorf("invalid hashcash")
		}

		//get random quote
		fmt.Printf("client %s succesfully computed hashcash %s\n", clientAddress, in.Payload)
		msg := proto.Message{
			Header:  proto.ResponseResource,
			Payload: []byte(quotes[rand.Intn(len(quotes))]),
		}

		s.cache.Delete(randVal)
		return msg, nil
	default:
		return proto.Message{}, fmt.Errorf("unknown header")
	}
}
