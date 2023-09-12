package luhn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLuhn(t *testing.T) {

	testCases := []struct {
		name    string
		number  string
		expExit bool
	}{
		{name: "Соответствующий номер", number: "12345678903", expExit: true},
		{name: "Не соответствующий номер", number: "12345678905", expExit: false},
		{name: "Не корректный номер", number: "12a4567i903", expExit: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expExit, Luhn(tc.number), "Результат проверки не совпадает с ожидаемым")
		})
	}
}
