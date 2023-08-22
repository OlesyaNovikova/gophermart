package handlers

import (
	"context"

	j "github.com/OlesyaNovikova/gophermart/internal/models/json"
)

type Storage interface {
	AddUser(ctx context.Context, userID string, pass []byte) error
	GetPass(ctx context.Context, userID string) ([]byte, error)
	AddOrder(ctx context.Context, order j.Orders) error
	GetOrder(ctx context.Context, number string) (j.Orders, error)
	GetOrders(ctx context.Context, userID string) ([]j.Orders, error)
	GetBalance(ctx context.Context, userID string) (j.Balance, error)
	AddWithdraw(ctx context.Context, withdraw j.Withdraws) error
	GetWithdraws(ctx context.Context, userID string) ([]j.Withdraws, error)
}

type storage struct {
	s Storage
}

var store storage

func InitStore(s Storage) {
	store.s = s
}
