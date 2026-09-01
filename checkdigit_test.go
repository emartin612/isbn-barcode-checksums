package checkdigit

import (
	"strings"
	"testing"
)

func TestISBN10(t *testing.T) {
	// Wikipedia's ISBN-10 worked example.
	if !ValidateISBN10("0-306-40615-2") {
		t.Error("expected valid ISBN-10")
	}
	if ValidateISBN10("0-306-40615-3") {
		t.Error("expected invalid ISBN-10 (wrong check digit)")
	}
	if ValidateISBN10("0-306-4061") {
		t.Error("expected invalid ISBN-10 (too short)")
	}
}

func TestISBN10CheckDigitX(t *testing.T) {
	// 0-8044-2957-X is a well known example with an 'X' check character.
	if !ValidateISBN10("0-8044-2957-X") {
		t.Error("expected valid ISBN-10 with X check digit")
	}
}

func TestISBN13(t *testing.T) {
	if !ValidateISBN13("978-0-306-40615-7") {
		t.Error("expected valid ISBN-13")
	}
	if ValidateISBN13("978-0-306-40615-8") {
		t.Error("expected invalid ISBN-13 (wrong check digit)")
	}
}

func TestUPCA(t *testing.T) {
	if !ValidateUPCA("036000291452") {
		t.Error("expected valid UPC-A")
	}
	if ValidateUPCA("036000291453") {
		t.Error("expected invalid UPC-A (wrong check digit)")
	}
}

func TestISSN(t *testing.T) {
	if !ValidateISSN("0378-5955") {
		t.Error("expected valid ISSN")
	}
	if ValidateISSN("0378-5954") {
		t.Error("expected invalid ISSN (wrong check digit)")
	}
	if ValidateISSN("0378-595") {
		t.Error("expected invalid ISSN (too short)")
	}
}

func TestISSNCheckDigitX(t *testing.T) {
	// Prefixes that land on a check value of 10 use 'X', the same
	// convention as ISBN-10.
	digit, err := ISSNCheckDigit("1000002")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if digit != 'X' {
		t.Errorf("got check digit %q, want 'X'", digit)
	}
	if !ValidateISSN("1000002X") {
		t.Error("expected valid ISSN with X check digit")
	}
}

func TestCode39(t *testing.T) {
	digit, err := Code39CheckChar("CODE39")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if digit != 'W' {
		t.Errorf("got check char %q, want 'W'", digit)
	}
	if !ValidateCode39("CODE39W") {
		t.Error("expected valid Code 39 data")
	}
	if ValidateCode39("CODE39X") {
		t.Error("expected invalid Code 39 data (wrong check char)")
	}
	if ValidateCode39("") {
		t.Error("expected invalid Code 39 data (empty)")
	}
	if ValidateCode39("ABC?1") {
		t.Error("expected invalid Code 39 data (character outside the set)")
	}
}

func TestValidateStream(t *testing.T) {
	input := strings.Join([]string{
		"0-306-40615-2",  // valid ISBN-10
		"978-0-306-40615-7", // valid ISBN-13
		"036000291452",   // valid UPC-A
		"1234567890",     // invalid ISBN-10 (check digit should be X)
		"",               // blank lines are skipped
		"not-a-barcode",  // wrong length entirely
	}, "\n")

	var results []Result
	err := ValidateStream(strings.NewReader(input), func(r Result) error {
		results = append(results, r)
		return nil
	})
	if err != nil {
		t.Fatalf("ValidateStream returned error: %v", err)
	}

	if len(results) != 5 {
		t.Fatalf("expected 5 non-blank results, got %d", len(results))
	}
	if results[0].Kind != KindISBN10 || results[0].Err != nil {
		t.Errorf("line 1: got kind %v err %v", results[0].Kind, results[0].Err)
	}
	if results[3].Err != ErrChecksumMismatch {
		t.Errorf("line 4: expected checksum mismatch, got %v", results[3].Err)
	}
	if results[4].Kind != KindUnknown {
		t.Errorf("line 6: expected unknown kind, got %v", results[4].Kind)
	}
}
