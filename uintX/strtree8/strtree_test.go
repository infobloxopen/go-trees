package strtree8

// !!!DON'T EDIT!!! Generated by infobloxopen/go-trees/etc from domaintree{{.bits}} with etc -s uint8 -d uintX.yaml -t ./domaintree\{\{.bits\}\}

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

func TestInsert(t *testing.T) {
	var r *Tree

	r = r.Insert("k", 1)
	assertTree(r, TestSingleNodeTree, "single node tree", t)

	r = NewTree()
	r = r.Insert("k", 1)
	assertTree(r, TestSingleNodeTree, "single node tree", t)

	r = NewTree()
	r = r.Insert("0", 0)
	r = r.Insert("1", 0)
	r = r.Insert("2", 0)
	assertTree(r, TestThreeNodeTreeRed, "tree 012", t)

	r = NewTree()
	r = r.Insert("1", 0)
	r = r.Insert("2", 0)
	r = r.Insert("0", 0)
	assertTree(r, TestThreeNodeTreeRed, "tree 120", t)

	r = NewTree()
	r = r.Insert("2", 0)
	r = r.Insert("0", 0)
	r = r.Insert("1", 0)
	assertTree(r, TestThreeNodeTreeRed, "tree 201", t)

	r = NewTree()
	r = r.Insert("0", 0)
	r = r.Insert("2", 0)
	r = r.Insert("1", 0)
	assertTree(r, TestThreeNodeTreeRed, "tree 021", t)

	r = NewTree()
	r = r.Insert("1", 0)
	r = r.Insert("0", 0)
	r = r.Insert("2", 0)
	assertTree(r, TestThreeNodeTreeRed, "tree 102", t)

	r = NewTree()
	r = r.Insert("2", 0)
	r = r.Insert("1", 0)
	r = r.Insert("0", 0)
	assertTree(r, TestThreeNodeTreeRed, "tree 210", t)

	r = NewTree()
	r = r.Insert("1", 0)
	r = r.Insert("0", 0)
	r = r.Insert("4", 0)
	r = r.Insert("2", 0)
	r = r.Insert("3", 0)
	assertTree(r, TestFiveNodeTree, "tree 10423", t)

	r = NewTree()
	r = r.Insert("F", 0)
	r = r.Insert("E", 0)
	r = r.Insert("D", 0)
	r = r.Insert("C", 0)
	r = r.Insert("B", 0)
	r = r.Insert("A", 0)
	r = r.Insert("9", 0)
	r = r.Insert("8", 0)
	r = r.Insert("7", 0)
	r = r.Insert("6", 0)
	r = r.Insert("5", 0)
	r = r.Insert("4", 0)
	r = r.Insert("3", 0)
	r = r.Insert("2", 0)
	r = r.Insert("1", 0)
	r = r.Insert("0", 0)
	assertTree(r, Test16InversedNodeTree, "tree inversed 16 nodes", t)

	r = NewTree()
	r = r.Insert("0", 0)
	r = r.Insert("1", 0)
	r = r.Insert("2", 0)
	r = r.Insert("3", 0)
	r = r.Insert("4", 0)
	r = r.Insert("5", 0)
	r = r.Insert("6", 0)
	r = r.Insert("7", 0)
	r = r.Insert("8", 0)
	r = r.Insert("9", 0)
	r = r.Insert("A", 0)
	r = r.Insert("B", 0)
	r = r.Insert("C", 0)
	r = r.Insert("D", 0)
	r = r.Insert("E", 0)
	r = r.Insert("F", 0)
	assertTree(r, Test16DirectNodeTree, "tree direct 16 nodes", t)

	r = NewTree()
	r = r.Insert("0", 0)
	r = r.Insert("2", 0)
	r = r.Insert("4", 0)
	r = r.Insert("6", 0)
	r = r.Insert("8", 0)
	r = r.Insert("A", 0)
	r = r.Insert("C", 0)
	r = r.Insert("E", 0)
	r = r.Insert("1", 0)
	r = r.Insert("3", 0)
	r = r.Insert("5", 0)
	r = r.Insert("7", 0)
	r = r.Insert("9", 0)
	r = r.Insert("B", 0)
	r = r.Insert("D", 0)
	r = r.Insert("F", 0)
	assertTree(r, Test16AlternatingNodeTree, "tree alternating 16 nodes", t)

	r = NewTree()
	r = r.Insert("F", 0)
	r = r.Insert("D", 0)
	r = r.Insert("B", 0)
	r = r.Insert("9", 0)
	r = r.Insert("7", 0)
	r = r.Insert("5", 0)
	r = r.Insert("3", 0)
	r = r.Insert("1", 0)
	r = r.Insert("E", 0)
	r = r.Insert("C", 0)
	r = r.Insert("A", 0)
	r = r.Insert("8", 0)
	r = r.Insert("6", 0)
	r = r.Insert("4", 0)
	r = r.Insert("2", 0)
	r = r.Insert("0", 0)
	assertTree(r, Test16AlternatingInversedNodeTree, "tree alternating inversed 16 nodes", t)

	r = NewTree()
	r = r.Insert("0", 0)
	r = r.Insert("3", 0)
	r = r.Insert("6", 0)
	r = r.Insert("9", 0)
	r = r.Insert("C", 0)
	r = r.Insert("F", 0)
	r = r.Insert("1", 0)
	r = r.Insert("2", 0)
	r = r.Insert("4", 0)
	r = r.Insert("5", 0)
	r = r.Insert("7", 0)
	r = r.Insert("8", 0)
	r = r.Insert("A", 0)
	r = r.Insert("B", 0)
	r = r.Insert("D", 0)
	r = r.Insert("E", 0)
	assertTree(r, Test16_3AltNodeTree, "tree alternating by 3 16 nodes", t)

	r = NewTree()
	r = r.Insert("00", 0)
	r = r.Insert("02", 0)
	r = r.Insert("04", 0)
	r = r.Insert("06", 0)
	r = r.Insert("08", 0)
	r = r.Insert("0A", 0)
	r = r.Insert("0C", 0)
	r = r.Insert("0E", 0)
	r = r.Insert("10", 0)
	r = r.Insert("12", 0)
	r = r.Insert("14", 0)
	r = r.Insert("16", 0)
	r = r.Insert("18", 0)
	r = r.Insert("1A", 0)
	r = r.Insert("1C", 0)
	r = r.Insert("1E", 0)
	r = r.Insert("01", 0)
	r = r.Insert("03", 0)
	r = r.Insert("05", 0)
	r = r.Insert("07", 0)
	r = r.Insert("09", 0)
	r = r.Insert("0B", 0)
	r = r.Insert("0D", 0)
	r = r.Insert("0F", 0)
	r = r.Insert("11", 0)
	r = r.Insert("13", 0)
	r = r.Insert("15", 0)
	r = r.Insert("17", 0)
	r = r.Insert("19", 0)
	r = r.Insert("1B", 0)
	r = r.Insert("1D", 0)
	r = r.Insert("1F", 0)
	assertTree(r, Test32AlternatingNodeTree, "tree with alternating 32 nodes", t)

	r = nil
	r = r.Insert("00", 1)
	assertTree(r, TestTreeSameNodeOnce, "tree with same node first insertion", t)
	r = r.Insert("00", 2)
	assertTree(r, TestTreeSameNodeTwice, "tree with same node second insertion", t)
}

