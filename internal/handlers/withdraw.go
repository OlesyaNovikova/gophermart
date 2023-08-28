package handlers

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"sync"

	j "github.com/OlesyaNovikova/gophermart/internal/models/json"
	a "github.com/OlesyaNovikova/gophermart/internal/utils/auth"
	l "github.com/OlesyaNovikova/gophermart/internal/utils/luhn"
)

var mut sync.Mutex

func Withdraw() http.HandlerFunc {
	fn := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		cookie, err := req.Cookie("authToken")
		if err != nil || cookie.Value == "" {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}
		name, err := a.GetUserID(cookie.Value)
		if err != nil {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var withdraw j.Withdraws
		var inBuf bytes.Buffer
		_, err = inBuf.ReadFrom(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(inBuf.Bytes(), &withdraw); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if !l.Luhn(withdraw.Order) {
			http.Error(res, "Invalid order number", http.StatusUnprocessableEntity)
			return
		}
		withdraw.Sum = math.Ceil(withdraw.Sum*100) / 100

		mut.Lock()
		balance, err := store.s.GetBalance(ctx, name)
		if err != nil {
			mut.Unlock()
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if balance.Current < withdraw.Sum {
			mut.Unlock()
			http.Error(res, "Not enough bonus points", http.StatusPaymentRequired)
			return
		}
		err = store.s.AddWithdraw(ctx, withdraw)
		mut.Unlock()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
}
