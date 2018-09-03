package sdntable64

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainsWriteAll(t *testing.T) {
	noLogs := normalizeLoggers{
		before: nil,
		sort:   nil,
		after:  nil,
	}

	ds := makeDomains()
	defer assert.NoError(t, ds.close())
	for i := 0; i < 20; i++ {
		d := makeDomainNameFromString(fmt.Sprintf("%02d.example.com", i))
		ds = ds.append(d.GetComparable(), 1<<uint(i))
	}

	ds = ds.normalize(noLogs)
	assert.Equal(t, []int64{
		// 0 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x303031, 0x454c504d41584501, 0x4d4f4341,
		// 0 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x303131, 0x454c504d41584501, 0x4d4f4341,
		// 1 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x313031, 0x454c504d41584501, 0x4d4f4341,
		// 1 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x313131, 0x454c504d41584501, 0x4d4f4341,
		// 2 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x323031, 0x454c504d41584501, 0x4d4f4341,
		// 2 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x323131, 0x454c504d41584501, 0x4d4f4341,
		// 3 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x333031, 0x454c504d41584501, 0x4d4f4341,
		// 3 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x333131, 0x454c504d41584501, 0x4d4f4341,
		// 4 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x343031, 0x454c504d41584501, 0x4d4f4341,
		// 4 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x343131, 0x454c504d41584501, 0x4d4f4341,
		// 5 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x353031, 0x454c504d41584501, 0x4d4f4341,
		// 5 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x353131, 0x454c504d41584501, 0x4d4f4341,
		// 6 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x363031, 0x454c504d41584501, 0x4d4f4341,
		// 6 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x363131, 0x454c504d41584501, 0x4d4f4341,
		// 7 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x373031, 0x454c504d41584501, 0x4d4f4341,
		// 7 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x373131, 0x454c504d41584501, 0x4d4f4341,
		// 8 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x383031, 0x454c504d41584501, 0x4d4f4341,
		// 8 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x383131, 0x454c504d41584501, 0x4d4f4341,
		// 9 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x393031, 0x454c504d41584501, 0x4d4f4341,
		// 9 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x393131, 0x454c504d41584501, 0x4d4f4341,
	}, ds.data.k)
	assert.Equal(t, []uint64{
		0x00001, 0x00400, 0x00002, 0x00800, 0x00004, 0x01000, 0x00008, 0x02000,
		0x00010, 0x04000, 0x00020, 0x08000, 0x00040, 0x10000, 0x00080, 0x20000,
		0x00100, 0x40000, 0x00200, 0x80000,
	}, ds.data.v)

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		assert.FailNow(t, "failed to create directory: %s", err)
	}
	ds.dir = &dir
	defer os.RemoveAll(dir)

	blks2 := blocks[2]
	blocks[2] = 5
	defer func() { blocks[2] = blks2 }()

	ds, err = ds.writeAll(nil)
	assert.NoError(t, err)
	assert.Equal(t, 4, ds.blks)
	assert.Equal(t, 0, ds.rem)

	assert.Equal(t, []uint32{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
	}, ds.data.i)
}

func TestDomainsDrop90(t *testing.T) {
	noLogs := normalizeLoggers{
		before: nil,
		sort:   nil,
		after:  nil,
	}

	ds := makeDomains()
	defer assert.NoError(t, ds.close())
	for i := 0; i < 20; i++ {
		d := makeDomainNameFromString(fmt.Sprintf("%02d.example.com", i))
		ds = ds.append(d.GetComparable(), 1<<uint(i))
	}

	ds = ds.normalize(noLogs)

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		assert.FailNow(t, "failed to create directory: %s", err)
	}
	ds.dir = &dir
	defer os.RemoveAll(dir)

	blks2 := blocks[2]
	blocks[2] = 5
	defer func() { blocks[2] = blks2 }()

	ds, err = ds.writeAll(nil)
	assert.NoError(t, err)

	ds = ds.drop90()
	assert.Equal(t, []int64{
		// 0 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x303031, 0x454c504d41584501, 0x4d4f4341,
		// 5 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x353031, 0x454c504d41584501, 0x4d4f4341,
	}, ds.data.k)
	assert.Equal(t, []uint64{0x00001, 0x00020}, ds.data.v)
	assert.Equal(t, []uint32{0, 10}, ds.data.i)
}

