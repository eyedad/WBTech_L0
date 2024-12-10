package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"example.com/m/v2/internal/order"
	"example.com/m/v2/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

func ViewHandler(w http.ResponseWriter, r http.Request, params httprouter.Params) {
	name := params.ByName("name")
	w.Write([]byte(fmt.Sprintf("Hello %s", name)))
}

func main() {
	logger := logging.GetLogger()

	logger.Info("Creating router")
	router := httprouter.New()

	logger.Info("Register handlers")
	handler := order.NewHandler(logger)
	handler.Register(router)

	logger.Info("Starting server")
	start(router)
}

func start(router *httprouter.Router) {
	logger := logging.GetLogger()

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 12 * time.Second,
		ReadTimeout:  12 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
