package v1

import (
	"encoding/json"
	"net/http"

	"example.com/m/v2/internal/entity"
	"example.com/m/v2/internal/usecase"
	"example.com/m/v2/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
	logger  *logging.Logger
	usecase *usecase.Usecase
}

func NewHandler(logger *logging.Logger, usecase *usecase.Usecase) Handler {
	return &handler{
		logger:  logger,
		usecase: usecase,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET("/orders", h.GetAllOrders)
	router.GET("/orders/:id", h.GetOrderById)
	router.POST("/orders/:id", h.AddOrder)
}

func (h *handler) GetAllOrders(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	h.logger.Info("New GET request")

	orders, err := h.usecase.GetAllOrders()
	if err != nil {
		h.logger.Errorf("Failed to get orders: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Info("All orders ids are being viewed")
}

func (h *handler) GetOrderById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	h.logger.Info("New GET request")
	orderUID := params.ByName("id")

	var order entity.Order

	err := h.usecase.GetOrderById(&order, orderUID)
	if err != nil {
		h.logger.Errorf("Failed to get order %s: %v", orderUID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Infof("Order is being viewed, id: %s", order.OrderUID)
}

func (h *handler) AddOrder(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	h.logger.Info("New POST request")

	var order entity.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Internal server error", http.StatusBadRequest)
		return
	}

	err := h.usecase.InsertOrder(&order)
	if err != nil {
		h.logger.Errorf("Failed to insert order: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order created successfully"})
	h.logger.Infof("Order created successfully, id: %s", order.OrderUID)
}
