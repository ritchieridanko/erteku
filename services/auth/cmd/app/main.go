package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ritchieridanko/erteku/services/auth/configs"
	"github.com/ritchieridanko/erteku/services/auth/internal/di"
	"github.com/ritchieridanko/erteku/services/auth/internal/infra"
	"github.com/ritchieridanko/erteku/services/auth/internal/transport/server"
)

func main() {
	cfg, err := configs.Init("./configs")
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}

	i, err := infra.Init(cfg)
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}

	defer func(i *infra.Infra) {
		if err := i.Close(); err != nil {
			log.Println("[WARN]:", err)
		}
	}(i)

	ctr := di.Init(cfg, i)
	srv := ctr.Server()

	// Start the server
	go func(srv *server.Server) {
		if err := srv.Start(); err != nil {
			log.Fatalln("[FATAL]:", err)
		}
	}(srv)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Printf("[%s] is shutting down...", strings.ToUpper(cfg.App.Name))

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.Timeout.Shutdown)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("[WARN]:", err)
	}
}
