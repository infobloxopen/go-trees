package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameMakeNameFromString(t *testing.T) {
	s := "looooooooooooooong.example.com"
	n, err := MakeNameFromString(s)
	assert.NoError(t, err)
	assert.Equal(t, s, n.h, "human-readable name should be the same as input string")
	assert.Equal(t, []int64{
		// O O O O O O L (3 bytes in incomplete dword and 3 dwords for the label)
		0x4f4f4f4f4f4f4c33,
		// O O O O O O O O
		0x4f4f4f4f4f4f4f4f,
		//           G N O
		0x0000000000474e4f,
		// E L P M A X E (0 bytes in incomplete dword and 1 dword for the label)
		0x454c504d41584501,
		//         M O C (4 bytes in incomplete dword and 1 dword for the label)
		0x000000004d4f4341,
	}, n.c)
}

func TestNameMakeNameFromStringFQDN(t *testing.T) {
	s := "looooooooooooooong.example.com."
	n, err := MakeNameFromString(s)
	assert.NoError(t, err)
	assert.Equal(t, s, n.h, "human-readable name should be the same as input string")
	assert.Equal(t, []int64{
		// O O O O O O L (3 bytes in incomplete dword and 3 dwords for the label)
		0x4f4f4f4f4f4f4c33,
		// O O O O O O O O
		0x4f4f4f4f4f4f4f4f,
		//           G N O
		0x0000000000474e4f,
		// E L P M A X E (0 bytes in incomplete dword and 1 dword for the label)
		0x454c504d41584501,
		//         M O C (4 bytes in incomplete dword and 1 dword for the label)
		0x000000004d4f4341,
	}, n.c)
}

func TestNameMakeNameFromStringEmpty(t *testing.T) {
	s := ""
	n, err := MakeNameFromString(s)
	assert.NoError(t, err)
	assert.Equal(t, s, n.h, "human-readable name should be the same as input string")
	assert.Equal(t, []int64(nil), n.c)
}

func TestNameMakeNameFromStringDot(t *testing.T) {
	s := "."
	n, err := MakeNameFromString(s)
	assert.NoError(t, err)
	assert.Equal(t, s, n.h, "human-readable name should be the same as input string")
	assert.Equal(t, []int64(nil), n.c)
}

func TestNameMakeNameFromStringWithEscapedDot(t *testing.T) {
	s := "www\\.example.com"
	n, err := MakeNameFromString(s)
	assert.NoError(t, err)
	assert.Equal(t, s, n.h, "human-readable name should be the same as input string")
	assert.Equal(t, []int64{
		// A X E . W W W (4 bytes in incomplete dword and 2 dwords for the label)
		0x4158452e57575742,
		//         E L P M
		0x00000000454c504d,
		//         M O C (4 bytes in incomplete dword and 1 dword for the label)
		0x000000004d4f4341,
	}, n.c)
}

func TestNameMakeNameFromStringWithEscapedChar(t *testing.T) {
	s := "www.e\\xample.com"
	n, err := MakeNameFromString(s)
	assert.NoError(t, err)
	assert.Equal(t, s, n.h, "human-readable name should be the same as input string")
	assert.Equal(t, []int64{
		//         W W W (4 bytes in incomplete dword and 1 dword for the label)
		0x0000000057575741,
		// E L P M A X E (0 bytes in incomplete dword and 1 dword for the label)
		0x454c504d41584501,
		//         M O C (4 bytes in incomplete dword and 1 dword for the label)
		0x000000004d4f4341,
	}, n.c)
}

func TestNameMakeNameFromStringWithEscapedCode(t *testing.T) {
	s := "www.e\\120ample.com"
	n, err := MakeNameFromString(s)
	assert.NoError(t, err)
	assert.Equal(t, s, n.h, "human-readable name should be the same as input string")
	assert.Equal(t, []int64{
		//         W W W (4 bytes in incomplete dword and 1 dword for the label)
		0x0000000057575741,
		// E L P M A X E (0 bytes in incomplete dword and 1 dword for the label)
		0x454c504d41584501,
		//         M O C (4 bytes in incomplete dword and 1 dword for the label)
		0x000000004d4f4341,
	}, n.c)
}

func TestNameMakeNameFromStringWithNameTooLong(t *testing.T) {
	s := "toooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo." +
		"loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong." +
		"doooooooooooooooooooooooooooooooooooooooooooooooooooooooooomain." +
		"naaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaame"

	n, err := MakeNameFromString(s)
	assert.Equal(t, ErrNameTooLong, err, "name %q %08x", n, n.c)
}

func TestNameMakeNameFromStringWithTooLongLabel(t *testing.T) {
	s := "www1.looooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong.com"
	n, err := MakeNameFromString(s)
	assert.Equal(t, ErrLabelTooLong, err, "name %q %08x", n, n.c)

	s = "www2.looooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong.com"
	n, err = MakeNameFromString(s)
	assert.Equal(t, ErrLabelTooLong, err, "name %q %08x", n, n.c)

	s = "www3.looooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong"
	n, err = MakeNameFromString(s)
	assert.Equal(t, ErrLabelTooLong, err, "name %q %08x", n, n.c)
}

func TestNameMakeNameFromStringWithEmptyLabel(t *testing.T) {
	s := "empty..label"
	n, err := MakeNameFromString(s)
	assert.Equal(t, ErrEmptyLabel, err, "name %q %08x", n, n.c)
}

func TestNameString(t *testing.T) {
	s := "example.com"
	n, err := MakeNameFromString(s)
	assert.NoError(t, err)
	assert.Equal(t, s, n.String(), "human-readable name should be the same as input string")
}

func TestNameLess(t *testing.T) {
	n1, err := MakeNameFromString("example.com")
	assert.NoError(t, err)

	n2, err := MakeNameFromString("example.com")
	assert.NoError(t, err)

	assert.False(t, n1.Less(n2))
	assert.False(t, n2.Less(n1))

	a, err := MakeNameFromString("short.example.com")
	assert.NoError(t, err)

	b, err := MakeNameFromString("loooooong.example.com")
	assert.NoError(t, err)

	assert.True(t, a.Less(b))
	assert.False(t, b.Less(a))

	a, err = MakeNameFromString("example.com")
	assert.NoError(t, err)

	b, err = MakeNameFromString("example.net")
	assert.NoError(t, err)

	assert.True(t, a.Less(b))
	assert.False(t, b.Less(a))
}