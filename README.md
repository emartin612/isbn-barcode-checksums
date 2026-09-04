# isbn-barcode-checksums

Go library for computing and verifying the check digits used by book and
retail barcodes: ISBN-10, ISBN-13/EAN-13, UPC-A, ISSN, and Code 39.

Every one of these formats ends in a digit that's derived from the ones
before it, so a single mistyped or misscanned character gets caught before
it turns into a bad lookup, a mis-shelved book, or a receipt with the wrong
item on it. This library implements the three algorithms and a way to run
them over a stream of codes without loading the whole input into memory.

No third-party dependencies - standard library only.

## Install

```
go get github.com/emartin612/isbn-barcode-checksums
```

## Usage

Validate a single code:

```go
package main

import (
	"fmt"

	"github.com/emartin612/isbn-barcode-checksums"
)

func main() {
	fmt.Println(checkdigit.ValidateISBN13("978-0-306-40615-7")) // true
	fmt.Println(checkdigit.ValidateUPCA("036000291452"))         // true

	if err := checkdigit.CheckISBN10("0-306-40615-3"); err != nil {
		fmt.Println("bad ISBN-10:", err) // checksum mismatch
	}
}
```

Compute a check digit for a code you're generating yourself:

```go
digit, err := checkdigit.ISBN13CheckDigit("978030640615")
// digit == '7'
```

Or get the full code back in one call, prefix and check digit together:

```go
code, err := checkdigit.GenerateISBN13("978030640615")
// code == "9780306406157"
```

Every format has a matching `Generate*` function (`GenerateISBN10`,
`GenerateISBN13`, `GenerateUPCA`, `GenerateISSN`, `GenerateCode39`), each
taking the same prefix length its `*CheckDigit` counterpart expects.

Check a whole file of codes, one per line, without reading it into memory
first:

```go
f, err := os.Open("codes.txt")
if err != nil {
	log.Fatal(err)
}
defer f.Close()

err = checkdigit.ValidateStream(f, func(r checkdigit.Result) error {
	if r.Err != nil {
		fmt.Printf("line %d: %s (%s): %v\n", r.Line, r.Raw, r.Kind, r.Err)
	}
	return nil // returning a non-nil error stops the scan early
})
if err != nil {
	log.Fatal(err)
}
```

`ValidateStream` reads line by line with a `bufio.Scanner`, so memory use
stays flat regardless of how many codes the file contains - it's the same
whether you feed it a hundred lines or a hundred million. It picks a format
per line from the cleaned digit count (8 -> ISSN, 10 -> ISBN-10, 12 ->
UPC-A, 13 -> ISBN-13), so it doesn't cover Code 39, whose data isn't purely
numeric.

Code 39 is checked separately, since its check character depends on a
43-character alphabet rather than plain digits:

```go
digit, err := checkdigit.Code39CheckChar("CODE39") // 'W', nil
fmt.Println(checkdigit.ValidateCode39("CODE39W"))  // true
```

## Formats supported

- ISBN-10 (mod 11, trailing check character can be `X`)
- ISBN-13 / EAN-13 (mod 10, alternating 1/3 weights)
- UPC-A (mod 10, alternating 3/1 weights)
- ISSN (mod 11, descending weights 8-2, trailing check character can be `X`)
- Code 39 (mod 43 over the 43-character Code 39 alphabet)

## License

MIT, see LICENSE.
