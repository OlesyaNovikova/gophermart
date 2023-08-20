package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	ac "github.com/OlesyaNovikova/gophermart/internal/integrations/accruals"
	j "github.com/OlesyaNovikova/gophermart/internal/models/json"
	l "github.com/OlesyaNovikova/gophermart/internal/utils/luhn"
)

func OrdersPost(ch chan j.Orders) http.HandlerFunc {
	fn := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		cookie, err := req.Cookie("authToken")
		if err != nil || cookie.Value == "" {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}
		name := cookie.Value

		var inBuf bytes.Buffer
		_, err = inBuf.ReadFrom(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		number := inBuf.String()

		if !l.Luhn(number) {
			http.Error(res, "Invalid order number", http.StatusUnprocessableEntity)
			return
		}

		ord, err := store.s.GetOrder(ctx, number)
		if err == nil {
			if ord.UserName == name {
				res.WriteHeader(http.StatusOK)
				return
			}
			http.Error(res, "The order is registered to another user", http.StatusConflict)
			return
		}

		accrual, err := ac.GetAccrual(ctx, number)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		status := accrual.Status
		if status == j.StatRegistered {
			status = j.StatNew
		}
		date := time.Now()
		dateStr := date.Format(time.RFC3339)

		order := j.Orders{
			UserName: name,
			Number:   number,
			Status:   status,
			Accrual:  accrual.Accrual,
			DateStr:  dateStr,
			Date:     date,
		}

		err = store.s.AddOrder(ctx, order)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if status != j.StatInvalid && status != j.StatProcessed {
			ch <- order
		}

		res.WriteHeader(http.StatusAccepted)
	}
	return http.HandlerFunc(fn)
}
