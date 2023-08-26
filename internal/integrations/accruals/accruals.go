package accruals

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	j "github.com/OlesyaNovikova/gophermart/internal/models/json"
)

type Storage interface {
	UpdateOrder(ctx context.Context, order j.Orders) error
	GetOrdersForUpd(ctx context.Context) ([]j.Orders, error)
}

type storage struct {
	s Storage
}

var accrualAddr string
var store storage
var client *http.Client

const updInt = time.Second * 10

func InitAccruals(ctx context.Context, addr string, s Storage) chan j.Orders {
	accrualAddr = addr + "/api/orders/"
	store.s = s
	client = &http.Client{}
	order := make(chan j.Orders)
	go AccrualRout(ctx, order)
	return order
}

func AccrualRout(ctx context.Context, or chan j.Orders) {
	orders, err := store.s.GetOrdersForUpd(ctx)
	if err == nil {
		for _, order := range orders {
			o := order
			go UpdAccrual(ctx, o)
		}
	}
	for {
		select {
		case <-ctx.Done():
			return
		case o, ok := <-or:
			if !ok {
				return
			}
			go UpdAccrual(ctx, o)
		}
	}
}

func UpdAccrual(ctx context.Context, order j.Orders) {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ticker.Reset(updInt)
			accrual, err := GetAccrual(ctx, order.Number)
			if err != nil {
				break
			}
			status := accrual.Status
			if status == j.StatRegistered {
				status = j.StatNew
			}
			if status != order.Status {
				date := time.Now()
				dateStr := date.Format(time.RFC3339)

				err = store.s.UpdateOrder(ctx, j.Orders{
					UserName: order.UserName,
					Number:   order.Number,
					Status:   status,
					Accrual:  accrual.Accrual,
					DateStr:  dateStr,
					Date:     date,
				})
				if err == nil {
					if status == j.StatInvalid || status == j.StatProcessed {
						return
					}
					order.Status = status
				}
			}
		}
	}
}

func GetAccrual(ctx context.Context, number string) (j.Accrual, error) {
	addr := accrualAddr + number
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
	if err != nil {
		return j.Accrual{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return j.Accrual{}, err
	}
	defer resp.Body.Close()

	stat := resp.StatusCode
	if stat == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return j.Accrual{}, err
		}
		var accrual j.Accrual
		err = json.Unmarshal(body, &accrual)
		if err != nil {
			return j.Accrual{}, err
		}
		return accrual, nil
	}
	/*if stat == http.StatusNoContent {
		return j.Accrual{Number: number, Status: j.StatInvalid}, nil
	}*/
	return j.Accrual{}, fmt.Errorf("информация не получена, статус %v", resp.StatusCode)
}
