package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/InternalTransfer/internal/config"
	"github.com/InternalTransfer/internal/database"
	"github.com/InternalTransfer/internal/handler"
	"github.com/InternalTransfer/internal/repository"
	"github.com/InternalTransfer/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("application exited with error", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := database.NewPool(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer pool.Close()
	logger.Info("connected to database")

	accountRepo := repository.NewAccountRepository(pool)
	transactionRepo := repository.NewTransactionRepository(pool)
	txManager := database.NewTxManager(pool)

	accountSvc := service.NewAccountService(accountRepo, logger)
	transferSvc := service.NewTransferService(accountRepo, transactionRepo, txManager, logger, cfg.MaxTransferAmount)

	accountHandler := handler.NewAccountHandler(accountSvc, logger)
	transactionHandler := handler.NewTransactionHandler(transferSvc, logger)

	router := handler.NewRouter(accountHandler, transactionHandler, logger)

	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		<-sigCh

		logger.Info("shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		srv.Shutdown(shutdownCtx)
	}()

	logger.Info("server starting", "port", cfg.ServerPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	logger.Info("server stopped")
	return nil
}
