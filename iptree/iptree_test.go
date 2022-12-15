package iptree

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/infobloxopen/go-trees/numtree"
)

func TestInsertNet(t *testing.T) {
	r := NewTree()

	newR := r.InsertNet(nil, "test")
	if newR != r {
		t.Errorf("Expected no changes inserting nil network but got:\n%s\n", newR.root32.Dot())
	}

	newR = r.InsertNet(&net.IPNet{IP: nil, Mask: nil}, "test")
	if newR != r {
		t.Errorf("Expected no changes inserting invalid network but got:\n%s\n", newR.root32.Dot())
	}

	_, n, _ := net.ParseCIDR("192.0.2.0/24")
	newR = r.InsertNet(n, "test")
	if newR == r {
		t.Errorf("Expected new root after insertion of new IPv4 address but got previous")
	} else {
		assertTree32Node(newR, 0xc0000200, 24, "test", "tree with single IPv4 address inserted", t)
	}

	_, n, _ = net.ParseCIDR("2001:db8::/32")
	r1 := r.InsertNet(n, "test")
	if r1 == r {
		t.Errorf("Expected new root after insertion of new IPv6 address but got previous")
	} else {
		assertTree64Node(r1, 0x20010db800000000, 32, 0x0, 0, "test",
			"tree with single IPv6 address inserted", t)
	}

	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	r2 := r1.InsertNet(n, "test 1")
	if r2 == r1 {
		t.Errorf("Expected new root after insertion of second IPv6 address but got previous")
	} else {
		assertTree64Node(r2, 0x20010db800000000, 64, 0x000000ff00000000, 32, "test 1",
			"tree with second IPv6 address inserted", t)
	}

	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:fe::/96")
	r3 := r2.InsertNet(n, "test 2")
	if r3 == r1 {
		t.Errorf("Expected new root after insertion of third IPv6 address but got previous")
	} else {
		assertTree64Node(r3, 0x20010db800000000, 64, 0x000000fe00000000, 32, "test 2",
			"tree with third IPv6 address inserted", t)
	}

	invR := NewTree()
	invR.root64 = invR.root64.Insert(0x20010db800000000, 64, "test")
	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	assertPanic(func() { invR.InsertNet(n, "panic") }, "inserting to invalid IPv6 tree", t)
}

func TestInplaceInsertNet(t *testing.T) {
	r := NewTree()

	r.InplaceInsertNet(nil, "test")
	if r.root32 != nil || r.root64 != nil {
		t.Error("Expected empty tree after inserting nil network")
	}

	r.InplaceInsertNet(&net.IPNet{IP: nil, Mask: nil}, "test")
	if r.root32 != nil || r.root64 != nil {
		t.Error("Expected empty tree after inserting invalid network")
	}

	_, n, _ := net.ParseCIDR("192.0.2.0/24")
	r.InplaceInsertNet(n, "test")
	if r.root32 == nil {
		t.Error("Expected some data in 32-bit tree")
	} else {
		assertTree32Node(r, 0xc0000200, 24, "test", "tree with single IPv4 address inserted", t)
	}

	_, n, _ = net.ParseCIDR("2001:db8::/32")
	r.InplaceInsertNet(n, "test")
	if r.root64 == nil {
		t.Error("Expected some data in 64-bit tree")
	} else {
		assertTree64Node(r, 0x20010db800000000, 32, 0x0, 0, "test",
			"tree with single IPv6 address inserted", t)
	}

	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	r.InplaceInsertNet(n, "test 1")
	if r.root64 == nil {
		t.Error("Expected some data in 64-bit tree")
	} else {
		assertTree64Node(r, 0x20010db800000000, 64, 0x000000ff00000000, 32, "test 1",
			"tree with second IPv6 address inserted", t)
	}

	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:fe::/96")
	r.InplaceInsertNet(n, "test 2")
	if r.root64 == nil {
		t.Error("Expected some data in 64-bit tree")
	} else {
		assertTree64Node(r, 0x20010db800000000, 64, 0x000000fe00000000, 32, "test 2",
			"tree with third IPv6 address inserted", t)
	}

	invR := NewTree()
	invR.root64 = invR.root64.Insert(0x20010db800000000, 64, "test")
	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	assertPanic(func() { invR.InplaceInsertNet(n, "panic") }, "inserting to invalid IPv6 tree", t)
}

func (p Pair) String() string {
	s, ok := p.Value.(string)
	if ok {
		return fmt.Sprintf("%s: %q", p.Key, s)
	}

	return fmt.Sprintf("%s: %T (%#v)", p.Key, p.Value, p.Value)
}

func TestEnumerate(t *testing.T) {
	var r *Tree

	for p := range r.Enumerate() {
		t.Errorf("Expected no nodes in empty tree but got at least one: %s", p)
		break
	}

	r = NewTree()

	_, n, _ := net.ParseCIDR("192.0.2.0/24")
	r = r.InsertNet(n, "test 1")

	_, n, _ = net.ParseCIDR("2001:db8::/32")
	r = r.InsertNet(n, "test 2.1")

	_, n, _ = net.ParseCIDR("2001:db8:1::/48")
	r = r.InsertNet(n, "test 2.2")

	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	r = r.InsertNet(n, "test 3")

	items := []string{}
	for p := range r.Enumerate() {
		items = append(items, p.String())
	}

	s := strings.Join(items, ", ")
	e := "192.0.2.0/24: \"test 1\", " +
		"2001:db8::/32: \"test 2.1\", " +
		"2001:db8::ff:0:0/96: \"test 3\", " +
		"2001:db8:1::/48: \"test 2.2\""
	if s != e {
		t.Errorf("Expected following nodes %q but got %q", e, s)
	}
}

