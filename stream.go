package checkdigit

import (
	"bufio"
	"io"
	"strings"
)

// Kind identifies which barcode format a scanned line looked like, based on
// its digit count after cleaning.
type Kind int

const (
	KindUnknown Kind = iota
	KindISBN10
	KindISBN13
	KindUPCA
	KindISSN
)

func (k Kind) String() string {
	switch k {
	case KindISBN10:
		return "ISBN-10"
	case KindISBN13:
		return "ISBN-13"
	case KindUPCA:
		return "UPC-A"
	case KindISSN:
		return "ISSN"
	default:
		return "unknown"
	}
}

// Result describes the outcome of checking one line from a stream.
type Result struct {
	Line int    // 1-based line number in the input
	Raw  string // the line as read, before cleaning
	Kind Kind
	Err  error // nil if the check digit is valid
}

// ValidateStream reads newline-separated codes from r, one at a time, and
// calls fn with the result of checking each non-blank line. The format is
// picked per line from its cleaned length (8 digits -> ISSN, 10 -> ISBN-10,
// 12 -> UPC-A, 13 -> ISBN-13); anything else is reported as KindUnknown
// with ErrInvalidLength. Code 39 isn't included here because its data
// isn't purely numeric, so length alone can't identify it.
//
// Only the current line is ever held in memory, so this is safe to run
// against a barcode list of any size without buffering the whole thing:
// callers that want a running tally or a filtered list should accumulate it
// themselves inside fn.
//
// ValidateStream stops and returns the first error fn returns. It also
// returns any error encountered while reading r.
func ValidateStream(r io.Reader, fn func(Result) error) error {
	scanner := bufio.NewScanner(r)
	line := 0
	for scanner.Scan() {
		line++
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" {
			continue
		}

		digits := Clean(raw)
		res := Result{Line: line, Raw: raw}

		switch len(digits) {
		case 8:
			res.Kind = KindISSN
			res.Err = CheckISSN(digits)
		case 10:
			res.Kind = KindISBN10
			res.Err = CheckISBN10(digits)
		case 12:
			res.Kind = KindUPCA
			res.Err = CheckUPCA(digits)
		case 13:
			res.Kind = KindISBN13
			res.Err = CheckISBN13(digits)
		default:
			res.Kind = KindUnknown
			res.Err = ErrInvalidLength
		}

		if err := fn(res); err != nil {
			return err
		}
	}
	return scanner.Err()
}
