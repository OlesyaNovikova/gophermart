package luhn

import "strconv"

// Luhn проверяет число на соответствие по алгоритму Луна
func Luhn(number string) bool {
	// Преобразуем строку в массив цифр
	digits := make([]int, len(number))
	for i, char := range number {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			// Возникла ошибка при преобразовании символа в число
			return false
		}
		digits[i] = digit
	}

	// Выполняем алгоритм Луна
	sum := 0
	for i := len(digits) - 2; i >= 0; i -= 2 {
		digits[i] = digits[i] * 2
		if digits[i] >= 10 {
			digits[i] = digits[i] - 9
		}
	}
	for _, digit := range digits {
		sum += digit
	}

	// Проверяем, что сумма делится нацело на 10
	if sum%10 == 0 {
		return true
	} else {
		return false
	}
}