func TestUpdateDescendants(t *testing.T) {
	type action struct {
		network        string
		value          string
		expectedResult bool
	}
	testCases := []struct {
		actions        []action
		tree           string
		target         string
		updateCallback func(Pair) (interface{}, bool)
		updatedTree    string
	}{
		{
			actions: []action{
				{"10.0.0.0/24", "0", false},
				{"10.0.0.0/28", "1", false},
				{"10.0.0.0/32", "2", false},
				{"10.0.0.0/16", "3", true},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"10.0.0.0/16 - (\"3\")\n" +
				"\t10.0.0.0/24 - (\"0\")\n" +
				"\t\t10.0.0.0/28 - (\"1\")\n" +
				"\t\t\t10.0.0.0/32 - (\"2\")\n" +
				"\n==== 128 bit ====\n" +
				"nil\n",
			target: "10.0.0.0/16",
			updateCallback: func(p Pair) (interface{}, bool) {
				newValues := map[string]string{
					"10.0.0.0/24": "1.0",
					"10.0.0.0/28": "1.1",
					"10.0.0.0/32": "1.2",
					"10.0.0.0/16": "1.3",
				}
				if v, ok := newValues[p.Key.String()]; ok {
					return v, true
				}
				panic("new value not provided")
				return nil, false
			},
			updatedTree: "" +
				"\n==== 32 bit ====\n" +
				"10.0.0.0/16 - (\"3\")\n" +
				"\t10.0.0.0/24 - (\"1.0\")\n" +
				"\t\t10.0.0.0/28 - (\"1.1\")\n" +
				"\t\t\t10.0.0.0/32 - (\"1.2\")\n" +
				"\n==== 128 bit ====\n" +
				"nil\n",
		},
		{
			actions: []action{
				{"192.168.0.0/24", "0", false},
				{"192.168.0.0/28", "1", false},
				{"10.0.0.0/16", "5", false},
				{"192.168.0.0/32", "2", false},
				{"192.168.0.0/16", "3", true},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"0.0.0.0/0 - (<nil>)\n" +
				"\t10.0.0.0/16 - (\"5\")\n" +
				"\t192.168.0.0/16 - (\"3\")\n" +
				"\t\t192.168.0.0/24 - (\"0\")\n" +
				"\t\t\t192.168.0.0/28 - (\"1\")\n" +
				"\t\t\t\t192.168.0.0/32 - (\"2\")\n" +
				"\n==== 128 bit ====\n" +
				"nil\n",
			target: "192.168.0.0/24",
			updateCallback: func(p Pair) (interface{}, bool) {
				newValues := map[string]string{
					"192.168.0.0/28": "1.1",
					"192.168.0.0/32": "1.2",
				}
				if v, ok := newValues[p.Key.String()]; ok {
					return v, true
				}
				panic("new value not provided")
				return nil, false
			},
			updatedTree: "" +
				"\n==== 32 bit ====\n" +
				"0.0.0.0/0 - (<nil>)\n" +
				"\t10.0.0.0/16 - (\"5\")\n" +
				"\t192.168.0.0/16 - (\"3\")\n" +
				"\t\t192.168.0.0/24 - (\"0\")\n" +
				"\t\t\t192.168.0.0/28 - (\"1.1\")\n" +
				"\t\t\t\t192.168.0.0/32 - (\"1.2\")\n" +
				"\n==== 128 bit ====\n" +
				"nil\n",
		},
		{
			actions: []action{
				{"192.168.0.0/16", "3", false},
				{"192.168.0.0/32", "2", false},
				{"192.168.0.0/28", "1", true},
				{"192.168.0.0/24", "0", true},
				{"10.0.0.0/16", "5", false},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"0.0.0.0/0 - (<nil>)\n" +
				"\t10.0.0.0/16 - (\"5\")\n" +
				"\t192.168.0.0/16 - (\"3\")\n" +
				"\t\t192.168.0.0/24 - (\"0\")\n" +
				"\t\t\t192.168.0.0/28 - (\"1\")\n" +
				"\t\t\t\t192.168.0.0/32 - (\"2\")\n" +
				"\n==== 128 bit ====\n" +
				"nil\n",
			target: "10.0.0.0/16",
			updateCallback: func(p Pair) (interface{}, bool) {
				newValues := map[string]string{}
				if v, ok := newValues[p.Key.String()]; ok {
					return v, true
				}
				panic("new value not provided")
				return nil, false
			},
			updatedTree: "" +
				"\n==== 32 bit ====\n" +
				"0.0.0.0/0 - (<nil>)\n" +
				"\t10.0.0.0/16 - (\"5\")\n" +
				"\t192.168.0.0/16 - (\"3\")\n" +
				"\t\t192.168.0.0/24 - (\"0\")\n" +
				"\t\t\t192.168.0.0/28 - (\"1\")\n" +
				"\t\t\t\t192.168.0.0/32 - (\"2\")\n" +
				"\n==== 128 bit ====\n" +
				"nil\n",
		},
		{
			actions: []action{
				{"192.168.0.0/16", "3", false},
				{"192.168.0.0/32", "2", false},
				{"192.168.0.0/28", "1", true},
				{"192.168.0.0/24", "0", true},
				{"10.0.0.0/16", "5", false},
				{"192.168.0.0/16", "3", true},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"0.0.0.0/0 - (<nil>)\n" +
				"\t10.0.0.0/16 - (\"5\")\n" +
				"\t192.168.0.0/16 - (\"3\")\n" +
				"\t\t192.168.0.0/24 - (\"0\")\n" +
				"\t\t\t192.168.0.0/28 - (\"1\")\n" +
				"\t\t\t\t192.168.0.0/32 - (\"2\")\n" +
				"\n==== 128 bit ====\n" +
				"nil\n",
			target: "192.168.0.0/32",
			updateCallback: func(p Pair) (interface{}, bool) {
				newValues := map[string]string{}
				if v, ok := newValues[p.Key.String()]; ok {
					return v, true
				}
				panic("new value not provided")
				return nil, false
			},
			updatedTree: "" +
				"\n==== 32 bit ====\n" +
				"0.0.0.0/0 - (<nil>)\n" +
				"\t10.0.0.0/16 - (\"5\")\n" +
				"\t192.168.0.0/16 - (\"3\")\n" +
				"\t\t192.168.0.0/24 - (\"0\")\n" +
				"\t\t\t192.168.0.0/28 - (\"1\")\n" +
				"\t\t\t\t192.168.0.0/32 - (\"2\")\n" +
				"\n==== 128 bit ====\n" +
				"nil\n",
		},

		// IPv6
		{
			actions: []action{
				{"2001:4860:4860::/48", "test 5", false},
				{"2001:4860:4860::/56", "test 4", false},
				{"2001:4860:4860::/64", "test 3", false},
				{"2001:4860:4860::/92", "test 2", false},
				{"2001:4860:4860::/128", "test 1", false},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 3\")\n" +
				"\t\t\t2001:4860:4860::/92 (\"test 2\")\n" +
				"\t\t\t\t2001:4860:4860::/128 (\"test 1\")\n",
			target: "2001:4860:4860::/48",
			updateCallback: func(p Pair) (interface{}, bool) {
				newValues := map[string]string{
					"2001:4860:4860::/56":  "test 1.4",
					"2001:4860:4860::/64":  "test 1.3",
					"2001:4860:4860::/92":  "test 1.2",
					"2001:4860:4860::/128": "test 1.1",
				}
				if v, ok := newValues[p.Key.String()]; ok {
					return v, true
				}
				fmt.Printf("DEBUG: KEY (%s) -- VALUE (%#v)\n", p.Key, p.Value)
				panic("new value not provided")
				return nil, false
			},
			updatedTree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 1.4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 1.3\")\n" +
				"\t\t\t2001:4860:4860::/92 (\"test 1.2\")\n" +
				"\t\t\t\t2001:4860:4860::/128 (\"test 1.1\")\n",
		},
		{
			actions: []action{
				{"2001:4860:4860::/48", "test 5", false},
				{"2001:4860:4860::/56", "test 4", false},
				{"2001:4860:4860::/64", "test 3", false},
				{"2001:4860:4860::/92", "test 2", false},
				{"2001:4860:4860::/128", "test 1", false},
				{"2001:4860:4860::/64", "test 3", true},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 3\")\n" +
				"\t\t\t2001:4860:4860::/92 (\"test 2\")\n" +
				"\t\t\t\t2001:4860:4860::/128 (\"test 1\")\n",
			target: "2001:4860:4860::/128",
			updateCallback: func(p Pair) (interface{}, bool) {
				newValues := map[string]string{}
				if v, ok := newValues[p.Key.String()]; ok {
					return v, true
				}
				panic("new value not provided")
				return nil, false
			},
			updatedTree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 3\")\n" +
				"\t\t\t2001:4860:4860::/92 (\"test 2\")\n" +
				"\t\t\t\t2001:4860:4860::/128 (\"test 1\")\n",
		},
		{
			actions: []action{
				{"2001:4860:4860::/128", "test 1", false},
				{"2001:4860:4860::/127", "test 6", true},
				{"2001:4860:4860::/92", "test 2", true},
				{"2001:4860:4860::/64", "test 3", true},
				{"2001:4860:4860::/56", "test 4", true},
				{"2001:4860:4860::/48", "test 5", true},
				{"2001:4860:4860::/127", "test 6", true},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 3\")\n" +
				"\t\t\t2001:4860:4860::/92 (\"test 2\")\n" +
				"\t\t\t\t2001:4860:4860::/127 (\"test 6\")\n" +
				"\t\t\t\t\t2001:4860:4860::/128 (\"test 1\")\n",
			target: "2001:4860:4860::/92",
			updateCallback: func(p Pair) (interface{}, bool) {
				newValues := map[string]string{
					"2001:4860:4860::/127": "test 1.6",
					"2001:4860:4860::/128": "test 1.1",
				}
				if v, ok := newValues[p.Key.String()]; ok {
					return v, true
				}
				panic("new value not provided")
				return nil, false
			},
			updatedTree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 3\")\n" +
				"\t\t\t2001:4860:4860::/92 (\"test 2\")\n" +
				"\t\t\t\t2001:4860:4860::/127 (\"test 1.6\")\n" +
				"\t\t\t\t\t2001:4860:4860::/128 (\"test 1.1\")\n",
		},
		{
			actions: []action{
				{"2001:4860:4860::/56", "test 4", false},
				{"2001:4860:4860::/128", "test 1", false},
				{"2001:4860:4860::001f:ffff:ffff/98", "test 7", false},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/56 (\"test 4\")\n" +
				"\t2001:4860:4860::/91 (<nil>)\n" +
				"\t\t2001:4860:4860::/128 (\"test 1\")\n" +
				"\t\t2001:4860:4860::1f:c000:0/98 (\"test 7\")\n",
			target: "2001:4860:4860::/56",
			updateCallback: func(p Pair) (interface{}, bool) {
				newValues := map[string]string{
					"2001:4860:4860::/128":         "test 1.1",
					"2001:4860:4860::1f:c000:0/98": "test 1.7",
				}
				if v, ok := newValues[p.Key.String()]; ok {
					return v, true
				}
				panic("new value not provided")
				return nil, false
			},
			updatedTree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/56 (\"test 4\")\n" +
				"\t2001:4860:4860::/91 (<nil>)\n" +
				"\t\t2001:4860:4860::/128 (\"test 1.1\")\n" +
				"\t\t2001:4860:4860::1f:c000:0/98 (\"test 1.7\")\n",
		},
		{
			actions: []action{
				{"2001:4860:4860::/56", "test 4", false},
				{"2001:4860:4860::/128", "test 1", false},
				{"2001:4860:4860::001f:ffff:ffff/98", "test 7", false},
				{"2001:4860:4860:0:1:0:ffff:ffff/98", "test 8", false},
				{"2001:4860:4860::/92", "test 2", true},
				{"2001:4860:4860::001f:ffff:ffff/128", "test 6", false},
				{"2001:4860:4860::/48", "test 5", true},
				{"2001:4860:4860::/64", "test 3", true},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 3\")\n" +
				"\t\t\t2001:4860:4860::/79 (<nil>)\n" +
				"\t\t\t\t2001:4860:4860::/91 (<nil>)\n" +
				"\t\t\t\t\t2001:4860:4860::/92 (\"test 2\")\n" +
				"\t\t\t\t\t\t2001:4860:4860::/128 (\"test 1\")\n" +
				"\t\t\t\t\t2001:4860:4860::1f:c000:0/98 (\"test 7\")\n" +
				"\t\t\t\t\t\t2001:4860:4860::1f:ffff:ffff/128 (\"test 6\")\n" +
				"\t\t\t\t2001:4860:4860:0:1:0:c000:0/98 (\"test 8\")\n",
			target: "2001:4860:4860::/91",
			updateCallback: func(p Pair) (interface{}, bool) {
				newValues := map[string]string{}
				if v, ok := newValues[p.Key.String()]; ok {
					return v, true
				}
				panic("new value not provided")
				return nil, false
			},
			updatedTree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 3\")\n" +
				"\t\t\t2001:4860:4860::/79 (<nil>)\n" +
				"\t\t\t\t2001:4860:4860::/91 (<nil>)\n" +
				"\t\t\t\t\t2001:4860:4860::/92 (\"test 2\")\n" +
				"\t\t\t\t\t\t2001:4860:4860::/128 (\"test 1\")\n" +
				"\t\t\t\t\t2001:4860:4860::1f:c000:0/98 (\"test 7\")\n" +
				"\t\t\t\t\t\t2001:4860:4860::1f:ffff:ffff/128 (\"test 6\")\n" +
				"\t\t\t\t2001:4860:4860:0:1:0:c000:0/98 (\"test 8\")\n",
		},
		{
			actions: []action{
				{"2001:4860:4860::/56", "test 4", false},
				{"2001:4860:4860::/128", "test 1", false},
				{"2001:4860:4860::001f:ffff:ffff/98", "test 7", false},
				{"2001:4860:4860:0:1:0:ffff:ffff/98", "test 8", false},
				{"2001:4860:4860::/92", "test 2", true},
				{"2001:4860:4860::001f:ffff:ffff/128", "test 6", false},
				{"2001:4860:4860::/48", "test 5", true},
				{"2001:4860:4860::/64", "test 3", true},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 3\")\n" +
				"\t\t\t2001:4860:4860::/79 (<nil>)\n" +
				"\t\t\t\t2001:4860:4860::/91 (<nil>)\n" +
				"\t\t\t\t\t2001:4860:4860::/92 (\"test 2\")\n" +
				"\t\t\t\t\t\t2001:4860:4860::/128 (\"test 1\")\n" +
				"\t\t\t\t\t2001:4860:4860::1f:c000:0/98 (\"test 7\")\n" +
				"\t\t\t\t\t\t2001:4860:4860::1f:ffff:ffff/128 (\"test 6\")\n" +
				"\t\t\t\t2001:4860:4860:0:1:0:c000:0/98 (\"test 8\")\n",
			target: "2001:4860:4860::/92",
			updateCallback: func(p Pair) (interface{}, bool) {
				newValues := map[string]string{
					"2001:4860:4860::/128": "test 1.1",
				}
				if v, ok := newValues[p.Key.String()]; ok {
					return v, true
				}
				panic("new value not provided")
				return nil, false
			},
			updatedTree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 3\")\n" +
				"\t\t\t2001:4860:4860::/79 (<nil>)\n" +
				"\t\t\t\t2001:4860:4860::/91 (<nil>)\n" +
				"\t\t\t\t\t2001:4860:4860::/92 (\"test 2\")\n" +
				"\t\t\t\t\t\t2001:4860:4860::/128 (\"test 1.1\")\n" +
				"\t\t\t\t\t2001:4860:4860::1f:c000:0/98 (\"test 7\")\n" +
				"\t\t\t\t\t\t2001:4860:4860::1f:ffff:ffff/128 (\"test 6\")\n" +
				"\t\t\t\t2001:4860:4860:0:1:0:c000:0/98 (\"test 8\")\n",
		},
		{
			actions: []action{
				{"2001:4860:4860::/56", "test 4", false},
				{"2001:4860:4860::/128", "test 1", false},
				{"2001:4860:4860::001f:ffff:ffff/98", "test 7", false},
				{"2001:4860:4860:0:1:0:ffff:ffff/98", "test 8", false},
				{"2001:4860:4860::/92", "test 2", true},
				{"2001:4860:4860::001f:ffff:ffff/128", "test 6", false},
				{"2001:4860:4860::/48", "test 5", true},
				{"2001:4860:4860::/64", "test 3", true},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 3\")\n" +
				"\t\t\t2001:4860:4860::/79 (<nil>)\n" +
				"\t\t\t\t2001:4860:4860::/91 (<nil>)\n" +
				"\t\t\t\t\t2001:4860:4860::/92 (\"test 2\")\n" +
				"\t\t\t\t\t\t2001:4860:4860::/128 (\"test 1\")\n" +
				"\t\t\t\t\t2001:4860:4860::1f:c000:0/98 (\"test 7\")\n" +
				"\t\t\t\t\t\t2001:4860:4860::1f:ffff:ffff/128 (\"test 6\")\n" +
				"\t\t\t\t2001:4860:4860:0:1:0:c000:0/98 (\"test 8\")\n",
			target: "2001:4860:4860::/64",
			updateCallback: func(p Pair) (interface{}, bool) {
				newValues := map[string]string{
					"2001:4860:4860::/92":              "test 1.2",
					"2001:4860:4860::/128":             "test 1.1",
					"2001:4860:4860::1f:c000:0/98":     "test 1.7",
					"2001:4860:4860::1f:ffff:ffff/128": "test 1.6",
					"2001:4860:4860:0:1:0:c000:0/98":   "test 1.8",
				}
				if v, ok := newValues[p.Key.String()]; ok {
					return v, true
				}
				panic("new value not provided")
				return nil, false
			},
			updatedTree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 3\")\n" +
				"\t\t\t2001:4860:4860::/79 (<nil>)\n" +
				"\t\t\t\t2001:4860:4860::/91 (<nil>)\n" +
				"\t\t\t\t\t2001:4860:4860::/92 (\"test 1.2\")\n" +
				"\t\t\t\t\t\t2001:4860:4860::/128 (\"test 1.1\")\n" +
				"\t\t\t\t\t2001:4860:4860::1f:c000:0/98 (\"test 1.7\")\n" +
				"\t\t\t\t\t\t2001:4860:4860::1f:ffff:ffff/128 (\"test 1.6\")\n" +
				"\t\t\t\t2001:4860:4860:0:1:0:c000:0/98 (\"test 1.8\")\n",
		},
	}

	parseCIDR := func(s string) *net.IPNet {
		_, n, _ := net.ParseCIDR(s)
		return n
	}

	for i, tc := range testCases {
		tc, i := tc, i+1
		t.Run(fmt.Sprintf("Test %d: UpdateDescendants", i), func(t *testing.T) {
			var r *Tree
			r = NewTree()
			for _, c := range tc.actions {
				actual := r.InplaceInsertNetCheckChildren(parseCIDR(c.network), c.value)
				if actual != c.expectedResult {
					t.Errorf(
						"Expected result does not match for new node entry (%s - %s)\n\t\t expected: %v\n\t\t   actual: %v",
						c.network, c.value, c.expectedResult, actual,
					)
				}
			}

			if tc.tree != "" {
				if r.String() != tc.tree {
					t.Errorf("Tree representation did not match\nexpected:\n%s\n\n  actual:\n%s\n", tc.tree, r.String())
				}
			}

			r.UpdateDescendants(parseCIDR(tc.target), tc.updateCallback)

			if tc.updatedTree != "" {
				if r.String() != tc.updatedTree {
					t.Errorf("Tree representation did not match\nexpected:\n%s\n\n  actual:\n%s\n", tc.updatedTree, r.String())
				}
			}
		})
	}
}

func TestInplaceInsertNet2(t *testing.T) {
	type action struct {
		network        string
		value          string
		expectedResult bool
	}
	testCases := []struct {
		actions []action
		tree    string
	}{
		{
			actions: []action{
				{"10.0.0.0/24", "0", false},
				{"10.0.0.0/28", "1", false},
				{"10.0.0.0/32", "2", false},
				{"10.0.0.0/16", "3", true},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"10.0.0.0/16 - (\"3\")\n" +
				"\t10.0.0.0/24 - (\"0\")\n" +
				"\t\t10.0.0.0/28 - (\"1\")\n" +
				"\t\t\t10.0.0.0/32 - (\"2\")\n" +
				"\n==== 128 bit ====\n" +
				"nil\n",
		},
		{
			actions: []action{
				{"192.168.0.0/24", "0", false},
				{"192.168.0.0/28", "1", false},
				{"10.0.0.0/16", "5", false},
				{"192.168.0.0/32", "2", false},
				{"192.168.0.0/16", "3", true},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"0.0.0.0/0 - (<nil>)\n" +
				"\t10.0.0.0/16 - (\"5\")\n" +
				"\t192.168.0.0/16 - (\"3\")\n" +
				"\t\t192.168.0.0/24 - (\"0\")\n" +
				"\t\t\t192.168.0.0/28 - (\"1\")\n" +
				"\t\t\t\t192.168.0.0/32 - (\"2\")\n" +
				"\n==== 128 bit ====\n" +
				"nil\n",
		},
		{
			actions: []action{
				{"192.168.0.0/16", "3", false},
				{"192.168.0.0/32", "2", false},
				{"192.168.0.0/28", "1", true},
				{"192.168.0.0/24", "0", true},
				{"10.0.0.0/16", "5", false},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"0.0.0.0/0 - (<nil>)\n" +
				"\t10.0.0.0/16 - (\"5\")\n" +
				"\t192.168.0.0/16 - (\"3\")\n" +
				"\t\t192.168.0.0/24 - (\"0\")\n" +
				"\t\t\t192.168.0.0/28 - (\"1\")\n" +
				"\t\t\t\t192.168.0.0/32 - (\"2\")\n" +
				"\n==== 128 bit ====\n" +
				"nil\n",
		},
		{
			actions: []action{
				{"192.168.0.0/16", "3", false},
				{"192.168.0.0/32", "2", false},
				{"192.168.0.0/28", "1", true},
				{"192.168.0.0/24", "0", true},
				{"10.0.0.0/16", "5", false},
				{"192.168.0.0/16", "3", true},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"0.0.0.0/0 - (<nil>)\n" +
				"\t10.0.0.0/16 - (\"5\")\n" +
				"\t192.168.0.0/16 - (\"3\")\n" +
				"\t\t192.168.0.0/24 - (\"0\")\n" +
				"\t\t\t192.168.0.0/28 - (\"1\")\n" +
				"\t\t\t\t192.168.0.0/32 - (\"2\")\n" +
				"\n==== 128 bit ====\n" +
				"nil\n",
		},
		{
			actions: []action{
				{"192.168.0.0/16", "1", false},
				{"10.0.0.0/16", "2", false},
				{"172.0.0.0/16", "3", false},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"0.0.0.0/0 - (<nil>)\n" +
				"\t10.0.0.0/16 - (\"2\")\n" +
				"\t128.0.0.0/1 - (<nil>)\n" +
				"\t\t172.0.0.0/16 - (\"3\")\n" +
				"\t\t192.168.0.0/16 - (\"1\")\n" +
				"\n==== 128 bit ====\n" +
				"nil\n",
		},

		// IPv6
		{
			actions: []action{
				{"2001:4860:4860::/48", "test 5", false},
				{"2001:4860:4860::/56", "test 4", false},
				{"2001:4860:4860::/64", "test 3", false},
				{"2001:4860:4860::/92", "test 2", false},
				{"2001:4860:4860::/128", "test 1", false},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 3\")\n" +
				"\t\t\t2001:4860:4860::/92 (\"test 2\")\n" +
				"\t\t\t\t2001:4860:4860::/128 (\"test 1\")\n",
		},
		{
			actions: []action{
				{"2001:4860:4860::/48", "test 5", false},
				{"2001:4860:4860::/56", "test 4", false},
				{"2001:4860:4860::/64", "test 3", false},
				{"2001:4860:4860::/92", "test 2", false},
				{"2001:4860:4860::/128", "test 1", false},
				{"2001:4860:4860::/64", "test 3", true},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 3\")\n" +
				"\t\t\t2001:4860:4860::/92 (\"test 2\")\n" +
				"\t\t\t\t2001:4860:4860::/128 (\"test 1\")\n",
		},
		{
			actions: []action{
				{"2001:4860:4860::/128", "test 1", false},
				{"2001:4860:4860::/127", "test 6", true},
				{"2001:4860:4860::/92", "test 2", true},
				{"2001:4860:4860::/64", "test 3", true},
				{"2001:4860:4860::/56", "test 4", true},
				{"2001:4860:4860::/48", "test 5", true},
				{"2001:4860:4860::/127", "test 6", true},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 3\")\n" +
				"\t\t\t2001:4860:4860::/92 (\"test 2\")\n" +
				"\t\t\t\t2001:4860:4860::/127 (\"test 6\")\n" +
				"\t\t\t\t\t2001:4860:4860::/128 (\"test 1\")\n",
		},
		{
			actions: []action{
				{"2001:4860:4860::/56", "test 4", false},
				{"2001:4860:4860::/128", "test 1", false},
				{"2001:4860:4860::001f:ffff:ffff/98", "test 7", false},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/56 (\"test 4\")\n" +
				"\t2001:4860:4860::/91 (<nil>)\n" +
				"\t\t2001:4860:4860::/128 (\"test 1\")\n" +
				"\t\t2001:4860:4860::1f:c000:0/98 (\"test 7\")\n",
		},
		{
			actions: []action{
				{"2001:4860:4860::/56", "test 4", false},
				{"2001:4860:4860::/128", "test 1", false},
				{"2001:4860:4860::001f:ffff:ffff/98", "test 7", false},
				{"2001:4860:4860:0:1:0:ffff:ffff/98", "test 8", false},
				{"2001:4860:4860::/92", "test 2", true},
				{"2001:4860:4860::001f:ffff:ffff/128", "test 6", false},
				{"2001:4860:4860::/48", "test 5", true},
				{"2001:4860:4860::/64", "test 3", true},
			},
			tree: "" +
				"\n==== 32 bit ====\n" +
				"nil\n" +
				"\n==== 128 bit ====\n" +
				"2001:4860:4860::/48 (\"test 5\")\n" +
				"\t2001:4860:4860::/56 (\"test 4\")\n" +
				"\t\t2001:4860:4860::/64 (\"test 3\")\n" +
				"\t\t\t2001:4860:4860::/79 (<nil>)\n" +
				"\t\t\t\t2001:4860:4860::/91 (<nil>)\n" +
				"\t\t\t\t\t2001:4860:4860::/92 (\"test 2\")\n" +
				"\t\t\t\t\t\t2001:4860:4860::/128 (\"test 1\")\n" +
				"\t\t\t\t\t2001:4860:4860::1f:c000:0/98 (\"test 7\")\n" +
				"\t\t\t\t\t\t2001:4860:4860::1f:ffff:ffff/128 (\"test 6\")\n" +
				"\t\t\t\t2001:4860:4860:0:1:0:c000:0/98 (\"test 8\")\n" +
				"",
		},
	}

	parseCIDR := func(s string) *net.IPNet {
		_, n, _ := net.ParseCIDR(s)
		return n
	}

	for i, tc := range testCases {
		tc, i := tc, i+1
		t.Run(fmt.Sprintf("Test %d: InplaceInsertNetCheckChildren", i), func(t *testing.T) {
			var r *Tree
			r = NewTree()
			for _, c := range tc.actions {
				actual := r.InplaceInsertNetCheckChildren(parseCIDR(c.network), c.value)
				if actual != c.expectedResult {
					t.Errorf(
						"Expected result does not match for new node entry (%s - %s)\n\t\t expected: %v\n\t\t   actual: %v",
						c.network, c.value, c.expectedResult, actual,
					)
				}
			}

			if tc.tree != "" {
				if r.String() != tc.tree {
					t.Errorf("Tree representation did not match\nexpected:\n%s\n\n  actual:\n%s\n", tc.tree, r.String())
				}
			}
		})
	}
}

func TestGetByNet(t *testing.T) {
	r := NewTree()

	_, n4, _ := net.ParseCIDR("192.0.2.0/24")
	r = r.InsertNet(n4, "test 1")

	_, n6Short1, _ := net.ParseCIDR("2001:db8::/32")
	r = r.InsertNet(n6Short1, "test 2.1")

	_, n6Short2, _ := net.ParseCIDR("2001:db8:1::/48")
	r = r.InsertNet(n6Short2, "test 2.2")

	_, n6Long, _ := net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	r = r.InsertNet(n6Long, "test 3")

	v, ok := r.GetByNet(nil)
	if ok {
		t.Errorf("Expected no result for nil network but got %T (%#v)", v, v)
	}

	v, ok = r.GetByNet(&net.IPNet{IP: nil, Mask: nil})
	if ok {
		t.Errorf("Expected no result for invalid network but got %T (%#v)", v, v)
	}

	v, ok = r.GetByNet(n4)
	assertResult(v, ok, "test 1", fmt.Sprintf("%s", n4), t)

	v, ok = r.GetByNet(n6Short1)
	assertResult(v, ok, "test 2.1", fmt.Sprintf("%s", n6Short1), t)

	v, ok = r.GetByNet(n6Long)
	assertResult(v, ok, "test 3", fmt.Sprintf("%s", n6Long), t)

	_, n6, _ := net.ParseCIDR("2001:db8:1::/64")
	v, ok = r.GetByNet(n6)
	assertResult(v, ok, "test 2.2", fmt.Sprintf("%s", n6), t)

	_, n6, _ = net.ParseCIDR("2001:db8:0:0:0:fe::/96")
	v, ok = r.GetByNet(n6)
	assertResult(v, ok, "test 2.1", fmt.Sprintf("%s", n6), t)
}

func TestDeleteByNet(t *testing.T) {
	var r *Tree

	_, n4, _ := net.ParseCIDR("192.0.2.0/24")
	r, ok := r.DeleteByNet(n4)
	if ok {
		t.Errorf("Expected no deletion in empty tree but got one")
	}

	r = r.InsertNet(n4, "test 1")

	_, n6Short1, _ := net.ParseCIDR("2001:db8::/32")
	r = r.InsertNet(n6Short1, "test 2.1")

	_, n6Short2, _ := net.ParseCIDR("2001:db8:1::/48")
	r = r.InsertNet(n6Short2, "test 2.2")

	_, n6Long1, _ := net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	r = r.InsertNet(n6Long1, "test 3.1")

	_, n6Long2, _ := net.ParseCIDR("2001:db8:0:0:0:fe::/96")
	r = r.InsertNet(n6Long2, "test 3.2")

	r, ok = r.DeleteByNet(nil)
	if ok {
		t.Errorf("Expected no deletion by nil network but got one")
	}

	r, ok = r.DeleteByNet(&net.IPNet{IP: nil, Mask: nil})
	if ok {
		t.Errorf("Expected no deletion by invalid network but got one")
	}

	r, ok = r.DeleteByNet(n6Long2)
	if !ok {
		t.Errorf("Expected deletion by %s but got nothing", n6Long2)
	}

	r, ok = r.DeleteByNet(n6Long1)
	if !ok {
		t.Errorf("Expected deletion by %s but got nothing", n6Long1)
	}

	v, ok := r.root64.ExactMatch(0x20010db800000000, 64)
	if ok {
		t.Errorf("Expected no subtree node at 0x%016x, %d after deleting all long mask addresses but got %#v",
			0x20010db800000000, 64, v)
	}

	r, ok = r.DeleteByNet(n6Short2)
	if !ok {
		t.Errorf("Expected deletion by %s but got nothing", n6Short2)
	}

	r, ok = r.DeleteByNet(n6Short1)
	if !ok {
		t.Errorf("Expected deletion by %s but got nothing", n6Short1)
	}

	r, ok = r.DeleteByNet(n4)
	if !ok {
		t.Errorf("Expected deletion by %s but got nothing", n4)
	}

	if r.root32 != nil || r.root64 != nil {
		t.Errorf("Expected expected empty tree at the end but have root32: %#v and root64: %#v", r.root32, r.root64)
	}

	r.root64 = r.root64.Insert(0x20010db800000000, 64, "panic")
	assertPanic(func() { r.DeleteByNet(n6Long1) }, "deletion from invalid tree", t)
}

func TestTreeByIP(t *testing.T) {
	ip := net.ParseIP("2001:db8::1")

	var r *Tree
	r = r.InsertIP(ip, "test")
	if r == nil {
		t.Errorf("Expected some tree after insert %s but got %#v", ip, r)
	}

	v, ok := r.GetByIP(ip)
	assertResult(v, ok, "test", fmt.Sprintf("address %s", ip), t)

	r, ok = r.DeleteByIP(ip)
	if !ok {
		t.Errorf("Expected deletion by address %s but got nothing", ip)
	}

	r.InplaceInsertIP(ip, "test")
	if r.root64 == nil {
		t.Errorf("Expected some tree after inplace insert %s", ip)
	}
}

func TestIPv4NetToUint32(t *testing.T) {
	_, n, _ := net.ParseCIDR("192.0.2.0/24")
	key, bits := iPv4NetToUint32(n)
	if key != 0xc0000200 || bits != 24 {
		t.Errorf("Expected 0xc0000200, 24 pair but got 0x%08x, %d", key, bits)
	}

	n = &net.IPNet{
		IP:   net.IP{0xc, 0x00},
		Mask: net.IPMask{0xff, 0xff, 0xff, 0x00}}
	key, bits = iPv4NetToUint32(n)
	if bits >= 0 {
		t.Errorf("Expected negative number of bits for invalid IPv4 address but got 0x%08x, %d", key, bits)
	}

	n = &net.IPNet{
		IP:   net.IP{0xc, 0x00, 0x02, 0x00},
		Mask: net.IPMask{0xff, 0x00, 0xff, 0x00}}
	key, bits = iPv4NetToUint32(n)
	if bits >= 0 {
		t.Errorf("Expected negative number of bits for invalid IPv4 mask but got 0x%08x, %d", key, bits)
	}
}

func TestIPv6NetToUint64Pair(t *testing.T) {
	_, n, _ := net.ParseCIDR("2001:db8::/32")
	MSKey, MSBits, LSKey, LSBits := iPv6NetToUint64Pair(n)
	if MSKey != 0x20010db800000000 || MSBits != 32 || LSKey != 0x0 || LSBits != 0 {
		t.Errorf("Expected 0x20010db800000000, 32 and 0x0000000000000000, 0 pairs bit got 0x%016x, %d and 0x%016x, %d",
			MSKey, MSBits, LSKey, LSBits)
	}

	_, n, _ = net.ParseCIDR("2001:db8:0:0:0:ff::/96")
	MSKey, MSBits, LSKey, LSBits = iPv6NetToUint64Pair(n)
	if MSKey != 0x20010db800000000 || MSBits != 64 || LSKey != 0x000000ff00000000 || LSBits != 32 {
		t.Errorf("Expected 0x20010db800000000, 32 and 0x0000000000000000, 0 pairs bit got 0x%016x, %d and 0x%016x, %d",
			MSKey, MSBits, LSKey, LSBits)
	}

	n = &net.IPNet{
		IP: net.IP{
			0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00},
		Mask: net.IPMask{
			0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}
	MSKey, MSBits, LSKey, LSBits = iPv6NetToUint64Pair(n)
	if MSBits >= 0 {
		t.Errorf("Expected negative number of bits for invalid IPv6 address but got 0x%016x, %d and 0x%016x, %d",
			MSKey, MSBits, LSKey, LSBits)
	}

	n = &net.IPNet{
		IP: net.IP{
			0x20, 0x01, 0x0d, 0xb8, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		Mask: net.IPMask{
			0x00, 0xff, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}
	MSKey, MSBits, LSKey, LSBits = iPv6NetToUint64Pair(n)
	if MSBits >= 0 {
		t.Errorf("Expected negative number of bits for invalid IPv6 mask but got 0x%016x, %d and 0x%016x, %d",
			MSKey, MSBits, LSKey, LSBits)
	}
}

func TestNewIPNetFromIP(t *testing.T) {
	n := newIPNetFromIP(net.ParseIP("192.0.2.1"))
	if n.String() != "192.0.2.1/32" {
		t.Errorf("Expected %s for IPv4 conversion but got %s", "192.0.2.1/32", n)
	}

	n = newIPNetFromIP(net.ParseIP("2001:db8::1"))
	if n.String() != "2001:db8::1/128" {
		t.Errorf("Expected %s for IPv6 conversion but got %s", "2001:db8::1/128", n)
	}

	n = newIPNetFromIP(net.IP{0xc, 0x00})
	if n != nil {
		t.Errorf("Expected %#v for invalid IP address but got %s", nil, n)
	}
}

func assertTree32Node(r *Tree, key uint32, bits int, e, desc string, t *testing.T) {
	v, ok := r.root32.ExactMatch(key, bits)
	assertResult(v, ok, e, fmt.Sprintf("0x%08x, %d for %s", key, bits, desc), t)
}

func assertTree64Node(r *Tree, MSKey uint64, MSBits int, LSKey uint64, LSBits int, e, desc string, t *testing.T) {
	desc = fmt.Sprintf("0x%016x, %d and 0x%016x, %d for %s", MSKey, MSBits, LSKey, LSBits, desc)
	v, ok := r.root64.ExactMatch(MSKey, MSBits)
	if ok {
		if MSBits < 64 {
			assertResult(v, ok, e, desc, t)
		} else {
			r, ok := v.(subTree64)
			if ok {
				v, ok := (*numtree.Node64)(r).ExactMatch(LSKey, LSBits)
				if ok {
					assertResult(v, ok, e, desc, t)
				} else {
					t.Errorf("Expected string %q at %s but got nothing at second hop", e, desc)
				}
			} else {
				t.Errorf("Expected subTree64 at %s (first hop) but got %T (%#v)", desc, v, v)
			}
		}
	} else {
		if MSBits < 64 {
			t.Errorf("Expected string %q at %s but got nothing", e, desc)
		} else {
			t.Errorf("Expected string %q at %s but got nothing even at first hop", e, desc)
		}
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

func assertResult(v interface{}, ok bool, e, desc string, t *testing.T) {
	t.Helper()
	if ok {
		s, ok := v.(string)
		if ok {
			if s != e {
				t.Errorf("Expected string %q at %s but got %q", e, desc, s)
			}
			return
		}

		t.Errorf("Expected string %q at %s but got %T (%#v)", e, desc, v, v)
		return
	}

	t.Errorf("Expected string %q at %s but got nothing", e, desc)
}

func (t *Tree) String() string {
	var sb strings.Builder
	sb.WriteString("\n==== 32 bit ====\n")
	if t.root32 != nil {
		sb.WriteString(debugNode32(t.root32))
	} else {
		sb.WriteString("nil\n")
	}
	sb.WriteString("\n==== 128 bit ====\n")
	if t.root64 != nil {
		sb.WriteString(debugNode64(t.root64))
	} else {
		sb.WriteString("nil\n")
	}
	return sb.String()
}

func debugNode32(tree *numtree.Node32) string {
	var sb strings.Builder

	var walk func(n *numtree.Node32, indent int)
	walk = func(n *numtree.Node32, indent int) {
		sb.WriteString(fmt.Sprintf("%s%s/%d - (%#v)\n", strings.Repeat("\t", indent), unpackUint32ToIP(n.Key), n.Bits, n.Value))
		c1, c2 := n.Children()
		if c1 != nil {
			walk(c1, indent+1)
		}
		if c2 != nil {
			walk(c2, indent+1)
		}
	}
	walk(tree, 0)

	return sb.String()
}

func debugNode64(tree *numtree.Node64) string {
	var sb strings.Builder

	var walk func(n *numtree.Node64, indent int)
	var walkSubtree func(MSIP net.IP, bits int, n *numtree.Node64, indent int)

	walkSubtree = func(MSIP net.IP, bits int, n *numtree.Node64, indent int) {
		LSIP := unpackUint64ToIP(n.Key)
		mask := net.CIDRMask(numtree.Key64BitSize+int(n.Bits), iPv6Bits)
		sb.WriteString(
			fmt.Sprintf("%s%s/%d (%#v)\n",
				strings.Repeat("\t", indent), append(MSIP[0:8], LSIP...).Mask(mask), bits+int(n.Bits), n.Value,
			),
		)

		if s, ok := n.Value.(subTree64); ok {
			walkSubtree(MSIP, bits, s, indent+1)
		} else {
			c1, c2 := n.Children()
			if c1 != nil {
				walkSubtree(MSIP, bits, c1, indent+1)
			}

			if c2 != nil {
				walkSubtree(MSIP, bits, c2, indent+1)
			}
		}
	}

	walk = func(n *numtree.Node64, indent int) {
		ip := unpackUint64ToIP(n.Key)
		MSIP := append(ip, make(net.IP, 8)...)
		bits := n.Bits

		c1, c2 := n.Children()

		if isNil(n.Value) {
			if c1 != nil {
				walk(c1, indent)
			}

			if c2 != nil {
				walk(c2, indent)
			}
		} else {
			indent2 := indent
			if s, ok := n.Value.(subTree64); ok {
				walkSubtree(MSIP, int(bits), s, indent2)
			} else {
				sb.WriteString(
					fmt.Sprintf("%s%s/%d (%#v)\n",
						strings.Repeat("\t", indent), MSIP, bits, n.Value,
					),
				)
			}

			if c1 != nil {
				walk(c1, indent+1)
			}

			if c2 != nil {
				walk(c2, indent+1)
			}
		}
	}
	walk(tree, 0)

	return sb.String()
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
