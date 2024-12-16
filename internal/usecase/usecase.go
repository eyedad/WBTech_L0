package usecase

import (
	"context"

	"example.com/m/v2/internal/entity"
	"example.com/m/v2/internal/repository"
	"github.com/go-redis/redis/v8"
)

type Usecase struct {
	r *repository.Repository
}

func New(r *repository.Repository) *Usecase {
	return &Usecase{r}
}

func (u *Usecase) GetAllOrders(ctx context.Context) ([]string, error) {

	orders, err := u.r.GetAllOrdersFromCache(ctx)
	if err == redis.Nil {
		orders, err = u.r.GetAllOrdersFromDB()
		if err != nil {
			return nil, err
		}
		err = u.r.InsertAllOrdersIntoCache(ctx, orders)
		if err != nil {
			return nil, err
		}
	}
	return orders, nil
}

func (u *Usecase) GetOrderById(ctx context.Context, order *entity.Order, orderUID string) error {
	order, err := u.r.GetOrderFromCache(ctx, order, orderUID)
	if err == redis.Nil {
		err := u.r.GetOrderFromDB(order, orderUID)

		if err != nil {
			return err
		}
	}
	err = u.r.InserOrderIntoCache(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) InsertOrder(ctx context.Context, order *entity.Order) error {
	err := u.r.InsertOrderIntoDB(order)
	if err != nil {
		return err
	}
	err = u.r.InserOrderIntoCache(ctx, order)
	if err != nil {
		return err
	}
	return nil
}
