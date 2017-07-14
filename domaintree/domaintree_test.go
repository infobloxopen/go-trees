package domaintree

import (
	"fmt"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func TestInsert(t *testing.T) {
	var r *Node

	r1 := r.Insert("com", "1")
	if r1 == nil {
		t.Errorf("Expected new tree but got nothing")
	}

	r2 := r1.Insert("test.com", "2")
	r3 := r2.Insert("test.net", "3")
	r4 := r3.Insert("example.com", "4")
	r5 := r4.Insert("www.test.com", "5")

	assertTree(r, "empty tree", t)

	assertTree(r1, "single element tree", t,
		"\"com\": \"1\"\n")

	assertTree(r2, "two elements tree", t,
		"\"com\": \"1\"\n",
		"\"test.com\": \"2\"\n")

	assertTree(r3, "three elements tree", t,
		"\"com\": \"1\"\n",
		"\"test.com\": \"2\"\n",
		"\"test.net\": \"3\"\n")

	assertTree(r4, "four elements tree", t,
		"\"com\": \"1\"\n",
		"\"example.com\": \"4\"\n",
		"\"test.com\": \"2\"\n",
		"\"test.net\": \"3\"\n")

	assertTree(r5, "five elements tree", t,
		"\"com\": \"1\"\n",
		"\"example.com\": \"4\"\n",
		"\"test.com\": \"2\"\n",
		"\"www.test.com\": \"5\"\n",
		"\"test.net\": \"3\"\n")
}

func assertTree(r *Node, desc string, t *testing.T, e ...string) {
	pairs := []string{}
	for p := range r.Enumerate() {
		s, ok := p.Value.(string)
		if ok {
			pairs = append(pairs, fmt.Sprintf("%q: %q\n", p.Key, s))
		} else {
			pairs = append(pairs, fmt.Sprintf("%q: %T (%#v)\n", p.Key, p.Value, p.Value))
		}
	}

	ctx := difflib.ContextDiff{
		A:        e,
		B:        pairs,
		FromFile: "Expected",
		ToFile:   "Got"}

	diff, err := difflib.GetContextDiffString(ctx)
	if err != nil {
		panic(fmt.Errorf("Can't compare \"%s\": %s", desc, err))
	}

	if len(diff) > 0 {
		t.Errorf("\"%s\" doesn't match:\n%s", desc, diff)
	}
}
