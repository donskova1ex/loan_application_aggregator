package validators

import (
	"app_aggregator/internal"
	"regexp"
	"strings"
)

func ValidPhone(phoneNumber string) bool {
	cleanedPhone := regexp.MustCompile(`\D`).ReplaceAllString(phoneNumber, "")

	pattern := `^(?:\+7|8|7)\s?9\d{2}\s?\d{3}\s?\d{2}\s?\d{2}$`
	re := regexp.MustCompile(pattern)

	return re.MatchString(cleanedPhone)
}

func PhoneNormalization(phoneNumber string) (string, error) {
	re := regexp.MustCompile(`\D`)
	digits := re.ReplaceAllString(phoneNumber, "")

	if digits == "" {
		return "", internal.ErrEmptyPhoneNumber
	}
	switch {
	case strings.HasPrefix(digits, "7") && len(digits) == 11:
		return digits, nil
	case strings.HasPrefix(digits, "8") && len(digits) == 11:
		return "7" + digits[1:], nil
	case strings.HasPrefix(digits, "9") && len(digits) == 10:
		return "7" + digits, nil
	}

	return "", internal.ErrInvalidPhoneNumber

}
