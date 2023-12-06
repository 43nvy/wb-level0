package service

import (
	"github.com/43nvy/wb_l0"
	"github.com/43nvy/wb_l0/internal/repository"
)

type OrderService struct {
	repo repository.WbOrder
}

func NewOrderService(repo repository.WbOrder) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(order wb_l0.Order) (int, error) {
	return s.repo.CreateOrder(order)
}

func (s *OrderService) GetOrder(id int) (wb_l0.Order, error) {
	return s.repo.GetOrder(id)
}

func (s *OrderService) GetTenOrders() ([]wb_l0.Order, error) {
	return s.repo.GetTenOrders()
}
