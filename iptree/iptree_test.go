package iptree

import (
	"fmt"
	"net"
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