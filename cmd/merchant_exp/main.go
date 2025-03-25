package main

import (
	"github.com/sergunchig/merchant_exp.git/config"
	"github.com/sergunchig/merchant_exp.git/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic("Cannt read config")
	}
	app.Run(cfg)
}
