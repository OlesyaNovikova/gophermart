package accruals

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

/*if or.Status == j.StatNew {
	or.Status = j.StatRegistered
}
if order.Status == j.StatRegistered {
	order.Status = j.StatNew
}*/

func InitAccruals(addr string, s Storage) {
	accrualAddr = addr
	store.s = s
	client = &http.Client{}
}

func GetAccrual(ctx context.Context, number string) (j.Accrual, error) {
	addr := "http://" + accrualAddr + "/api/orders/" + number
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
	return j.Accrual{}, fmt.Errorf("информация не получена, статус %v", resp.StatusCode)
}
