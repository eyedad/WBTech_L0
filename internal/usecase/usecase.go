package usecase

import (
	"time"

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

func (u *Usecase) GetAllOrders() ([]string, error) {

	orders, err := u.r.GetAllOrdersFromCache()
	if err == redis.Nil {
		orders, err = u.r.GetAllOrdersFromDB()
		time.Sleep(3 * time.Second)
		if err != nil {
			return nil, err
		}
		err = u.r.InsertAllOrdersIntoCache(orders)
		if err != nil {
			return nil, err
		}
	}
	return orders, nil
}

func (u *Usecase) GetOrderById(order *entity.Order, orderUID string) error {
	order, err := u.r.GetOrderFromCache(order, orderUID)
	if err == redis.Nil {
		err := u.r.GetOrderFromDB(order, orderUID)
		time.Sleep(3 * time.Second)

		if err != nil {
			return err
		}
	}
	err = u.r.InserOrderIntoCache(order)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) InsertOrder(order *entity.Order) error {
	err := u.r.InsertOrderIntoDB(order)
	time.Sleep(3 * time.Second)
	if err != nil {
		return err
	}
	err = u.r.InserOrderIntoCache(order)
	if err != nil {
		return err
	}
	return nil
}
