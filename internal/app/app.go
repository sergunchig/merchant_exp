package app

import (
	"context"
	"fmt"
	"merchant_exp/config"
	"merchant_exp/internal/server"
	"time"
)

func Run(cfg *config.Config) {
	fmt.Println("start project...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5) //WithCancel(context.Background())
	defer cancel()

	srv, err := server.New(cfg)
	fmt.Println("server was started...")
	if err != nil {
		fmt.Println(err)
	}

	<-ctx.Done()

	srv.Shutdown()

	fmt.Println("app was stoped")
}
