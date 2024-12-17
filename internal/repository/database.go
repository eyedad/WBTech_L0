package repository

import (
	"fmt"

	"example.com/m/v2/internal/entity"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) InsertOrderIntoDB(order *entity.Order) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err := insertIntoOrders(tx, order); err != nil {
		return err
	}
	if err := insertIntoDeliveries(tx, order); err != nil {
		return err
	}
	if err := insertIntoPayments(tx, order); err != nil {
		return err
	}
	if err := insertIntoItems(tx, order); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func insertIntoOrders(tx *sqlx.Tx, order *entity.Order) error {
	query := `
        INSERT INTO orders (order_uid, track_number, entry, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, locale, internal_signature)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `
	_, err := tx.Exec(query, order.OrderUID, order.TrackNumber, order.Entry, order.CustomerID,
		order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard,
		order.Locale, order.InternalSignature)
	return err
}

func insertIntoDeliveries(tx *sqlx.Tx, order *entity.Order) error {
	query := `
        INSERT INTO deliveries (name, phone, zip, city, address, region, email, order_uid)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	_, err := tx.Exec(query, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
		order.OrderUID)
	return err
}

func insertIntoPayments(tx *sqlx.Tx, order *entity.Order) error {
	query := `
        INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee, order_uid)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `
	_, err := tx.Exec(query, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee, order.OrderUID)
	return err
}

func insertIntoItems(tx *sqlx.Tx, order *entity.Order) error {
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

func (r *Repository) GetOrderFromDB(order *entity.Order, orderUID string) error {

	if err := getFromOrders(r.db, order, orderUID); err != nil {
		fmt.Printf("1, %s", err.Error())
		return err
	}
	if err := getFromDeliveries(r.db, order, orderUID); err != nil {
		fmt.Printf("2, %s", err.Error())
		return err
	}
	if err := getFromPayments(r.db, order, orderUID); err != nil {
		fmt.Printf("3, %s", err.Error())
		return err
	}
	if err := getFromItems(r.db, order, orderUID); err != nil {
		fmt.Printf("4, %s", err.Error())
		return err
	}

	return nil
}

func getFromOrders(db *sqlx.DB, order *entity.Order, orderUID string) error {
	query := `
	SELECT order_uid, track_number, entry, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, locale, internal_signature
	FROM orders WHERE order_uid = $1
`

	return db.Get(order, query, orderUID)
}

func getFromDeliveries(db *sqlx.DB, order *entity.Order, orderUID string) error {
	query := `
        SELECT name, phone, zip, city, address, region, email
        FROM deliveries WHERE order_uid = $1
    `

	return db.Get(&order.Delivery, query, orderUID)

}

func getFromPayments(db *sqlx.DB, order *entity.Order, orderUID string) error {
	query := `
        SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
        FROM payments WHERE order_uid = $1
    `

	return db.Get(&order.Payment, query, orderUID)
}

func getFromItems(db *sqlx.DB, order *entity.Order, orderUID string) error {
	query := `
	SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	FROM items WHERE order_uid = $1
`

	return db.Select(&order.Items, query, orderUID)
}

func (r *Repository) GetAllOrdersFromDB() ([]string, error) {
	var orders []string

	err := r.db.Select(&orders, "SELECT order_uid FROM orders")

	return orders, err
}
