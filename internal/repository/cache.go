package repository

import (
	"context"
	"encoding/json"
	"time"

	"example.com/m/v2/internal/entity"
	"github.com/go-redis/redis/v8"
)

func (r *Repository) GetAllOrdersFromCache() ([]string, error) {
	data, err := r.cache.Get(context.Background(), "orders").Result()
	if err != nil {
		return nil, err
	}
	if data == "null" {
		return nil, redis.Nil
	}

	var orders []string

	err = json.Unmarshal([]byte(data), &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
func (r *Repository) InsertAllOrdersIntoCache(orders []string) error {
	jsonData, err := json.Marshal(orders)
	if err != nil {
		return err
	}

	err = r.cache.Set(context.Background(), "orders", jsonData, 10*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetOrderFromCache(order *entity.Order, key string) (*entity.Order, error) {
	data, err := r.cache.Get(context.Background(), key).Result()
	if err != nil {
		return order, err
	}

	err = json.Unmarshal([]byte(data), &order)
	if err != nil {
		return order, err
	}

	return order, nil
}

func (r *Repository) InserOrderIntoCache(order *entity.Order) error {
	jsonData, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = r.cache.Set(context.Background(), order.OrderUID, jsonData, 10*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}