func TestInplaceInsert(t *testing.T) {
	r := NewTree()

	r.InplaceInsert("k", 1)
	assertTree(r, TestSingleNodeTree, "single node inplace tree", t)

	r = NewTree()
	r.InplaceInsert("k", 1)
	assertTree(r, TestSingleNodeTree, "single node inplace tree", t)

	r = NewTree()
	r.InplaceInsert("0", 0)
	r.InplaceInsert("1", 0)
	r.InplaceInsert("2", 0)
	assertTree(r, TestThreeNodeTreeRed, "inplace tree 012", t)

	r = NewTree()
	r.InplaceInsert("1", 0)
	r.InplaceInsert("2", 0)
	r.InplaceInsert("0", 0)
	assertTree(r, TestThreeNodeTreeRed, "inplace tree 120", t)

	r = NewTree()
	r.InplaceInsert("2", 0)
	r.InplaceInsert("0", 0)
	r.InplaceInsert("1", 0)
	assertTree(r, TestThreeNodeTreeRed, "inplace tree 201", t)

	r = NewTree()
	r.InplaceInsert("0", 0)
	r.InplaceInsert("2", 0)
	r.InplaceInsert("1", 0)
	assertTree(r, TestThreeNodeTreeRed, "inplace tree 021", t)

	r = NewTree()
	r.InplaceInsert("1", 0)
	r.InplaceInsert("0", 0)
	r.InplaceInsert("2", 0)
	assertTree(r, TestThreeNodeTreeRed, "inplace tree 102", t)

	r = NewTree()
	r.InplaceInsert("2", 0)
	r.InplaceInsert("1", 0)
	r.InplaceInsert("0", 0)
	assertTree(r, TestThreeNodeTreeRed, "inplace tree 210", t)

	r = NewTree()
	r.InplaceInsert("1", 0)
	r.InplaceInsert("0", 0)
	r.InplaceInsert("4", 0)
	r.InplaceInsert("2", 0)
	r.InplaceInsert("3", 0)
	assertTree(r, TestFiveNodeTree, "inplace tree 10423", t)

	r = NewTree()
	r.InplaceInsert("F", 0)
	r.InplaceInsert("E", 0)
	r.InplaceInsert("D", 0)
	r.InplaceInsert("C", 0)
	r.InplaceInsert("B", 0)
	r.InplaceInsert("A", 0)
	r.InplaceInsert("9", 0)
	r.InplaceInsert("8", 0)
	r.InplaceInsert("7", 0)
	r.InplaceInsert("6", 0)
	r.InplaceInsert("5", 0)
	r.InplaceInsert("4", 0)
	r.InplaceInsert("3", 0)
	r.InplaceInsert("2", 0)
	r.InplaceInsert("1", 0)
	r.InplaceInsert("0", 0)
	assertTree(r, Test16InversedNodeTree, "inplace tree inversed 16 nodes", t)

	r = NewTree()
	r.InplaceInsert("0", 0)
	r.InplaceInsert("1", 0)
	r.InplaceInsert("2", 0)
	r.InplaceInsert("3", 0)
	r.InplaceInsert("4", 0)
	r.InplaceInsert("5", 0)
	r.InplaceInsert("6", 0)
	r.InplaceInsert("7", 0)
	r.InplaceInsert("8", 0)
	r.InplaceInsert("9", 0)
	r.InplaceInsert("A", 0)
	r.InplaceInsert("B", 0)
	r.InplaceInsert("C", 0)
	r.InplaceInsert("D", 0)
	r.InplaceInsert("E", 0)
	r.InplaceInsert("F", 0)
	assertTree(r, Test16DirectNodeTree, "inplace tree direct 16 nodes", t)

	r = NewTree()
	r.InplaceInsert("0", 0)
	r.InplaceInsert("2", 0)
	r.InplaceInsert("4", 0)
	r.InplaceInsert("6", 0)
	r.InplaceInsert("8", 0)
	r.InplaceInsert("A", 0)
	r.InplaceInsert("C", 0)
	r.InplaceInsert("E", 0)
	r.InplaceInsert("1", 0)
	r.InplaceInsert("3", 0)
	r.InplaceInsert("5", 0)
	r.InplaceInsert("7", 0)
	r.InplaceInsert("9", 0)
	r.InplaceInsert("B", 0)
	r.InplaceInsert("D", 0)
	r.InplaceInsert("F", 0)
	assertTree(r, Test16AlternatingNodeTree, "inplace tree alternating 16 nodes", t)

	r = NewTree()
	r.InplaceInsert("F", 0)
	r.InplaceInsert("D", 0)
	r.InplaceInsert("B", 0)
	r.InplaceInsert("9", 0)
	r.InplaceInsert("7", 0)
	r.InplaceInsert("5", 0)
	r.InplaceInsert("3", 0)
	r.InplaceInsert("1", 0)
	r.InplaceInsert("E", 0)
	r.InplaceInsert("C", 0)
	r.InplaceInsert("A", 0)
	r.InplaceInsert("8", 0)
	r.InplaceInsert("6", 0)
	r.InplaceInsert("4", 0)
	r.InplaceInsert("2", 0)
	r.InplaceInsert("0", 0)
	assertTree(r, Test16AlternatingInversedNodeTree, "inplace tree alternating inversed 16 nodes", t)

	r = NewTree()
	r.InplaceInsert("0", 0)
	r.InplaceInsert("3", 0)
	r.InplaceInsert("6", 0)
	r.InplaceInsert("9", 0)
	r.InplaceInsert("C", 0)
	r.InplaceInsert("F", 0)
	r.InplaceInsert("1", 0)
	r.InplaceInsert("2", 0)
	r.InplaceInsert("4", 0)
	r.InplaceInsert("5", 0)
	r.InplaceInsert("7", 0)
	r.InplaceInsert("8", 0)
	r.InplaceInsert("A", 0)
	r.InplaceInsert("B", 0)
	r.InplaceInsert("D", 0)
	r.InplaceInsert("E", 0)
	assertTree(r, Test16_3AltNodeTree, "inplace tree alternating by 3 16 nodes", t)

	r = NewTree()
	r.InplaceInsert("00", 0)
	r.InplaceInsert("02", 0)
	r.InplaceInsert("04", 0)
	r.InplaceInsert("06", 0)
	r.InplaceInsert("08", 0)
	r.InplaceInsert("0A", 0)
	r.InplaceInsert("0C", 0)
	r.InplaceInsert("0E", 0)
	r.InplaceInsert("10", 0)
	r.InplaceInsert("12", 0)
	r.InplaceInsert("14", 0)
	r.InplaceInsert("16", 0)
	r.InplaceInsert("18", 0)
	r.InplaceInsert("1A", 0)
	r.InplaceInsert("1C", 0)
	r.InplaceInsert("1E", 0)
	r.InplaceInsert("01", 0)
	r.InplaceInsert("03", 0)
	r.InplaceInsert("05", 0)
	r.InplaceInsert("07", 0)
	r.InplaceInsert("09", 0)
	r.InplaceInsert("0B", 0)
	r.InplaceInsert("0D", 0)
	r.InplaceInsert("0F", 0)
	r.InplaceInsert("11", 0)
	r.InplaceInsert("13", 0)
	r.InplaceInsert("15", 0)
	r.InplaceInsert("17", 0)
	r.InplaceInsert("19", 0)
	r.InplaceInsert("1B", 0)
	r.InplaceInsert("1D", 0)
	r.InplaceInsert("1F", 0)
	assertTree(r, Test32AlternatingNodeTree, "inplace tree with alternating 32 nodes", t)

	r = nil
	assertPanic(func() { r.InplaceInsert("00", 0) }, "nil tree inplace insertion", t)
}

