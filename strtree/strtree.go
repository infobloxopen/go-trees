// Package strtree implements red-black tree for key value pairs with string keys and custom comparison.
package strtree

import "strings"

// Compare defines function interface for custom comparison. Function implementing the interface should return value less than zero if its first argument precedes second one, zero if both are equal and positive if the second precedes.
type Compare func(a, b string) int

// Tree is a red-black tree for key-value pairs where key is string.
type Tree struct {
	root    node
	compare Compare
}

// NewTree creates empty tree with default comparison operation (strings.Compare).
func NewTree() *Tree {
	return &Tree{compare: strings.Compare}
}

// NewTreeWithCustomComparison creates empty tree with given comparison operation.
func NewTreeWithCustomComparison(compare Compare) *Tree {
	return &Tree{compare: compare}
}

// Dot dumps tree to Graphviz .dot format.
func (t *Tree) Dot() string {
	return "digraph d {\n}\n"
}
