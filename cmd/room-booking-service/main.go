package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/app"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/core/logger"
)

func main() {
	cmd := osArgsOrDefault(1, "serve")

	application, err := app.NewApp(context.Background())
	if err != nil {
		log.Fatalf("failed to build app: %v", err)
	}

	switch cmd {
	case "serve":
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		if err := application.Run(ctx); err != nil {
			logger.Error("serve", "err", err)
			os.Exit(1)
		}

	case "seed":
		if err := application.Seed(context.Background()); err != nil {
			logger.Error("seed", "err", err)
			os.Exit(1)
		}
		logger.Info("seed completed")

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		os.Exit(2)
	}
}

func osArgsOrDefault(idx int, fallback string) string {
	if len(os.Args) > idx {
		return os.Args[idx]
	}
	return fallback
}
