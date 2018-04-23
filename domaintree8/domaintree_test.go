package domaintree8

import (
	"fmt"
	"testing"

	"github.com/pmezard/go-difflib/difflib"

	"github.com/infobloxopen/go-trees/domain"
)

func TestInsert(t *testing.T) {
	var r *Node

	r1 := r.Insert("com", 1)
	if r1 == nil {
		t.Error("Expected new tree but got nothing")
	}

	r2 := r1.Insert("test.com", 2)
	r3 := r2.Insert("test.net", 3)
	r4 := r3.Insert("example.com", 4)
	r5 := r4.Insert("www.test.com.", 5)
	r6 := r5.Insert(".", 6)
	r7 := r6.Insert("", 7)

	assertTree(r, "empty tree", t)

	assertTree(r1, "single element tree", t,
		"\"com\": 1\n")

	assertTree(r2, "two elements tree", t,
		"\"com\": 1\n",
		"\"test.com\": 2\n")

	assertTree(r3, "three elements tree", t,
		"\"com\": 1\n",
		"\"test.com\": 2\n",
		"\"test.net\": 3\n")

	assertTree(r4, "four elements tree", t,
		"\"com\": 1\n",
		"\"test.com\": 2\n",
		"\"example.com\": 4\n",
		"\"test.net\": 3\n")

	assertTree(r5, "five elements tree", t,
		"\"com\": 1\n",
		"\"test.com\": 2\n",
		"\"www.test.com\": 5\n",
		"\"example.com\": 4\n",
		"\"test.net\": 3\n")

	assertTree(r6, "siz elements tree", t,
		"\"\": 6\n",
		"\"com\": 1\n",
		"\"test.com\": 2\n",
		"\"www.test.com\": 5\n",
		"\"example.com\": 4\n",
		"\"test.net\": 3\n")

	assertTree(r7, "five elements tree", t,
		"\"\": 7\n",
		"\"com\": 1\n",
		"\"test.com\": 2\n",
		"\"www.test.com\": 5\n",
		"\"example.com\": 4\n",
		"\"test.net\": 3\n")

	r = r.Insert("AbCdEfGhIjKlMnOpQrStUvWxYz.aBcDeFgHiJkLmNoPqRsTuVwXyZ", 255)
	assertTree(r, "case-check tree", t,
		"\"abcdefghijklmnopqrstuvwxyz.abcdefghijklmnopqrstuvwxyz\": 255\n")
}

func TestInplaceInsert(t *testing.T) {
	r := &Node{}
	assertTree(r, "empty inplace tree", t)

	r.InplaceInsert("com", 1)
	assertTree(r, "single element inplace tree", t,
		"\"com\": 1\n")

	r.InplaceInsert("test.com", 2)
	assertTree(r, "two elements inplace tree", t,
		"\"com\": 1\n",
		"\"test.com\": 2\n")

	r.InplaceInsert("test.net", 3)
	assertTree(r, "three elements inplace tree", t,
		"\"com\": 1\n",
		"\"test.com\": 2\n",
		"\"test.net\": 3\n")

	r.InplaceInsert("example.com", 4)
	assertTree(r, "four elements inplace tree", t,
		"\"com\": 1\n",
		"\"test.com\": 2\n",
		"\"example.com\": 4\n",
		"\"test.net\": 3\n")

	r.InplaceInsert("www.test.com", 5)
	assertTree(r, "five elements tree", t,
		"\"com\": 1\n",
		"\"test.com\": 2\n",
		"\"www.test.com\": 5\n",
		"\"example.com\": 4\n",
		"\"test.net\": 3\n")
}

func TestGet(t *testing.T) {
	var r *Node

	v, ok := r.Get("test.com")
	assertValue(v, ok, 0, false, "fetching from empty tree", t)

	r = r.Insert("com", 1)
	r = r.Insert("test.com", 2)
	r = r.Insert("test.net", 3)
	r = r.Insert("example.com", 4)
	r = r.Insert("www.test.com", 5)

	v, ok = r.Get("test.com")
	assertValue(v, ok, 2, true, "fetching \"test.com\" from tree", t)

	v, ok = r.Get("www.test.com")
	assertValue(v, ok, 5, true, "fetching \"www.test.com\" from tree", t)

	v, ok = r.Get("ns.test.com")
	assertValue(v, ok, 2, true, "fetching \"ns.test.com\" from tree", t)

	v, ok = r.Get("test.org")
	assertValue(v, ok, 0, false, "fetching \"test.org\" from tree", t)

	v, ok = r.Get("nS.tEsT.cOm")
	assertValue(v, ok, 2, true, "fetching \"nS.tEsT.cOm\" from tree", t)
}

