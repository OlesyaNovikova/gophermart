package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"net/http"

	j "github.com/OlesyaNovikova/gophermart/internal/models/json"
	a "github.com/OlesyaNovikova/gophermart/internal/utils/auth"
)

func Login() http.HandlerFunc {
	fn := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		var name j.User
		var inBuf bytes.Buffer

		_, err := inBuf.ReadFrom(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(inBuf.Bytes(), &name); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		basePass, err := store.s.GetPass(ctx, name.UserName)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		h := sha256.New()
		h.Write([]byte(name.Pass))
		pass := h.Sum(nil)

		if !bytes.Equal(basePass, pass) {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := a.BuildJWTString(name.UserName)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		cookie := &http.Cookie{
			Name:  "authToken",
			Value: token,
		}
		http.SetCookie(res, cookie)
		res.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
}
