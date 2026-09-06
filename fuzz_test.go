package checkdigit

import "testing"

// FuzzClean checks the one invariant Clean promises: it only ever removes
// hyphens and spaces, never anything else, and never grows the input.
func FuzzClean(f *testing.F) {
	f.Add("978-0-306-40615-7")
	f.Add("  036 000 291 452  ")
	f.Add("")

	f.Fuzz(func(t *testing.T, s string) {
		out := Clean(s)
		if len(out) > len(s) {
			t.Fatalf("Clean(%q) grew the input to %q", s, out)
		}
		for i := 0; i < len(out); i++ {
			if out[i] == '-' || out[i] == ' ' {
				t.Fatalf("Clean(%q) left a separator in %q", s, out)
			}
		}
	})
}

// FuzzCheckFunctions exists to catch panics (out-of-range indexing, mainly)
// rather than to assert particular outcomes: arbitrary input should always
// come back as one of the sentinel errors, never crash the process.
func FuzzCheckFunctions(f *testing.F) {
	seeds := []string{
		"0-306-40615-2",
		"978-0-306-40615-7",
		"036000291452",
		"0378-5955",
		"1000002X",
		"CODE39W",
		"",
		"x",
		"X",
		"-----",
		"978030640615X",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		_ = CheckISBN10(s)
		_ = CheckISBN13(s)
		_ = CheckUPCA(s)
		_ = CheckISSN(s)
		_ = CheckCode39(s)
	})
}

// FuzzISBN10RoundTrip only exercises inputs GenerateISBN10 accepts, but that
// still covers the full 9-digit prefix space; a check digit it produces
// should always satisfy CheckISBN10.
func FuzzISBN10RoundTrip(f *testing.F) {
	f.Add("030640615")
	f.Fuzz(func(t *testing.T, prefix string) {
		code, err := GenerateISBN10(prefix)
		if err != nil {
			return
		}
		if err := CheckISBN10(code); err != nil {
			t.Fatalf("GenerateISBN10(%q) = %q, which CheckISBN10 rejects: %v", prefix, code, err)
		}
	})
}

func FuzzISBN13RoundTrip(f *testing.F) {
	f.Add("978030640615")
	f.Fuzz(func(t *testing.T, prefix string) {
		code, err := GenerateISBN13(prefix)
		if err != nil {
			return
		}
		if err := CheckISBN13(code); err != nil {
			t.Fatalf("GenerateISBN13(%q) = %q, which CheckISBN13 rejects: %v", prefix, code, err)
		}
	})
}

func FuzzUPCARoundTrip(f *testing.F) {
	f.Add("03600029145")
	f.Fuzz(func(t *testing.T, prefix string) {
		code, err := GenerateUPCA(prefix)
		if err != nil {
			return
		}
		if err := CheckUPCA(code); err != nil {
			t.Fatalf("GenerateUPCA(%q) = %q, which CheckUPCA rejects: %v", prefix, code, err)
		}
	})
}

func FuzzISSNRoundTrip(f *testing.F) {
	f.Add("0378595")
	f.Add("1000002")
	f.Fuzz(func(t *testing.T, prefix string) {
		code, err := GenerateISSN(prefix)
		if err != nil {
			return
		}
		if err := CheckISSN(code); err != nil {
			t.Fatalf("GenerateISSN(%q) = %q, which CheckISSN rejects: %v", prefix, code, err)
		}
	})
}

func FuzzCode39RoundTrip(f *testing.F) {
	f.Add("CODE39")
	f.Fuzz(func(t *testing.T, data string) {
		code, err := GenerateCode39(data)
		if err != nil {
			return
		}
		if err := CheckCode39(code); err != nil {
			t.Fatalf("GenerateCode39(%q) = %q, which CheckCode39 rejects: %v", data, code, err)
		}
	})
}