func TestWireGet(t *testing.T) {
	var r *Node

	v, ok, err := r.WireGet(domain.WireNameLower("\x04test\x03com\x00"))
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		assertValue(v, ok, 0, false, "fetching from empty tree", t)
	}

	r = r.Insert("com", 1)
	r = r.Insert("test.com", 2)
	r = r.Insert("test.net", 3)
	r = r.Insert("example.com", 4)
	r = r.Insert("www.test.com", 5)

	_, _, err = r.WireGet(domain.WireNameLower("\xC0\x2F"))
	if err != domain.ErrCompressedName {
		t.Errorf("Expected %q error but got %q", domain.ErrCompressedName, err)
	}

	_, _, err = r.WireGet(domain.WireNameLower("\x04test\x20com\x00"))
	if err != domain.ErrLabelTooLong {
		t.Errorf("Expected %q error but got %q", domain.ErrLabelTooLong, err)
	}

	_, _, err = r.WireGet(domain.WireNameLower("\x04test\x00\x03com\x00"))
	if err != domain.ErrEmptyLabel {
		t.Errorf("Expected %q error but got %q", domain.ErrEmptyLabel, err)
	}

	_, _, err = r.WireGet(domain.WireNameLower(
		"\x3ftoooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo" +
			"\x3floooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong" +
			"\x3fdoooooooooooooooooooooooooooooooooooooooooooooooooooooooooomain" +
			"\x3fnaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaame" +
			"\x00"))
	if err != domain.ErrNameTooLong {
		t.Errorf("Expected %q error but got %q", domain.ErrNameTooLong, err)
	}

	v, ok, err = r.WireGet(domain.WireNameLower("\x04test\x03com\x00"))
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		assertValue(v, ok, 2, true, "fetching \"test.com\" from tree", t)
	}

	v, ok, err = r.WireGet(domain.WireNameLower("\x03www\x04test\x03com\x00"))
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		assertValue(v, ok, 5, true, "fetching \"www.test.com\" from tree", t)
	}

	v, ok, err = r.WireGet(domain.WireNameLower("\x02ns\x04test\x03com\x00"))
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		assertValue(v, ok, 2, true, "fetching \"ns.test.com\" from tree", t)
	}

	v, ok, err = r.WireGet(domain.WireNameLower("\x04test\x03org\x00"))
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	} else {
		assertValue(v, ok, 0, false, "fetching \"test.org\" from tree", t)
	}
}

func TestDeleteSubdomains(t *testing.T) {
	var r *Node

	r, ok := r.DeleteSubdomains("test.com")
	if ok {
		t.Error("Expected no deletion from empty tree but got deleted something")
	}

	r = r.Insert("com", 1)
	r = r.Insert("test.com", 2)
	r = r.Insert("test.net", 3)
	r = r.Insert("example.com", 4)
	r = r.Insert("www.test.com", 5)
	r = r.Insert("www.test.org", 6)

	r, ok = r.DeleteSubdomains("ns.test.com")
	if ok {
		t.Error("Expected \"ns.test.com\" to be not deleted as it's absent in the tree")
	}

	r, ok = r.DeleteSubdomains("test.com")
	if !ok {
		t.Error("Expected \"test.com\" to be deleted")
	}

	r, ok = r.DeleteSubdomains("www.test.com")
	if ok {
		t.Error("Expected \"www.test.com\" to be not deleted as it should be deleted with \"test.com\"")
	}

	r, ok = r.DeleteSubdomains("com")
	if !ok {
		t.Error("Expected \"com\" to be deleted")
	}

	assertTree(r, "tree with no \"com\"", t,
		"\"test.net\": 3\n",
		"\"www.test.org\": 6\n")

	r, ok = r.DeleteSubdomains("test.net")
	if !ok {
		t.Error("Expected \"test.net\" to be deleted")
	}

	r, ok = r.DeleteSubdomains("")
	if !ok {
		t.Error("Expected not empty tree to be cleaned up")
	}

	r, ok = r.DeleteSubdomains("")
	if ok {
		t.Error("Expected nothing to clean up from empty tree")
	}

	r = r.Insert("com", 1)
	r = r.Insert("test.com", 2)
	r = r.Insert("test.net", 3)
	r = r.Insert("example.com", 4)
	r = r.Insert("www.test.com", 5)
	r = r.Insert("www.test.org", 6)

	r, ok = r.DeleteSubdomains("WwW.tEsT.cOm")
	if !ok {
		t.Error("Expected \"WwW.tEsT.cOm\" to be deleted")
	}
}

