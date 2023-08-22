package handlers

import (
	"encoding/json"
	"net/http"

	a "github.com/OlesyaNovikova/gophermart/internal/utils/auth"
)

func Withdrawals() http.HandlerFunc {
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
		withdraws, err := store.s.GetWithdraws(ctx, name)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(withdraws) == 0 {
			res.WriteHeader(http.StatusNoContent)
			return
		}
		body, err := json.Marshal(withdraws)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(body)
	}
	return http.HandlerFunc(fn)
}
