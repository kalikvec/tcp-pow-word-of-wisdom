package main

import (
	"context"
	"fmt"

	"github.com/kalikvec/tcp-pow-word-of-wisdom/internal/cache"
	"github.com/kalikvec/tcp-pow-word-of-wisdom/internal/env"
	"github.com/kalikvec/tcp-pow-word-of-wisdom/internal/server"
)

func main() {
	fmt.Println("start server")

	cfg, err := env.LoadConfig("config/config.json")
	if err != nil {
		fmt.Println("load config error:", err)
		return
	}

	c := cache.NewMemory()
	srv := server.NewServer(cfg, c)

	err = srv.Run(context.Background())
	if err != nil {
		fmt.Println("client error:", err)
	}
}