func TestDomainsMergeAndReadFromDisk(t *testing.T) {
	noLogs := normalizeLoggers{
		before: nil,
		sort:   nil,
		after:  nil,
	}

	ds := makeDomains()
	defer assert.NoError(t, ds.close())
	for i := 0; i < 20; i++ {
		d := makeDomainNameFromString(fmt.Sprintf("%02d.example.com", i))
		ds = ds.append(d.GetComparable(), 1<<uint(i))
	}

	ds = ds.normalize(noLogs)

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		assert.FailNow(t, "failed to create directory: %s", err)
	}
	ds.dir = &dir
	defer os.RemoveAll(dir)

	blks2 := blocks[2]
	blocks[2] = 5
	defer func() { blocks[2] = blks2 }()

	ds, err = ds.writeAll(nil)
	assert.NoError(t, err)

	ds = ds.drop90()
	for i := 20; i < 30; i++ {
		d := makeDomainNameFromString(fmt.Sprintf("%02d.example.com", i))
		ds = ds.append(d.GetComparable(), 1<<uint(i))
	}
	ds = ds.normalize(noLogs)
	assert.Equal(t, []int64{
		// 0 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x303031, 0x454c504d41584501, 0x4d4f4341,
		// 0 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x303231, 0x454c504d41584501, 0x4d4f4341,
		// 1 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x313231, 0x454c504d41584501, 0x4d4f4341,
		// 2 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x323231, 0x454c504d41584501, 0x4d4f4341,
		// 3 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x333231, 0x454c504d41584501, 0x4d4f4341,
		// 4 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x343231, 0x454c504d41584501, 0x4d4f4341,
		// 5 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x353031, 0x454c504d41584501, 0x4d4f4341,
		// 5 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x353231, 0x454c504d41584501, 0x4d4f4341,
		// 6 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x363231, 0x454c504d41584501, 0x4d4f4341,
		// 7 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x373231, 0x454c504d41584501, 0x4d4f4341,
		// 8 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x383231, 0x454c504d41584501, 0x4d4f4341,
		// 9 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x393231, 0x454c504d41584501, 0x4d4f4341,
	}, ds.data.k)
	assert.Equal(t, []uint64{
		0x0000001, 0x00100000, 0x00200000, 0x00400000, 0x00800000, 0x01000000, 0x00000020, 0x02000000,
		0x4000000, 0x08000000, 0x10000000, 0x20000000,
	}, ds.data.v)
	assert.Equal(t, []uint32{
		0, math.MaxUint32, math.MaxUint32, math.MaxUint32, math.MaxUint32, math.MaxUint32,
		10, math.MaxUint32, math.MaxUint32, math.MaxUint32, math.MaxUint32, math.MaxUint32,
	}, ds.data.i)

	ds, g, err := ds.merge(nil)
	assert.NoError(t, err)
	assert.Equal(t, []int64{
		// 0 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x303031, 0x454c504d41584501, 0x4d4f4341,
		// 0 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x303231, 0x454c504d41584501, 0x4d4f4341,
		// 1 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x313231, 0x454c504d41584501, 0x4d4f4341,
		// 2 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x323231, 0x454c504d41584501, 0x4d4f4341,
		// 3 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x333231, 0x454c504d41584501, 0x4d4f4341,
		// 4 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x343231, 0x454c504d41584501, 0x4d4f4341,
		// 5 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x353031, 0x454c504d41584501, 0x4d4f4341,
		// 5 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x353231, 0x454c504d41584501, 0x4d4f4341,
		// 6 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x363231, 0x454c504d41584501, 0x4d4f4341,
		// 7 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x373231, 0x454c504d41584501, 0x4d4f4341,
		// 8 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x383231, 0x454c504d41584501, 0x4d4f4341,
		// 9 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x393231, 0x454c504d41584501, 0x4d4f4341,
	}, ds.data.k)
	assert.Equal(t, []uint64{
		0x0000001, 0x00100000, 0x00200000, 0x00400000, 0x00800000, 0x01000000, 0x00000020, 0x02000000,
		0x4000000, 0x08000000, 0x10000000, 0x20000000,
	}, ds.data.v)
	if g != nil {
		assert.NoError(t, g.Stop())
	}

	// 00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29
	// 00 10 20 01 11 21 02 12 22 03 13 23 04 14 24 05 15 25 06 16 26 07 17 27 08 18 28 09 19 29
	assert.Equal(t, []uint32{
		0, 2, 5, 8, 11, 14,
		15, 17, 20, 23, 26, 29,
	}, ds.data.i)

	f, ok := ds.getter.f.(*os.File)
	assert.True(t, ok)
	tmp, err := ds.getter.read(f, 0, math.MaxUint32,
		run{
			k: []int64{},
			v: []uint64{},
		},
		run{
			k: make([]int64, g.m*g.blk),
			v: make([]uint64, g.blk),
		},
		nil,
	)
	assert.NoError(t, err)
	assert.Equal(t, []int64{
		// 0 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x303031, 0x454c504d41584501, 0x4d4f4341,
		// 0 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x303131, 0x454c504d41584501, 0x4d4f4341,
		// 0 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x303231, 0x454c504d41584501, 0x4d4f4341,
		// 1 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x313031, 0x454c504d41584501, 0x4d4f4341,
		// 1 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x313131, 0x454c504d41584501, 0x4d4f4341,
		// 1 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x313231, 0x454c504d41584501, 0x4d4f4341,
		// 2 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x323031, 0x454c504d41584501, 0x4d4f4341,
		// 2 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x323131, 0x454c504d41584501, 0x4d4f4341,
		// 2 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x323231, 0x454c504d41584501, 0x4d4f4341,
		// 3 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x333031, 0x454c504d41584501, 0x4d4f4341,
		// 3 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x333131, 0x454c504d41584501, 0x4d4f4341,
		// 3 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x333231, 0x454c504d41584501, 0x4d4f4341,
		// 4 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x343031, 0x454c504d41584501, 0x4d4f4341,
		// 4 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x343131, 0x454c504d41584501, 0x4d4f4341,
		// 4 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x343231, 0x454c504d41584501, 0x4d4f4341,
		// 5 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x353031, 0x454c504d41584501, 0x4d4f4341,
		// 5 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x353131, 0x454c504d41584501, 0x4d4f4341,
		// 5 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x353231, 0x454c504d41584501, 0x4d4f4341,
		// 6 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x363031, 0x454c504d41584501, 0x4d4f4341,
		// 6 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x363131, 0x454c504d41584501, 0x4d4f4341,
		// 6 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x363231, 0x454c504d41584501, 0x4d4f4341,
		// 7 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x373031, 0x454c504d41584501, 0x4d4f4341,
		// 7 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x373131, 0x454c504d41584501, 0x4d4f4341,
		// 7 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x373231, 0x454c504d41584501, 0x4d4f4341,
		// 8 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x383031, 0x454c504d41584501, 0x4d4f4341,
		// 8 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x383131, 0x454c504d41584501, 0x4d4f4341,
		// 8 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x383231, 0x454c504d41584501, 0x4d4f4341,
		// 9 0 1+3   E L P M A X E 1+0   M O C 1+4
		0x393031, 0x454c504d41584501, 0x4d4f4341,
		// 9 1 1+3   E L P M A X E 1+0   M O C 1+4
		0x393131, 0x454c504d41584501, 0x4d4f4341,
		// 9 2 1+3   E L P M A X E 1+0   M O C 1+4
		0x393231, 0x454c504d41584501, 0x4d4f4341,
	}, tmp.k)
}

