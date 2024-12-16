package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"example.com/m/v2/config"
	"example.com/m/v2/internal/app"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
		<-sigChan
		cancel()
	}()

	app.Run(ctx, cfg)
}
