package store

import (
	"context"
	"fmt"

	j "github.com/OlesyaNovikova/gophermart/internal/models/json"
)

type TestStore struct {
	us map[string][]byte      //ключ UserName
	or map[string]j.Orders    //ключ Number
	wd map[string]j.Withdraws //ключ Order
}

func NewStore() (TestStore, error) {
	return TestStore{
		us: make(map[string][]byte),
		or: make(map[string]j.Orders),
		wd: make(map[string]j.Withdraws),
	}, nil
}

func (t *TestStore) AddUser(ctx context.Context, userID string, pass []byte) error {
	if _, ok := t.us[userID]; !ok {
		t.us[userID] = pass
		return nil
	}
	return fmt.Errorf("пользователь с таким именем существует")
}

func (t *TestStore) GetPass(ctx context.Context, userID string) ([]byte, error) {
	if pass, ok := t.us[userID]; ok {
		return pass, nil
	}
	return nil, fmt.Errorf("пароль не найден")
}

func (t *TestStore) AddOrder(ctx context.Context, order j.Orders) error {
	if _, ok := t.or[order.Number]; !ok {
		t.or[order.Number] = order
		return nil
	}
	return fmt.Errorf("заказ был добавлен ранее")
}

func (t *TestStore) GetOrder(ctx context.Context, number string) (j.Orders, error) {
	if or, ok := t.or[number]; ok {
		return or, nil
	}
	return j.Orders{}, fmt.Errorf("заказ не найден")
}

func (t *TestStore) GetOrders(ctx context.Context, userID string) ([]j.Orders, error) {
	var userOrders []j.Orders
	for _, or := range t.or {
		if or.UserName == userID {
			userOrders = append(userOrders, or)
		}
	}
	//тут должна быть сортировка от старых к новым, но ее не будет в тестовом варианте
	return userOrders, nil
}

func (t *TestStore) UpdateOrder(ctx context.Context, order j.Orders) error {
	if _, ok := t.or[order.Number]; ok {
		t.or[order.Number] = order
		return nil
	}
	return fmt.Errorf("заказ не найден")
}

func (t *TestStore) GetOrdersForUpd(ctx context.Context) ([]j.Orders, error) {
	var updOrders []j.Orders
	for _, or := range t.or {
		if !(or.Status == j.StatInvalid || or.Status == j.StatProcessed) {
			updOrders = append(updOrders, or)
		}
	}
	return updOrders, nil
}

func (t *TestStore) GetBalance(ctx context.Context, userID string) (j.Balance, error) {
	orders, err := t.GetOrders(ctx, userID)
	if err != nil {
		return j.Balance{}, err
	}
	var balance j.Balance
	for _, order := range orders {
		balance.Current += order.Accrual
	}
	withdraws, err := t.GetWithdraws(ctx, userID)
	if err != nil {
		return j.Balance{}, err
	}
	for _, withdraw := range withdraws {
		balance.Current -= withdraw.Sum
		balance.Withdrawn += withdraw.Sum
	}
	return balance, nil
}

func (t *TestStore) AddWithdraw(ctx context.Context, withdraw j.Withdraws) error {
	if _, ok := t.wd[withdraw.Order]; !ok {
		t.wd[withdraw.Order] = withdraw
		return nil
	}
	return fmt.Errorf("заказ был добавлен ранее")
}

func (t *TestStore) GetWithdraws(ctx context.Context, userID string) ([]j.Withdraws, error) {
	var userWithdraws []j.Withdraws
	for _, wd := range t.wd {
		if wd.UserName == userID {
			userWithdraws = append(userWithdraws, wd)
		}
	}
	//тут должна быть сортировка от старых к новым, но ее не будет в тестовом варианте
	return userWithdraws, nil
}
