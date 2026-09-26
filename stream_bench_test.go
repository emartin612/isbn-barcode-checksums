package checkdigit

import (
	"strconv"
	"strings"
	"testing"
)

// buildStreamInput generates n lines of valid codes, cycling through the
// four formats ValidateStream dispatches on by digit count, so the
// benchmark exercises every branch rather than just whichever format
// happens to be first.
func buildStreamInput(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		switch i % 4 {
		case 0:
			b.WriteString("0-306-40615-2\n")
		case 1:
			b.WriteString("978-0-306-40615-7\n")
		case 2:
			b.WriteString("036000291452\n")
		case 3:
			b.WriteString("0378-5955\n")
		}
	}
	return b.String()
}

// BenchmarkValidateStream checks that memory use stays flat as the input
// grows, which is the whole point of scanning line by line instead of
// slurping the file. Run with -benchmem and compare B/op across sizes: it
// should stay constant, not scale with n.
func BenchmarkValidateStream(b *testing.B) {
	for _, n := range []int{100, 10_000, 1_000_000} {
		input := buildStreamInput(n)
		b.Run(strconv.Itoa(n)+"_lines", func(b *testing.B) {
			b.SetBytes(int64(len(input)))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r := strings.NewReader(input)
				if err := ValidateStream(r, func(Result) error { return nil }); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
