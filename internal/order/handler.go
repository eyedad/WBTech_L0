package order

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"example.com/m/v2/internal/handlers"
	"example.com/m/v2/pkg/logging"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

const (
	ordersURl  = "/orders"
	orderURl   = "/orders/:id"
	queryOrder = `
        SELECT order_uid, track_number, entry, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, locale, internal_signature
        FROM orders
        WHERE order_uid = $1
    `
	queryDelivery = `
        SELECT name, phone, zip, city, address, region, email
        FROM deliveries
        WHERE order_uid = $1
    `
	queryPayment = `
        SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
        FROM payments
        WHERE order_uid = $1
    `
	queryItems = `
        SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
        FROM items
        WHERE order_uid = $1
    `
)

type handler struct {
	logger *logging.Logger
	db     *sqlx.DB
}

func NewHandler(logger *logging.Logger, db *sqlx.DB) handlers.Handler {
	return &handler{
		logger: logger,
		db:     db,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(ordersURl, h.GetAllOrders)
	router.GET(orderURl, h.GetOrderById)
	router.POST(ordersURl, h.AddOrder)
}

func (h *handler) GetAllOrders(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	h.logger.Info("New GET request")

	var orders []string

	err := h.db.Select(&orders, "SELECT order_uid FROM orders")
	if err != nil {
		h.logger.Fatal(err)
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

	var order Order

	err := h.db.Get(&order, queryOrder, orderUID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		h.logger.Errorf("Failed to get order: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.db.Get(&order.Delivery, queryDelivery, orderUID)
	if err != nil {
		h.logger.Errorf("Failed to get delivery details: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.db.Get(&order.Payment, queryPayment, orderUID)
	if err != nil {
		h.logger.Errorf("Failed to get payment details: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.db.Select(&order.Items, queryItems, orderUID)
	if err != nil {
		h.logger.Errorf("Failed to get items: %v", err)
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

	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Internal server error", http.StatusBadRequest)
		return
	}

	tx, err := h.db.Beginx()
	if err != nil {
		h.logger.Errorf("Failed to begin transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusBadRequest)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err := h.insertOrder(tx, order); err != nil {
		h.logger.Error(err)
		http.Error(w, "Internal server error", http.StatusBadRequest)
		return
	}

	if err := tx.Commit(); err != nil {
		h.logger.Errorf("Failed to commit transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order created successfully"})
	h.logger.Infof("Order created successfully, id: %s", order.OrderUID)
}

func (h *handler) insertOrder(tx *sqlx.Tx, order Order) error {
	if err := h.insertIntoOrders(tx, order); err != nil {
		return err
	}
	if err := h.insertIntoDeliveries(tx, order); err != nil {
		return err
	}
	if err := h.insertIntoPayments(tx, order); err != nil {
		return err
	}
	if err := h.insertIntoItems(tx, order); err != nil {
		return err
	}
	return nil
}

func (h *handler) insertIntoOrders(tx *sqlx.Tx, order Order) error {
	query := `
        INSERT INTO orders (order_uid, track_number, entry, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, locale, internal_signature)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `
	_, err := tx.Exec(query, order.OrderUID, order.TrackNumber, order.Entry, order.CustomerID,
		order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard,
		order.Locale, order.InternalSignature)
	return err
}

func (h *handler) insertIntoDeliveries(tx *sqlx.Tx, order Order) error {
	query := `
        INSERT INTO deliveries (name, phone, zip, city, address, region, email, order_uid)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	_, err := tx.Exec(query, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
		order.OrderUID)
	return err
}

func (h *handler) insertIntoPayments(tx *sqlx.Tx, order Order) error {
	query := `
        INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee, order_uid)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `
	_, err := tx.Exec(query, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee, order.OrderUID)
	return err
}

func (h *handler) insertIntoItems(tx *sqlx.Tx, order Order) error {
	query := `
        INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_uid)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `
	for _, item := range order.Items {
		if _, err := tx.Exec(query, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status, order.OrderUID); err != nil {
			return err
		}
	}
	return nil
}