func TestGet(t *testing.T) {
	var r *Tree

	v, ok := r.Get("0")
	if ok {
		t.Errorf("Expected nothing but got %T (%#v)", v, v)
	}

	r = NewTree()
	r = r.Insert("1", 1)
	r = r.Insert("0", 0)
	r = r.Insert("4", 4)
	r = r.Insert("2", 2)
	r = r.Insert("3", 3)

	v, ok = r.Get("3")
	if !ok {
		t.Errorf("Expected 3 but got nothing")
	} else if v != 3 {
		t.Errorf("Expected 3 but got %d", v)
	}

	v, ok = r.Get("F")
	if ok {
		t.Errorf("Expected nothing but got %T (%#v)", v, v)
	}
}

func TestEnumerate(t *testing.T) {
	var r *Tree

	assertEnumerate(r.Enumerate(), "empty tree", t)

	r = NewTree()
	r = r.Insert("1", 1)
	r = r.Insert("0", 0)
	r = r.Insert("4", 4)
	r = r.Insert("2", 2)
	r = r.Insert("3", 3)
	assertEnumerate(r.Enumerate(), "enumeration of tree 10423", t,
		"\"0\": 0\n",
		"\"1\": 1\n",
		"\"2\": 2\n",
		"\"3\": 3\n",
		"\"4\": 4\n")
}

