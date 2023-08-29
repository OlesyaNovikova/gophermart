package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {

	testCases := []struct {
		name     string
		key      string
		userName string
		err      error
	}{
		{name: "Пример 1", key: "default", userName: "Иван", err: nil},
		{name: "Пример 2", key: "jgflkkjjjjjjjfffyjjklkkkkkkjllllkkujjd", userName: "Василий", err: nil},
		{name: "Пример 3", key: "", userName: "Пётр", err: nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			InitAuth(tc.key)
			token, err := BuildJWTString(tc.userName)
			assert.Equal(t, tc.err, err, "Ошибка получения токена")
			name, err := GetUserID(token)
			assert.Equal(t, tc.err, err, "Ошибка распаковки токена")
			assert.Equal(t, tc.userName, name, "Не совпали вход и выход")
		})
	}
}
