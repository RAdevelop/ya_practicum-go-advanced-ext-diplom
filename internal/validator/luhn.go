package validator

// IsValidLuhn проверяет, что строка из цифр проходит алгоритм Луна.
// Возвращает false, если строка пуста, содержит не цифровые символы
// или не проходит проверку.
func IsValidLuhn(number string) bool {
	if number == "" {
		return false
	}

	sum := 0
	double := false

	// Идём справа налево.
	for i := len(number) - 1; i >= 0; i-- {
		c := number[i]

		// Только цифры.
		if c < '0' || c > '9' {
			return false
		}

		digit := int(c - '0')

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		double = !double
	}

	return sum%10 == 0
}
