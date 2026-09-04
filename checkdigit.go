// Package checkdigit computes and verifies the check digits used by common
// book and retail barcode formats: ISBN-10, ISBN-13 (which shares its
// algorithm with EAN-13), UPC-A, ISSN, and Code 39.
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

// GenerateISBN10 appends the check character to a 9-digit ISBN-10 prefix,
// returning the full 10-character code.
func GenerateISBN10(prefix string) (string, error) {
	check, err := ISBN10CheckDigit(prefix)
	if err != nil {
		return "", err
	}
	return prefix + string(check), nil
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

// GenerateISBN13 appends the check digit to a 12-digit ISBN-13 (or EAN-13)
// prefix, returning the full 13-digit code.
func GenerateISBN13(prefix string) (string, error) {
	check, err := ISBN13CheckDigit(prefix)
	if err != nil {
		return "", err
	}
	return prefix + string(check), nil
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

// GenerateUPCA appends the check digit to an 11-digit UPC-A prefix,
// returning the full 12-digit code.
func GenerateUPCA(prefix string) (string, error) {
	check, err := UPCACheckDigit(prefix)
	if err != nil {
		return "", err
	}
	return prefix + string(check), nil
}

// ISSNCheckDigit computes the check character for a 7-digit ISSN prefix
// using descending weights 8 down to 2. The result is '0'-'9' or 'X' (for a
// check value of 10).
func ISSNCheckDigit(prefix string) (byte, error) {
	if len(prefix) != 7 {
		return 0, ErrInvalidLength
	}
	sum := 0
	for i := 0; i < 7; i++ {
		c := prefix[i]
		if c < '0' || c > '9' {
			return 0, ErrInvalidDigit
		}
		sum += int(c-'0') * (8 - i)
	}
	check := 11 - sum%11
	switch check {
	case 11:
		return '0', nil
	case 10:
		return 'X', nil
	default:
		return byte('0' + check), nil
	}
}

// CheckISSN reports whether s (after removing hyphens and spaces) is a
// valid 8-digit ISSN, including the trailing check character.
func CheckISSN(s string) error {
	digits := Clean(s)
	if len(digits) != 8 {
		return ErrInvalidLength
	}
	want, err := ISSNCheckDigit(digits[:7])
	if err != nil {
		return err
	}
	got := digits[7]
	if got == 'x' {
		got = 'X'
	}
	if got != 'X' && (got < '0' || got > '9') {
		return ErrInvalidDigit
	}
	if got != want {
		return ErrChecksumMismatch
	}
	return nil
}

// ValidateISSN is a convenience wrapper around CheckISSN.
func ValidateISSN(s string) bool {
	return CheckISSN(s) == nil
}

// GenerateISSN appends the check character to a 7-digit ISSN prefix,
// returning the full 8-character code.
func GenerateISSN(prefix string) (string, error) {
	check, err := ISSNCheckDigit(prefix)
	if err != nil {
		return "", err
	}
	return prefix + string(check), nil
}

// code39Value returns the position (0-42) of c in the Code 39 character
// set, used both as the character's point value and, run in reverse, to
// turn a computed check value back into a character.
func code39Value(c byte) (int, bool) {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0'), true
	case c >= 'A' && c <= 'Z':
		return int(c-'A') + 10, true
	}
	switch c {
	case '-':
		return 36, true
	case '.':
		return 37, true
	case ' ':
		return 38, true
	case '$':
		return 39, true
	case '/':
		return 40, true
	case '+':
		return 41, true
	case '%':
		return 42, true
	default:
		return 0, false
	}
}

func code39Char(v int) byte {
	switch {
	case v < 10:
		return byte('0' + v)
	case v < 36:
		return byte('A' + v - 10)
	}
	switch v {
	case 36:
		return '-'
	case 37:
		return '.'
	case 38:
		return ' '
	case 39:
		return '$'
	case 40:
		return '/'
	case 41:
		return '+'
	default:
		return '%'
	}
}

// Code39CheckChar computes the mod-43 check character for data, the encoded
// content of a Code 39 barcode between (but not including) its start and
// stop '*' delimiters.
func Code39CheckChar(data string) (byte, error) {
	if data == "" {
		return 0, ErrInvalidLength
	}
	sum := 0
	for i := 0; i < len(data); i++ {
		v, ok := code39Value(data[i])
		if !ok {
			return 0, ErrInvalidDigit
		}
		sum += v
	}
	return code39Char(sum % 43), nil
}

// CheckCode39 reports whether s is Code 39 data with a valid trailing
// mod-43 check character. s should not include the start/stop '*'
// delimiters; strip those before calling.
func CheckCode39(s string) error {
	if len(s) < 2 {
		return ErrInvalidLength
	}
	want, err := Code39CheckChar(s[:len(s)-1])
	if err != nil {
		return err
	}
	got := s[len(s)-1]
	if _, ok := code39Value(got); !ok {
		return ErrInvalidDigit
	}
	if got != want {
		return ErrChecksumMismatch
	}
	return nil
}

// ValidateCode39 is a convenience wrapper around CheckCode39.
func ValidateCode39(s string) bool {
	return CheckCode39(s) == nil
}

// GenerateCode39 appends the mod-43 check character to Code 39 data,
// returning the full string (still without the start/stop '*' delimiters).
func GenerateCode39(data string) (string, error) {
	check, err := Code39CheckChar(data)
	if err != nil {
		return "", err
	}
	return data + string(check), nil
}
