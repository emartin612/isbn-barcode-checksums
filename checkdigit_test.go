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
