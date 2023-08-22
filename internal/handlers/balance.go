package handlers

import (
	"encoding/json"
	"net/http"

	a "github.com/OlesyaNovikova/gophermart/internal/utils/auth"
)

func Balance() http.HandlerFunc {
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
		balance, err := store.s.GetBalance(ctx, name)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		body, err := json.Marshal(balance)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(body)
	}
	return http.HandlerFunc(fn)
}
