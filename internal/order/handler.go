package order

import (
	"fmt"
	"net/http"

	"example.com/m/v2/internal/handlers"
	"example.com/m/v2/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

const (
	ordersURl = "/orders"
	orderURl  = "/orders/:id"
)

type handler struct {
	loger logging.Logger
}

func NewHandler(logger logging.Logger) handlers.Handler {
	return &handler{
		loger: logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(ordersURl, h.GetAllOrders)
	router.GET(orderURl, h.GetOrderById)
	router.POST(ordersURl, h.AddOrder)
}

func (h *handler) GetAllOrders(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("all orders"))
}

func (h *handler) GetOrderById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	w.Write([]byte(fmt.Sprintf("its %s order", id)))
}

func (h *handler) AddOrder(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("add orders"))
}
