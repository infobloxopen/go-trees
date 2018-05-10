//+build !amd64

package domain

import "strings"

// Compare function compares two strings and returns -1 if a < b, +1 if a > b and 0 otherwise. Function is optimized for speed and doesn't follow alphabetic order. Caller must ensure that given strings are of the same length.
func Compare(a, b string) int {
	return strings.Compare(a, b)
}
