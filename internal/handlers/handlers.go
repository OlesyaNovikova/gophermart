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
}

type storage struct {
	s Storage
}

var store storage

func InitStore(s Storage) {
	store.s = s
}
