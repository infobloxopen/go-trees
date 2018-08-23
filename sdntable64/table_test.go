package sdntable64

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infobloxopen/go-trees/udomain"
)

var (
	rootDN domain.Name
	orgDN  domain.Name
	comDN  domain.Name
	netDN  domain.Name
	govDN  domain.Name
	ioDN   domain.Name
	ioZone domain.Name
)

func init() {
	s := "example.org"
	n, err := domain.MakeNameFromString(s)
	if err != nil {
		panic(fmt.Errorf("failed to make domain name from %q: %s", s, err))
	}
	orgDN = n

	s = "example.com"
	n, err = domain.MakeNameFromString(s)
	if err != nil {
		panic(fmt.Errorf("failed to make domain name from %q: %s", s, err))
	}
	comDN = n

	s = "example.net"
	n, err = domain.MakeNameFromString(s)
	if err != nil {
		panic(fmt.Errorf("failed to make domain name from %q: %s", s, err))
	}
	netDN = n

	s = "example.gov"
	n, err = domain.MakeNameFromString(s)
	if err != nil {
		panic(fmt.Errorf("failed to make domain name from %q: %s", s, err))
	}
	govDN = n

	s = "example.io"
	n, err = domain.MakeNameFromString(s)
	if err != nil {
		panic(fmt.Errorf("failed to make domain name from %q: %s", s, err))
	}
	ioDN = n

	s = "io"
	n, err = domain.MakeNameFromString(s)
	if err != nil {
		panic(fmt.Errorf("failed to make domain name from %q: %s", s, err))
	}
	ioZone = n
}

func TestTable64GetDomainNameIndex(t *testing.T) {
	assert.Equal(t,
		getDomainNameIndex(domain.MaxLabels-orgDN.GetLabelCount()),
		getDomainNameIndex(domain.MaxLabels-comDN.GetLabelCount()),
	)

	assert.Equal(t,
		getDomainNameIndex(domain.MaxLabels-comDN.GetLabelCount()),
		getDomainNameIndex(domain.MaxLabels-netDN.GetLabelCount()),
	)

	assert.Equal(t,
		getDomainNameIndex(domain.MaxLabels-netDN.GetLabelCount()),
		getDomainNameIndex(domain.MaxLabels-govDN.GetLabelCount()),
	)

	assert.Equal(t,
		getDomainNameIndex(domain.MaxLabels-govDN.GetLabelCount()),
		getDomainNameIndex(domain.MaxLabels-ioDN.GetLabelCount()),
	)
}

func TestTable64InplaceInsert(t *testing.T) {
	dnt := NewTable64()
	cnt := orgDN.GetLabelCount()
	idx := getDomainNameIndex(domain.MaxLabels-cnt) + len(orgDN.GetComparable()) - cnt

	dnt.InplaceInsert(comDN, 2)
	assert.Equal(t, []uint64{
		// E L P M A X E 1+0   M O C 1+4  v: 2
		0x454c504d41584501, 0x4d4f4341, 0x2,
	}, dnt.body[idx])

	dnt.InplaceInsert(govDN, 8)
	assert.Equal(t, []uint64{
		// E L P M A X E 1+0   M O C 1+4  v: 2
		0x454c504d41584501, 0x4d4f4341, 0x2,
		// E L P M A X E 1+0   V O G 1+4  v: 8
		0x454c504d41584501, 0x564f4741, 0x8,
	}, dnt.body[idx])

	dnt.InplaceInsert(netDN, 4)
	assert.Equal(t, []uint64{
		// E L P M A X E 1+0   M O C 1+4  v: 2
		0x454c504d41584501, 0x4d4f4341, 0x2,
		// E L P M A X E 1+0   T E N 1+4  v: 4
		0x454c504d41584501, 0x54454e41, 0x4,
		// E L P M A X E 1+0   V O G 1+4  v: 8
		0x454c504d41584501, 0x564f4741, 0x8,
	}, dnt.body[idx])

	dnt.InplaceInsert(orgDN, 1)
	assert.Equal(t, []uint64{
		// E L P M A X E 1+0   G R O 1+4  v: 1
		0x454c504d41584501, 0x47524f41, 0x1,
		// E L P M A X E 1+0   M O C 1+4  v: 2
		0x454c504d41584501, 0x4d4f4341, 0x2,
		// E L P M A X E 1+0   T E N 1+4  v: 4
		0x454c504d41584501, 0x54454e41, 0x4,
		// E L P M A X E 1+0   V O G 1+4  v: 8
		0x454c504d41584501, 0x564f4741, 0x8,
	}, dnt.body[idx])

	dnt.InplaceInsert(netDN, 16)
	assert.Equal(t, []uint64{
		// E L P M A X E 1+0   G R O 1+4  v: 1
		0x454c504d41584501, 0x47524f41, 0x1,
		// E L P M A X E 1+0   M O C 1+4  v: 2
		0x454c504d41584501, 0x4d4f4341, 0x2,
		// E L P M A X E 1+0   T E N 1+4  v: 14
		0x454c504d41584501, 0x54454e41, 0x14,
		// E L P M A X E 1+0   V O G 1+4  v: 8
		0x454c504d41584501, 0x564f4741, 0x8,
	}, dnt.body[idx])
}

func TestTable64Get(t *testing.T) {
	dnt := NewTable64()

	dnt.InplaceInsert(comDN, 2)
	dnt.InplaceInsert(govDN, 8)
	dnt.InplaceInsert(netDN, 4)
	dnt.InplaceInsert(orgDN, 1)
	dnt.InplaceInsert(ioZone, 16)

	n, err := domain.MakeNameFromString("wjgsapgatlmody.umguqdiw.mnppqimge")
	assert.NoError(t, err)
	dnt.InplaceInsert(n, 32)

	if v, ok := dnt.Get(orgDN); assert.True(t, ok) {
		assert.EqualValues(t, 1, v)
	}

	if v, ok := dnt.Get(comDN); assert.True(t, ok) {
		assert.EqualValues(t, 2, v)
	}

	if v, ok := dnt.Get(netDN); assert.True(t, ok) {
		assert.EqualValues(t, 4, v)
	}

	if v, ok := dnt.Get(govDN); assert.True(t, ok) {
		assert.EqualValues(t, 8, v)
	}

	if v, ok := dnt.Get(ioDN); assert.True(t, ok) {
		assert.EqualValues(t, 16, v)
	}

	if v, ok := dnt.Get(n); assert.True(t, ok) {
		assert.EqualValues(t, 32, v)
	}

	_, ok := dnt.Get(rootDN)
	assert.False(t, ok)
}
