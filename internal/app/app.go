package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Orendev/go-loyalty/internal/client"
	"github.com/Orendev/go-loyalty/internal/config"
	"github.com/Orendev/go-loyalty/internal/logger"
	"github.com/Orendev/go-loyalty/internal/models"
	"github.com/Orendev/go-loyalty/internal/repository"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type App struct {
	repo        repository.Storage
	accrualChan chan models.Accrual
}

var shutdownTimeout = 10 * time.Second

func Run(cfg *config.Config) {
	ctx := gracefulShutdown()

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

	a := NewApp(ctx, repo, accrualChain)

	_, err = client.NewHTTPClient(context.Background(), repo, cfg.AccrualSystem.Addr, accrualChain)
	if err != nil {
		logger.Log.Error("failed to start client", zap.Error(err))
	}

	startServer(ctx, &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: a.Routes(chi.NewRouter()),
	})
}

func NewApp(_ context.Context, repo repository.Storage, accrualChan chan models.Accrual) *App {
	return &App{repo: repo, accrualChan: accrualChan}
}

func startServer(ctx context.Context, srv *http.Server) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start server %s", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatalf("failed to shudown server %s", err)
	}

	wg.Wait()
}

func gracefulShutdown() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	irqSig := make(chan os.Signal, 1)
	// Получено сообщение о завершении работы от операционной системы.
	signal.Notify(irqSig, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-irqSig
		cancel()
	}()
	return ctx
}