func TestDelete(t *testing.T) {
	var r *Tree

	r, ok := r.Delete("test")
	if ok {
		t.Errorf("Expected nothing to be deleted from empty tree but something has been deleted:\n%s", r.Dot())
	}

	r = NewTree()
	r = r.Insert("0", 0)
	r = r.Insert("3", 0)
	r = r.Insert("6", 0)
	r = r.Insert("9", 0)
	r = r.Insert("C", 0)
	r = r.Insert("F", 0)
	r = r.Insert("1", 0)
	r = r.Insert("2", 0)
	r = r.Insert("4", 0)
	r = r.Insert("5", 0)
	r = r.Insert("7", 0)
	r = r.Insert("8", 0)
	r = r.Insert("A", 0)
	r = r.Insert("B", 0)
	r = r.Insert("D", 0)
	r = r.Insert("E", 0)

	r, ok = r.Delete("81")
	if ok {
		t.Errorf("Expected nothing to be deleted by key \"81\" but something has been deleted")
	}
	assertTree(r, TestTreeAfterNonExistingNodeDel, "tree after non-existing node deletion", t)

	r, ok = r.Delete("6")
	if !ok {
		t.Errorf("Expected node \"6\" to be deleted but got nothing")
	}
	assertTree(r, TestTreeAfterNode6Deletion, "tree after node 6 deletion", t)

	r, ok = r.Delete("7")
	if !ok {
		t.Errorf("Expected node \"7\" to be deleted but got nothing")
	}
	r, ok = r.Delete("8")
	if !ok {
		t.Errorf("Expected node \"8\" to be deleted but got nothing")
	}
	r, ok = r.Delete("5")
	if !ok {
		t.Errorf("Expected node \"5\" to be deleted but got nothing")
	}
	r, ok = r.Delete("9")
	if !ok {
		t.Errorf("Expected node \"9\" to be deleted but got nothing")
	}
	assertTree(r, TestTreeAfterNodes7859Deletion, "tree after nodes 7, 8, 5 and 9 deletion", t)

	r, ok = r.Delete("C")
	if !ok {
		t.Errorf("Expected node \"C\" to be deleted but got nothing")
	}
	r, ok = r.Delete("E")
	if !ok {
		t.Errorf("Expected node \"E\" to be deleted but got nothing")
	}
	r, ok = r.Delete("D")
	if !ok {
		t.Errorf("Expected node \"D\" to be deleted but got nothing")
	}
	r, ok = r.Delete("A")
	if !ok {
		t.Errorf("Expected node \"A\" to be deleted but got nothing")
	}
	r, ok = r.Delete("B")
	if !ok {
		t.Errorf("Expected node \"B\" to be deleted but got nothing")
	}
	r, ok = r.Delete("4")
	if !ok {
		t.Errorf("Expected node \"4\" to be deleted but got nothing")
	}
	r, ok = r.Delete("F")
	if !ok {
		t.Errorf("Expected node \"F\" to be deleted but got nothing")
	}
	r, ok = r.Delete("0")
	if !ok {
		t.Errorf("Expected node \"0\" to be deleted but got nothing")
	}
	r, ok = r.Delete("3")
	if !ok {
		t.Errorf("Expected node \"3\" to be deleted but got nothing")
	}
	r, ok = r.Delete("1")
	if !ok {
		t.Errorf("Expected node \"1\" to be deleted but got nothing")
	}
	r, ok = r.Delete("2")
	if !ok {
		t.Errorf("Expected node \"2\" to be deleted but got nothing")
	}
	assertTree(r, TestEmptyTree, "tree after rest nodes deletion", t)
}

