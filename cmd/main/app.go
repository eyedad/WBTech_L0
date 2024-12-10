package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"example.com/m/v2/internal/config"
	"example.com/m/v2/internal/order"
	"example.com/m/v2/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()

	logger.Info("Creating router")
	router := httprouter.New()

	cfg := config.GetConfig()

	handler := order.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
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

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 12 * time.Second,
		ReadTimeout:  12 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
