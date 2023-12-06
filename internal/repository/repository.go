package repository

import (
	"github.com/43nvy/wb_l0"
	"github.com/jmoiron/sqlx"
)

type WbOrder interface {
	CreateOrder(message wb_l0.Order) (int, error)
	GetOrder(id int) (wb_l0.Order, error)
	GetTenOrders() ([]wb_l0.Order, error)
}

type Repository struct {
	WbOrder
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		WbOrder: NewOrderPG(db),
	}
}
