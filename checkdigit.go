// Package checkdigit computes and verifies the check digits used by common
// book and retail barcode formats: ISBN-10, ISBN-13 (which shares its
// algorithm with EAN-13), and UPC-A.
package checkdigit

import "errors"

var (
	// ErrInvalidLength is returned when the cleaned input does not have the
	// number of digits the format requires.
	ErrInvalidLength = errors.New("checkdigit: invalid length")

	// ErrInvalidDigit is returned when the input contains a byte that isn't
	// a decimal digit (or, for ISBN-10, the trailing 'X').
	ErrInvalidDigit = errors.New("checkdigit: invalid digit")

	// ErrChecksumMismatch is returned when every digit is well formed but the
	// check digit doesn't match what the algorithm computes.
	ErrChecksumMismatch = errors.New("checkdigit: checksum mismatch")
)

// Clean strips hyphens and spaces, the two separators publishers and
// packaging actually use when printing these numbers. It does not validate
// the remaining characters; Check* functions do that.
func Clean(s string) string {
	b := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '-' || c == ' ' {
			continue
		}
		b = append(b, c)
	}
	return string(b)
}

// ISBN10CheckDigit computes the check character for a 9-digit ISBN-10
// prefix. The result is '0'-'9' or 'X' (for a check value of 10).
func ISBN10CheckDigit(prefix string) (byte, error) {
	if len(prefix) != 9 {
		return 0, ErrInvalidLength
	}
	sum := 0
	for i := 0; i < 9; i++ {
		c := prefix[i]
		if c < '0' || c > '9' {
			return 0, ErrInvalidDigit
		}
		sum += int(c-'0') * (10 - i)
	}
	check := 11 - sum%11
	switch {
	case check == 11:
		return '0', nil
	case check == 10:
		return 'X', nil
	default:
		return byte('0' + check), nil
	}
}

// CheckISBN10 reports whether s (after removing hyphens and spaces) is a
// valid 10-digit ISBN, including the trailing check character.
func CheckISBN10(s string) error {
	digits := Clean(s)
	if len(digits) != 10 {
		return ErrInvalidLength
	}
	want, err := ISBN10CheckDigit(digits[:9])
	if err != nil {
		return err
	}
	got := digits[9]
	if got == 'x' {
		got = 'X'
	}
	if got != '0' && got != 'X' && (got < '0' || got > '9') {
		return ErrInvalidDigit
	}
	if got != want {
		return ErrChecksumMismatch
	}
	return nil
}

// ValidateISBN10 is a convenience wrapper around CheckISBN10 for callers who
// only care whether the number is valid, not why it failed.
func ValidateISBN10(s string) bool {
	return CheckISBN10(s) == nil
}

// mod10CheckDigit implements the check-digit algorithm shared by EAN-13
// (and therefore ISBN-13) and UPC-A: alternating weights of 1 and 3 (or 3
// and 1, depending on where the format starts counting), then rounding the
// sum up to the next multiple of ten.
func mod10CheckDigit(digits string, firstWeight int) byte {
	sum := 0
	w := firstWeight
	for i := 0; i < len(digits); i++ {
		sum += int(digits[i]-'0') * w
		if w == 1 {
			w = 3
		} else {
			w = 1
		}
	}
	check := (10 - sum%10) % 10
	return byte('0' + check)
}

func allDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// ISBN13CheckDigit computes the check digit for a 12-digit ISBN-13 (or
// EAN-13) prefix.
func ISBN13CheckDigit(prefix string) (byte, error) {
	if len(prefix) != 12 {
		return 0, ErrInvalidLength
	}
	if !allDigits(prefix) {
		return 0, ErrInvalidDigit
	}
	return mod10CheckDigit(prefix, 1), nil
}

// CheckISBN13 reports whether s is a valid 13-digit ISBN/EAN-13 number.
func CheckISBN13(s string) error {
	digits := Clean(s)
	if len(digits) != 13 {
		return ErrInvalidLength
	}
	want, err := ISBN13CheckDigit(digits[:12])
	if err != nil {
		return err
	}
	if !allDigits(digits[12:]) {
		return ErrInvalidDigit
	}
	if digits[12] != want {
		return ErrChecksumMismatch
	}
	return nil
}

// ValidateISBN13 is a convenience wrapper around CheckISBN13.
func ValidateISBN13(s string) bool {
	return CheckISBN13(s) == nil
}

// UPCACheckDigit computes the check digit for an 11-digit UPC-A prefix.
func UPCACheckDigit(prefix string) (byte, error) {
	if len(prefix) != 11 {
		return 0, ErrInvalidLength
	}
	if !allDigits(prefix) {
		return 0, ErrInvalidDigit
	}
	return mod10CheckDigit(prefix, 3), nil
}

// CheckUPCA reports whether s is a valid 12-digit UPC-A number.
func CheckUPCA(s string) error {
	digits := Clean(s)
	if len(digits) != 12 {
		return ErrInvalidLength
	}
	want, err := UPCACheckDigit(digits[:11])
	if err != nil {
		return err
	}
	if !allDigits(digits[11:]) {
		return ErrInvalidDigit
	}
	if digits[11] != want {
		return ErrChecksumMismatch
	}
	return nil
}

// ValidateUPCA is a convenience wrapper around CheckUPCA.
func ValidateUPCA(s string) bool {
	return CheckUPCA(s) == nil
}