func TestDomainsGet(t *testing.T) {
	noLogs := normalizeLoggers{
		before: nil,
		sort:   nil,
		after:  nil,
	}

	ds := makeDomains()
	defer assert.NoError(t, ds.close())
	for i := 0; i < 20; i++ {
		d := makeDomainNameFromString(fmt.Sprintf("%02d.example.com", i))
		ds = ds.append(d.GetComparable(), 1<<uint(i))
	}

	ds = ds.normalize(noLogs)

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		assert.FailNow(t, "failed to create directory: %s", err)
	}
	ds.dir = &dir
	defer os.RemoveAll(dir)

	blks2 := blocks[2]
	blocks[2] = 5
	defer func() { blocks[2] = blks2 }()

	ds, err = ds.writeAll(nil)
	assert.NoError(t, err)

	ds = ds.drop90()
	for i := 20; i < 30; i++ {
		d := makeDomainNameFromString(fmt.Sprintf("%02d.example.com", i))
		ds = ds.append(d.GetComparable(), 1<<uint(i))
	}
	ds = ds.normalize(noLogs)

	assert.Equal(t, uint64(0), ds.get(makeDomainNameFromString("--.example.com").GetComparable(), nil))
	assert.Equal(t, uint64(1<<15), ds.get(makeDomainNameFromString("15.example.com").GetComparable(), nil))

	for i := 30; i < 32; i++ {
		d := makeDomainNameFromString(fmt.Sprintf("%02d.example.com", i))
		ds = ds.append(d.GetComparable(), 1<<uint(i))
	}
	ds = ds.normalize(noLogs)

	ds, g, err := ds.merge(nil)
	assert.NoError(t, err)
	if g != nil {
		assert.NoError(t, g.Stop())
	}

	assert.Equal(t, uint64(1<<19), ds.get(makeDomainNameFromString("19.example.com").GetComparable(), nil))
}
