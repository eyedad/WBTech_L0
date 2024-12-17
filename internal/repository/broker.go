package repository

import (
	"encoding/json"

	"example.com/m/v2/internal/entity"
)

func (r *Repository) Produce(topic string, order *entity.Order) error {
	orderMessage, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return r.producer.Produce(topic, order.OrderUID, orderMessage)
}