func TestDelete(t *testing.T) {
	var r *Node

	r, ok := r.Delete("test.com")
	if ok {
		t.Error("Expected no deletion from empty tree but got deleted something")
	}

	r = r.Insert("com", 1)
	r = r.Insert("test.com", 2)
	r = r.Insert("test.net", 3)
	r = r.Insert("example.com", 4)
	r = r.Insert("www.test.com", 5)
	r = r.Insert("www.test.org", 6)

	r, ok = r.Delete("ns.test.com")
	if ok {
		t.Error("Expected \"ns.test.com\" to be not deleted as it's absent in the tree")
	}

	r, ok = r.Delete("test.com")
	if !ok {
		t.Error("Expected \"test.com\" to be deleted")
	}

	r, ok = r.Delete("www.test.com")
	if !ok {
		t.Error("Expected \"www.test.com\" to be deleted")
	}

	r, ok = r.Delete("com")
	if !ok {
		t.Error("Expected \"com\" to be deleted")
	}

	assertTree(r, "tree", t,
		"\"example.com\": 4\n",
		"\"test.net\": 3\n",
		"\"www.test.org\": 6\n")

	r, ok = r.Delete("test.net")
	if !ok {
		t.Error("Expected \"test.net\" to be deleted")
	}

	r, ok = r.Delete("")
	if ok {
		t.Error("Expected nothing to clean up from tree which hasn't set value for root domain")
	}

	r = r.Insert("", 1)
	assertTree(r, "tree", t,
		"\"\": 1\n",
		"\"example.com\": 4\n",
		"\"www.test.org\": 6\n")

	r, ok = r.Delete("")
	if !ok {
		t.Error("Expected root domain to be deleted")
	}

	assertTree(r, "tree", t,
		"\"example.com\": 4\n",
		"\"www.test.org\": 6\n")

	r = r.Insert("", 1)
	r, ok = r.Delete("example.com")
	if !ok {
		t.Error("Expected \"example.com\" to be deleted")
	}

	r, ok = r.Delete("www.test.org")
	if !ok {
		t.Error("Expected \"www.test.org\" to be deleted")
	}

	r, ok = r.Delete("")
	if !ok {
		t.Error("Expected root domain to be deleted")
	}

	r, ok = r.Delete("")
	if ok {
		t.Error("Expected nothing to be deleted from empty tree")
	}

	r = r.Insert("com", 1)
	r = r.Insert("test.com", 2)
	r = r.Insert("test.net", 3)
	r = r.Insert("example.com", 4)
	r = r.Insert("www.test.com", 5)
	r = r.Insert("www.test.org", 6)

	r, ok = r.Delete("WwW.tEsT.cOm")
	if !ok {
		t.Error("Expected \"WwW.tEsT.cOm\" to be deleted")
	}
}

func assertTree(r *Node, desc string, t *testing.T, e ...string) {
	pairs := []string{}
	for p := range r.Enumerate() {
		pairs = append(pairs, fmt.Sprintf("%q: %d\n", p.Key, p.Value))
	}

	ctx := difflib.ContextDiff{
		A:        e,
		B:        pairs,
		FromFile: "Expected",
		ToFile:   "Got"}

	diff, err := difflib.GetContextDiffString(ctx)
	if err != nil {
		panic(fmt.Errorf("can't compare \"%s\": %s", desc, err))
	}

	if len(diff) > 0 {
		t.Errorf("\"%s\" doesn't match:\n%s", desc, diff)
	}
}

func assertValue(v uint8, vok bool, e uint8, eok bool, desc string, t *testing.T) {
	if eok {
		if vok {
			if v != e {
				t.Errorf("Expected %d for %s but got %d", e, desc, v)
			}
		} else {
			t.Errorf("Expected %d for %s but got nothing", e, desc)
		}
	} else {
		if vok {
			t.Errorf("Expected no value for %s but got %d", desc, v)
		}
	}
}