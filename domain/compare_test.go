// +build amd64

package domain

import "testing"

type domainCompareTestCase struct {
	a string
	b string
	n int
}

var domainCompareTestCases = []domainCompareTestCase{
	{"z", "z", 0},
	{"a", "b", -1},
	{"b", "a", 1},
	{"yz", "yz", 0},
	{"ab", "ba", 255},
	{"ba", "ab", -255},
	{"xyz", "xyz", 0},
	{"abc", "bca", -257},
	{"bca", "abc", 257},
    {"xya", "xyb", -1},
    {"xyb", "xya", 1},
	{"wxyz", "wxyz", 0},
	{"abcd", "bcda", 50265855},
	{"bcda", "abcd", -50265855},
}

func TestDomainCompare(t *testing.T) {
	for i, c := range domainCompareTestCases {
		n := Compare(c.a, c.b)
		if n != c.n {
			t.Errorf("%d: expected %d for %q and %q but got %d", i, c.n, c.a, c.b, n)
		}
	}
}
