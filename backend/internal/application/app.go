package application

import (
	"unicode"
)

func ValidExpression(expression string) bool {
	for _, r := range expression {
		if !unicode.Is(unicode.Digit, r) {
			/*проверяем наличие любых символов кроме ( ) + - / *  */
			if !(r >= 0x08 && r <= 0x0B || r == 0x0D || r == 0x0F) {
				return false
			}
		}
	}
	return true
}
