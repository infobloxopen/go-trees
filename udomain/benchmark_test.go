package domain

import (
	"fmt"
	"testing"
)

func BenchmarkNameLess(b *testing.B) {
	for _, s := range []string{
		"first-level-label",
		"second-level-label.first-level-label",
		"third-level-label.second-level-label.first-level-label",
		"fourth-level-label.third-level-label.second-level-label.first-level-label",
	} {
		n, err := MakeNameFromString(s)
		if err != nil {
			b.Fatalf("got error for name %q: %s", s, err)
		}

		b.Run(fmt.Sprintf("%d-dwords", len(n.c)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if n.Less(n) {
					b.Fatalf("expected %q to be not less than itself but it is", n)
				}
			}
		})
	}
}
