package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/m/v2/internal/config"
	"example.com/m/v2/internal/order"
	"example.com/m/v2/pkg/logging"
	"example.com/m/v2/pkg/postgers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()

	logger.Info("Creating router")
	router := httprouter.New()

	logger.Info("Reading configuration")
	cfg := config.GetConfig()

	logger.Info("Connecting to database")
	db := postgers.New(cfg, logger)

	logger.Info("Register handlers")
	handler := order.NewHandler(logger, db)
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
	if err := db.Close(); err != nil {
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
