package strtree

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func TestNewTree(t *testing.T) {
	r := NewTree()

	assertTree(r, TestEmptyTree, "empty tree", t)
}

func TestNewTreeWithCustomComparison(t *testing.T) {
	r := NewTreeWithCustomComparison(strings.Compare)

	assertTree(r, TestEmptyTree, "empty tree", t)
}

const TestEmptyTree = `digraph d {
}
`

func assertTree(r *Tree, e, desc string, t *testing.T) {
	assertStringLists(difflib.SplitLines(r.Dot()), difflib.SplitLines(e), desc, t)
}

func assertStringLists(v, e []string, desc string, t *testing.T) {
	ctx := difflib.ContextDiff{
		A:        e,
		B:        v,
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
