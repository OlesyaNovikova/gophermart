package handlers

import (
	"encoding/json"
	"net/http"
)

func OrdersGet() http.HandlerFunc {
	fn := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		cookie, err := req.Cookie("authToken")
		if err != nil || cookie.Value == "" {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}
		name := cookie.Value
		orders, err := store.s.GetOrders(ctx, name)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(orders) == 0 {
			res.WriteHeader(http.StatusNoContent)
			return
		}
		body, err := json.Marshal(orders)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(body)
	}
	return http.HandlerFunc(fn)
}
