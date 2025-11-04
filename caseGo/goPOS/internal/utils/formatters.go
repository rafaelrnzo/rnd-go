package utils

import (
	"fmt"
	"strings"
)

func FormatRupiah(n int) string {
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	s := fmt.Sprintf("%d", n)
	var b strings.Builder
	for i, c := range s {
		if i != 0 && (len(s)-i)%3 == 0 {
			b.WriteRune('.')
		}
		b.WriteRune(c)
	}
	return "Rp" + sign + b.String()
}
