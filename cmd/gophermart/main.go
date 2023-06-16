package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Orendev/go-loyalty/internal/app"
	"github.com/Orendev/go-loyalty/internal/client"
	"github.com/Orendev/go-loyalty/internal/config"
	"github.com/Orendev/go-loyalty/internal/logger"
	"github.com/Orendev/go-loyalty/internal/models"
	"github.com/Orendev/go-loyalty/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := logger.NewLogger("info"); err != nil {
		log.Fatal(err)
	}

	shutdownTimeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	repo, err := repository.NewRepository(cfg.Database.URI)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = repo.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = repo.Bootstrap(ctx)
	if err != nil {
		log.Fatal(err)
	}

	accrualChain := make(chan models.Accrual, cfg.Size)

	a := app.NewApp(ctx, repo, accrualChain)

	_, err = client.NewHTTPClient(context.Background(), repo, cfg.AccrualSystem.Addr, accrualChain)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: a.Routes(chi.NewRouter()),
	}

	log.Fatal(srv.ListenAndServe())
}
