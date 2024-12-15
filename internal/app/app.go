package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/m/v2/config"
	v1 "example.com/m/v2/internal/delivery/http/v1"
	"example.com/m/v2/internal/repository"
	"example.com/m/v2/internal/usecase"
	"example.com/m/v2/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

func Run(cfg *config.Config) {
	logger := logging.GetLogger()

	logger.Info("Creating router")
	router := httprouter.New()

	logger.Info("Register handlers")
	repository := repository.New(cfg, logger)
	usecase := usecase.New(repository)
	handler := v1.NewHandler(logger, usecase)
	handler.Register(router)

	logger.Info("Connecting to server")
	server := &http.Server{
		Handler:      router,
		WriteTimeout: time.Duration(cfg.Listen.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Listen.ReadTimeout) * time.Second,
	}
	go start(cfg, server)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("Shutting down")
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Errorf("Error occured while server shutting down, error: %s", err.Error())
	}
	if err := repository.Close(); err != nil {
		logger.Errorf("Error occured while closing db connection, error: %s", err.Error())
	}
}

func start(cfg *config.Config, server *http.Server) {
	logger := logging.GetLogger()
	logger.Info("Starting server")

	var listener net.Listener
	var err error

	logger.Info("Listen tcp")
	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	if err != nil {
		panic(err)
	}
	logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)

	logger.Fatal(server.Serve(listener))
}
