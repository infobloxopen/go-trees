package domaintree

import (
	"fmt"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func TestInsert(t *testing.T) {
	var r *Node

	r = r.Insert("com", "1")
	r = r.Insert("test.com", "2")
	r = r.Insert("test.net", "3")
	r = r.Insert("example.com", "4")
	r = r.Insert("www.test.com", "5")
	if r == nil {
		t.Errorf("Expected new tree but got nothing")
	}

	assertTree(r, "single element tree", t,
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
