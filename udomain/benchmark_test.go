package domain

import "testing"

func BenchmarkNameLess(b *testing.B) {
	n1, err := MakeNameFromString("fourth-level-label.third-level-label.second-level-label.first-level-label-A")
	if err != nil {
		b.Fatal(err)
	}

	n2, err := MakeNameFromString("fourth-level-label.third-level-label.second-level-label.first-level-label-B")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		if !n1.Less(n2) {
			b.Fatalf("expected %q to be less than %q but it is not", n1, n2)
		}
	}
}
