package util

import (
	"strconv"
	"time"
	"unicode"
)

// IsValidPAN checks if the given PAN (Primary Account Number) is valid.
// It uses the Luhn algorithm to validate the number.
// The PAN should consist of digits only and can be of any length.
// The function returns true if the PAN is valid, otherwise false.
func IsValidPAN(pan string) bool {
	var sum int
	var alt bool

	// Process digits from right to left
	for i := len(pan) - 1; i >= 0; i-- {
		r := rune(pan[i])
		if !unicode.IsDigit(r) {
			continue // skip non-digit characters like spaces or dashes
		}

		digit := int(r - '0')
		if alt {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alt = !alt
	}

	return sum%10 == 0
}

// IsValidExpiry checks if the given expiry date is valid.
// The expiry date should be in the format "YYYYMM" and should not be in the past.
// It also checks if the month is between 1 and 12.
// The function returns true if the expiry date is valid, otherwise false.
func IsValidExpiry(expiry string) bool {
	if len(expiry) != 6 {
		return false
	}

	year, err1 := strconv.Atoi(expiry[:4])
	month, err2 := strconv.Atoi(expiry[4:])
	if err1 != nil || err2 != nil || month < 1 || month > 12 {
		return false
	}

	now := time.Now()
	expiryTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	// Move to last day of expiry month
	expiryTime = expiryTime.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return expiryTime.After(now)
}
