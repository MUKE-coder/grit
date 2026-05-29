package handlers

import (
	"fmt"
	"strconv"
	"strings"
)

// formatMoney renders a UGX amount with thousands separators (no currency
// prefix). Used for human-readable audit-log descriptions like
// "Sold 3 item(s) for UGX 52,500". Caller adds the "UGX " prefix.
func formatMoney(n float64) string {
	whole := int64(n)
	s := strconv.FormatInt(abs64(whole), 10)
	if len(s) <= 3 {
		if whole < 0 {
			return "-" + s
		}
		return s
	}

	var b strings.Builder
	if whole < 0 {
		b.WriteByte('-')
	}
	first := len(s) % 3
	if first == 0 {
		first = 3
	}
	b.WriteString(s[:first])
	for i := first; i < len(s); i += 3 {
		b.WriteByte(',')
		b.WriteString(s[i : i+3])
	}
	// preserve cents only if they're non-zero
	cents := n - float64(whole)
	if cents > 0.005 || cents < -0.005 {
		b.WriteString(fmt.Sprintf(".%02d", int((cents*100)+0.5)))
	}
	return b.String()
}

func abs64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
