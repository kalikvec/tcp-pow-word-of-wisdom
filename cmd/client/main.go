package main

import (
	"context"
	"fmt"

	"github.com/kalikvec/tcp-pow-word-of-wisdom/internal/client"
	"github.com/kalikvec/tcp-pow-word-of-wisdom/internal/env"
)

func main() {
	fmt.Println("start client")

	cfg, err := env.LoadConfig("config/config.json")
	if err != nil {
		fmt.Println("load config error:", err)
		return
	}

	c := client.NewClient(cfg)
	err = c.Run(context.Background())
	if err != nil {
		fmt.Println("client error:", err)
	}
}
