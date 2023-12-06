package service

import (
	"github.com/43nvy/wb_l0"
	"github.com/43nvy/wb_l0/internal/repository"
	"github.com/nats-io/stan.go"
)

type WbOrder interface {
	CreateOrder(order wb_l0.Order) (int, error)
	GetOrder(id int) (wb_l0.Order, error)
	GetTenOrders() ([]wb_l0.Order, error)
}

type NATS interface {
	SubscribeClient(client_id string) error
	UnsubscribeClient(client_id string)
	NotifyNewOrder(order_id int)
}

type Service struct {
	WbOrder
	NATS
}

func NewService(repos *repository.Repository, sc stan.Conn) *Service {
	return &Service{
		WbOrder: NewOrderService(repos.WbOrder),
		NATS:    NewNATS(sc),
	}
}