func TestIsEmpty(t *testing.T) {
	var r *Tree

	if !r.IsEmpty() {
		t.Errorf("Expected nil tree to be empty")
	}

	r = NewTree()
	r = r.Insert("0", 0)
	r = r.Insert("3", 0)
	r = r.Insert("6", 0)
	if r.IsEmpty() {
		t.Errorf("Expected three nodes tree to be not empty")
	}

	r, ok := r.Delete("3")
	if !ok {
		t.Errorf("Expected element \"3\" to be deleted")
	}

	if r.IsEmpty() {
		t.Errorf("Expected two nodes tree to be not empty")
	}

	r, ok = r.Delete("0")
	r, ok = r.Delete("6")

	if !r.IsEmpty() {
		t.Errorf("Expected empty non-nil tree to be empty")
	}
}

const (
	TestEmptyTree = `digraph d {
N0 [label="nil" style=filled fontcolor=white fillcolor=black]
}
`

	TestSingleNodeTree = `digraph d {
N0 [label="k: \"k\" v: \"1\"" style=filled fontcolor=white fillcolor=black]
}
`

	TestThreeNodeTreeRed = `digraph d {
N0 [label="k: \"1\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="k: \"0\" v: \"0\"" style=filled fillcolor=red]
N2 [label="k: \"2\" v: \"0\"" style=filled fillcolor=red]
}
`

	TestFiveNodeTree = `digraph d {
N0 [label="k: \"1\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="k: \"0\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N2 [label="k: \"3\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N2 -> { N3 N4 }
N3 [label="k: \"2\" v: \"0\"" style=filled fillcolor=red]
N4 [label="k: \"4\" v: \"0\"" style=filled fillcolor=red]
}
`

	Test16InversedNodeTree = `digraph d {
N0 [label="k: \"C\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="k: \"8\" v: \"0\"" style=filled fillcolor=red]
N1 -> { N3 N4 }
N2 [label="k: \"E\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="k: \"4\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N3 -> { N7 N8 }
N4 [label="k: \"A\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N4 -> { N9 N10 }
N5 [label="k: \"D\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N6 [label="k: \"F\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N7 [label="k: \"2\" v: \"0\"" style=filled fillcolor=red]
N7 -> { N11 N12 }
N8 [label="k: \"6\" v: \"0\"" style=filled fillcolor=red]
N8 -> { N13 N14 }
N9 [label="k: \"9\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N10 [label="k: \"B\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N11 [label="k: \"1\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N11 -> { N15 N16 }
N12 [label="k: \"3\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N13 [label="k: \"5\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N14 [label="k: \"7\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N15 [label="k: \"0\" v: \"0\"" style=filled fillcolor=red]
N16 [label="nil" style=filled fontcolor=white fillcolor=black]
}
`

	Test16DirectNodeTree = `digraph d {
N0 [label="k: \"3\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="k: \"1\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="k: \"7\" v: \"0\"" style=filled fillcolor=red]
N2 -> { N5 N6 }
N3 [label="k: \"0\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N4 [label="k: \"2\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N5 [label="k: \"5\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N5 -> { N7 N8 }
N6 [label="k: \"B\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N6 -> { N9 N10 }
N7 [label="k: \"4\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N8 [label="k: \"6\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N9 [label="k: \"9\" v: \"0\"" style=filled fillcolor=red]
N9 -> { N11 N12 }
N10 [label="k: \"D\" v: \"0\"" style=filled fillcolor=red]
N10 -> { N13 N14 }
N11 [label="k: \"8\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N12 [label="k: \"A\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N13 [label="k: \"C\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N14 [label="k: \"E\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N14 -> { N15 N16 }
N15 [label="nil" style=filled fontcolor=white fillcolor=black]
N16 [label="k: \"F\" v: \"0\"" style=filled fillcolor=red]
}
`

	Test16AlternatingNodeTree = `digraph d {
N0 [label="k: \"6\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="k: \"2\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="k: \"A\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="k: \"0\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N3 -> { N7 N8 }
N4 [label="k: \"4\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N4 -> { N9 N10 }
N5 [label="k: \"8\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N5 -> { N11 N12 }
N6 [label="k: \"C\" v: \"0\"" style=filled fillcolor=red]
N6 -> { N13 N14 }
N7 [label="nil" style=filled fontcolor=white fillcolor=black]
N8 [label="k: \"1\" v: \"0\"" style=filled fillcolor=red]
N9 [label="k: \"3\" v: \"0\"" style=filled fillcolor=red]
N10 [label="k: \"5\" v: \"0\"" style=filled fillcolor=red]
N11 [label="k: \"7\" v: \"0\"" style=filled fillcolor=red]
N12 [label="k: \"9\" v: \"0\"" style=filled fillcolor=red]
N13 [label="k: \"B\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N14 [label="k: \"E\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N14 -> { N15 N16 }
N15 [label="k: \"D\" v: \"0\"" style=filled fillcolor=red]
N16 [label="k: \"F\" v: \"0\"" style=filled fillcolor=red]
}
`

	Test16AlternatingInversedNodeTree = `digraph d {
N0 [label="k: \"9\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="k: \"5\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="k: \"D\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="k: \"3\" v: \"0\"" style=filled fillcolor=red]
N3 -> { N7 N8 }
N4 [label="k: \"7\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N4 -> { N9 N10 }
N5 [label="k: \"B\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N5 -> { N11 N12 }
N6 [label="k: \"F\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N6 -> { N13 N14 }
N7 [label="k: \"1\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N7 -> { N15 N16 }
N8 [label="k: \"4\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N9 [label="k: \"6\" v: \"0\"" style=filled fillcolor=red]
N10 [label="k: \"8\" v: \"0\"" style=filled fillcolor=red]
N11 [label="k: \"A\" v: \"0\"" style=filled fillcolor=red]
N12 [label="k: \"C\" v: \"0\"" style=filled fillcolor=red]
N13 [label="k: \"E\" v: \"0\"" style=filled fillcolor=red]
N14 [label="nil" style=filled fontcolor=white fillcolor=black]
N15 [label="k: \"0\" v: \"0\"" style=filled fillcolor=red]
N16 [label="k: \"2\" v: \"0\"" style=filled fillcolor=red]
}
`

	Test16_3AltNodeTree = `digraph d {
N0 [label="k: \"5\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="k: \"3\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="k: \"9\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="k: \"1\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N3 -> { N7 N8 }
N4 [label="k: \"4\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N5 [label="k: \"7\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N5 -> { N9 N10 }
N6 [label="k: \"C\" v: \"0\"" style=filled fillcolor=red]
N6 -> { N11 N12 }
N7 [label="k: \"0\" v: \"0\"" style=filled fillcolor=red]
N8 [label="k: \"2\" v: \"0\"" style=filled fillcolor=red]
N9 [label="k: \"6\" v: \"0\"" style=filled fillcolor=red]
N10 [label="k: \"8\" v: \"0\"" style=filled fillcolor=red]
N11 [label="k: \"A\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N11 -> { N13 N14 }
N12 [label="k: \"E\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N12 -> { N15 N16 }
N13 [label="nil" style=filled fontcolor=white fillcolor=black]
N14 [label="k: \"B\" v: \"0\"" style=filled fillcolor=red]
N15 [label="k: \"D\" v: \"0\"" style=filled fillcolor=red]
N16 [label="k: \"F\" v: \"0\"" style=filled fillcolor=red]
}
`

	Test32AlternatingNodeTree = `digraph d {
N0 [label="k: \"0E\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="k: \"06\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="k: \"16\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="k: \"02\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N3 -> { N7 N8 }
N4 [label="k: \"0A\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N4 -> { N9 N10 }
N5 [label="k: \"12\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N5 -> { N11 N12 }
N6 [label="k: \"1A\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N6 -> { N13 N14 }
N7 [label="k: \"00\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N7 -> { N15 N16 }
N8 [label="k: \"04\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N8 -> { N17 N18 }
N9 [label="k: \"08\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N9 -> { N19 N20 }
N10 [label="k: \"0C\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N10 -> { N21 N22 }
N11 [label="k: \"10\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N11 -> { N23 N24 }
N12 [label="k: \"14\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N12 -> { N25 N26 }
N13 [label="k: \"18\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N13 -> { N27 N28 }
N14 [label="k: \"1C\" v: \"0\"" style=filled fillcolor=red]
N14 -> { N29 N30 }
N15 [label="nil" style=filled fontcolor=white fillcolor=black]
N16 [label="k: \"01\" v: \"0\"" style=filled fillcolor=red]
N17 [label="k: \"03\" v: \"0\"" style=filled fillcolor=red]
N18 [label="k: \"05\" v: \"0\"" style=filled fillcolor=red]
N19 [label="k: \"07\" v: \"0\"" style=filled fillcolor=red]
N20 [label="k: \"09\" v: \"0\"" style=filled fillcolor=red]
N21 [label="k: \"0B\" v: \"0\"" style=filled fillcolor=red]
N22 [label="k: \"0D\" v: \"0\"" style=filled fillcolor=red]
N23 [label="k: \"0F\" v: \"0\"" style=filled fillcolor=red]
N24 [label="k: \"11\" v: \"0\"" style=filled fillcolor=red]
N25 [label="k: \"13\" v: \"0\"" style=filled fillcolor=red]
N26 [label="k: \"15\" v: \"0\"" style=filled fillcolor=red]
N27 [label="k: \"17\" v: \"0\"" style=filled fillcolor=red]
N28 [label="k: \"19\" v: \"0\"" style=filled fillcolor=red]
N29 [label="k: \"1B\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N30 [label="k: \"1E\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N30 -> { N31 N32 }
N31 [label="k: \"1D\" v: \"0\"" style=filled fillcolor=red]
N32 [label="k: \"1F\" v: \"0\"" style=filled fillcolor=red]
}
`

	TestTreeAfterNonExistingNodeDel = `digraph d {
N0 [label="k: \"5\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="k: \"3\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="k: \"C\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="k: \"1\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N3 -> { N7 N8 }
N4 [label="k: \"4\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N5 [label="k: \"9\" v: \"0\"" style=filled fillcolor=red]
N5 -> { N9 N10 }
N6 [label="k: \"E\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N6 -> { N11 N12 }
N7 [label="k: \"0\" v: \"0\"" style=filled fillcolor=red]
N8 [label="k: \"2\" v: \"0\"" style=filled fillcolor=red]
N9 [label="k: \"7\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N9 -> { N13 N14 }
N10 [label="k: \"A\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N10 -> { N15 N16 }
N11 [label="k: \"D\" v: \"0\"" style=filled fillcolor=red]
N12 [label="k: \"F\" v: \"0\"" style=filled fillcolor=red]
N13 [label="k: \"6\" v: \"0\"" style=filled fillcolor=red]
N14 [label="k: \"8\" v: \"0\"" style=filled fillcolor=red]
N15 [label="nil" style=filled fontcolor=white fillcolor=black]
N16 [label="k: \"B\" v: \"0\"" style=filled fillcolor=red]
}
`

	TestTreeAfterNode6Deletion = `digraph d {
N0 [label="k: \"5\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="k: \"3\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="k: \"C\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="k: \"1\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N3 -> { N7 N8 }
N4 [label="k: \"4\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N5 [label="k: \"9\" v: \"0\"" style=filled fillcolor=red]
N5 -> { N9 N10 }
N6 [label="k: \"E\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N6 -> { N11 N12 }
N7 [label="k: \"0\" v: \"0\"" style=filled fillcolor=red]
N8 [label="k: \"2\" v: \"0\"" style=filled fillcolor=red]
N9 [label="k: \"7\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N9 -> { N13 N14 }
N10 [label="k: \"A\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N10 -> { N15 N16 }
N11 [label="k: \"D\" v: \"0\"" style=filled fillcolor=red]
N12 [label="k: \"F\" v: \"0\"" style=filled fillcolor=red]
N13 [label="nil" style=filled fontcolor=white fillcolor=black]
N14 [label="k: \"8\" v: \"0\"" style=filled fillcolor=red]
N15 [label="nil" style=filled fontcolor=white fillcolor=black]
N16 [label="k: \"B\" v: \"0\"" style=filled fillcolor=red]
}
`

	TestTreeAfterNodes7859Deletion = `digraph d {
N0 [label="k: \"A\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N0 -> { N1 N2 }
N1 [label="k: \"2\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N1 -> { N3 N4 }
N2 [label="k: \"C\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N2 -> { N5 N6 }
N3 [label="k: \"1\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N3 -> { N7 N8 }
N4 [label="k: \"4\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N4 -> { N9 N10 }
N5 [label="k: \"B\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N6 [label="k: \"E\" v: \"0\"" style=filled fontcolor=white fillcolor=black]
N6 -> { N11 N12 }
N7 [label="k: \"0\" v: \"0\"" style=filled fillcolor=red]
N8 [label="nil" style=filled fontcolor=white fillcolor=black]
N9 [label="k: \"3\" v: \"0\"" style=filled fillcolor=red]
N10 [label="nil" style=filled fontcolor=white fillcolor=black]
N11 [label="k: \"D\" v: \"0\"" style=filled fillcolor=red]
N12 [label="k: \"F\" v: \"0\"" style=filled fillcolor=red]
}
`

	TestTreeSameNodeOnce = `digraph d {
N0 [label="k: \"00\" v: \"1\"" style=filled fontcolor=white fillcolor=black]
}
`

	TestTreeSameNodeTwice = `digraph d {
N0 [label="k: \"00\" v: \"2\"" style=filled fontcolor=white fillcolor=black]
}
`
)

func assertTree(r *Tree, e, desc string, t *testing.T) {
	assertStringLists(difflib.SplitLines(r.Dot()), difflib.SplitLines(e), desc, t)
}

func assertEnumerate(ch chan Pair, desc string, t *testing.T, e ...string) {
	pairs := []string{}
	for p := range ch {
		pairs = append(pairs, fmt.Sprintf("%q: %d\n", p.Key, p.Value))
	}

	assertStringLists(pairs, e, desc, t)
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

func assertPanic(f func(), desc string, t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic from %s but got nothing", desc)
		}
	}()

	f()
}
