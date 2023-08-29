package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	s "github.com/OlesyaNovikova/gophermart/internal/store"
)

func TestRegister(t *testing.T) {

	// описываем набор данных: метод запроса, ожидаемый код ответа, ожидаемое тело
	testCases := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{name: "Корректные данные", body: `{"login": "vasya", "password": "123" }`, expectedCode: http.StatusOK},
		{name: "Повторные данные", body: `{"login": "vasya", "password": "123" }`, expectedCode: http.StatusConflict},
		{name: "Пустые данные", body: "", expectedCode: http.StatusBadRequest},
		{name: "Не корректные данные", body: `{"login": "petya"}`, expectedCode: http.StatusBadRequest},
	}

	store, _ := s.NewStore()
	InitStore(&store)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(tc.body))
			w := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/api/user/register", Register())
			router.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")

		})
	}
}
