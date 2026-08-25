# isbn-barcode-checksums

Go library for computing and verifying the check digits used by book and
retail barcodes: ISBN-10, ISBN-13/EAN-13, and UPC-A.

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
whether you feed it a hundred lines or a hundred million.

## Formats supported

- ISBN-10 (mod 11, trailing check character can be `X`)
- ISBN-13 / EAN-13 (mod 10, alternating 1/3 weights)
- UPC-A (mod 10, alternating 3/1 weights)

## License

MIT, see LICENSE.
