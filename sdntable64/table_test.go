package sdntable64

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infobloxopen/go-trees/udomain"
)

func TestTable64NewTable64(t *testing.T) {
	dnt := NewTable64()
	assert.NotNil(t, dnt)
	assert.True(t, dnt.ready)
}

func TestTable64InplaceInsert(t *testing.T) {
	dnt := NewTable64()
	dnt.InplaceInsert(makeDomainNameFromString("."), 1)
	assert.Equal(t, uint64(1), dnt.root)
	assert.True(t, dnt.ready)

	n := makeDomainNameFromString("example.com")
	i := len(n.GetComparable()) - 1

	dnt.InplaceInsert(n, 2)
	dnt.InplaceInsert(makeDomainNameFromString("example.gov"), 16)
	dnt.InplaceInsert(makeDomainNameFromString("example.net"), 4)
	dnt.InplaceInsert(makeDomainNameFromString("example.org"), 1)
	dnt.InplaceInsert(makeDomainNameFromString("example.gov"), 8)

	assert.True(t, dnt.ready)
	assert.Equal(t, []int64{
		// E L P M A X E 1+0   G R O 1+4
		0x454c504d41584501, 0x47524f41,
		// E L P M A X E 1+0   M O C 1+4
		0x454c504d41584501, 0x4d4f4341,
		// E L P M A X E 1+0   T E N 1+4
		0x454c504d41584501, 0x54454e41,
		// E L P M A X E 1+0   V O G 1+4
		0x454c504d41584501, 0x564f4741,
	}, dnt.body[i].data.k)
	assert.Equal(t, []uint64{1, 2, 4, 8}, dnt.body[i].data.v)
}

func TestTable64Append(t *testing.T) {
	dnt := NewTable64()
	dnt, _ = dnt.Append(makeDomainNameFromString("."), 1)
	assert.Equal(t, uint64(1), dnt.root)
	assert.True(t, dnt.ready)

	n := makeDomainNameFromString("example.com")
	i := len(n.GetComparable()) - 1

	dnt, _ = dnt.Append(n, 2)
	assert.False(t, dnt.ready)
	assert.Equal(t, []int64{
		// E L P M A X E 1+0   M O C 1+4
		0x454c504d41584501, 0x4d4f4341,
	}, dnt.body[i].data.k)
	assert.Equal(t, []uint64{2}, dnt.body[i].data.v)

	n = makeDomainNameFromString("example.net")
	j := len(n.GetComparable()) - 1
	assert.Equal(t, i, j)

	dnt, _ = dnt.Append(n, 4)
	assert.False(t, dnt.ready)
	assert.Equal(t, []int64{
		// E L P M A X E 1+0   M O C 1+4
		0x454c504d41584501, 0x4d4f4341,
		// E L P M A X E 1+0   T E N 1+4
		0x454c504d41584501, 0x54454e41,
	}, dnt.body[i].data.k)
	assert.Equal(t, []uint64{2, 4}, dnt.body[i].data.v)
}

func TestTable64Normalize(t *testing.T) {
	dnt := NewTable64()

	n := makeDomainNameFromString("example.com")
	i := len(n.GetComparable()) - 1

	dnt, _ = dnt.Append(n, 2)
	dnt, _ = dnt.Append(makeDomainNameFromString("example.gov"), 16)
	dnt, _ = dnt.Append(makeDomainNameFromString("example.net"), 4)
	dnt, _ = dnt.Append(makeDomainNameFromString("example.org"), 1)
	dnt, _ = dnt.Append(makeDomainNameFromString("example.gov"), 8)
	assert.False(t, dnt.ready)

	dnt = dnt.Normalize()
	assert.True(t, dnt.ready)
	assert.Equal(t, []int64{
		// E L P M A X E 1+0   G R O 1+4
		0x454c504d41584501, 0x47524f41,
		// E L P M A X E 1+0   M O C 1+4
		0x454c504d41584501, 0x4d4f4341,
		// E L P M A X E 1+0   T E N 1+4
		0x454c504d41584501, 0x54454e41,
		// E L P M A X E 1+0   V O G 1+4
		0x454c504d41584501, 0x564f4741,
	}, dnt.body[i].data.k)
	assert.Equal(t, []uint64{1, 2, 4, 8}, dnt.body[i].data.v)
}

func TestTable64Get(t *testing.T) {
	dnt := NewTable64()
	dnt.InplaceInsert(makeDomainNameFromString("."), 1)
	assert.Equal(t, uint64(1), dnt.root)

	assert.Equal(t, uint64(1), dnt.Get(makeDomainNameFromString(".")))

	dnt = NewTable64()
	dnt, _ = dnt.Append(makeDomainNameFromString("example.com"), 2)
	dnt, _ = dnt.Append(makeDomainNameFromString("example.gov"), 16)
	dnt, _ = dnt.Append(makeDomainNameFromString("example.net"), 4)
	dnt, _ = dnt.Append(makeDomainNameFromString("example.org"), 1)
	dnt, _ = dnt.Append(makeDomainNameFromString("example.gov"), 8)
	dnt = dnt.Normalize()

	assert.Equal(t, uint64(0), dnt.Get(makeDomainNameFromString("example.edu")))
	assert.Equal(t, uint64(1), dnt.Get(makeDomainNameFromString("example.org")))
	assert.Equal(t, uint64(2), dnt.Get(makeDomainNameFromString("example.com")))
	assert.Equal(t, uint64(4), dnt.Get(makeDomainNameFromString("example.net")))
	assert.Equal(t, uint64(8), dnt.Get(makeDomainNameFromString("example.gov")))
	assert.Equal(t, uint64(8), dnt.Get(makeDomainNameFromString("www.example.gov")))
}

func makeDomainNameFromString(s string) domain.Name {
	n, err := domain.MakeNameFromString(s)
	if err != nil {
		panic(fmt.Errorf("can't make domain name from string %q: %s", s, err))
	}

	return n
}
