package main

import (
	"merchant_exp/config"
	"merchant_exp/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic("Cannt read config")
	}
	app.Run(cfg)
}
