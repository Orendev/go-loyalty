package main

import (
	"log"

	"github.com/Orendev/go-loyalty/internal/app"
	"github.com/Orendev/go-loyalty/internal/config"
	"github.com/Orendev/go-loyalty/internal/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := logger.NewLogger("info"); err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
