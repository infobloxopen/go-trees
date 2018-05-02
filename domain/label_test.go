package domain

import (
	"strings"
	"testing"
)

func TestLabelMarkLabels(t *testing.T) {
	s := "one.two.three.four.five"

	var offs [5]int
	n, err := markLabels(s, offs[:])
	if err != nil {
		t.Fatal(err)
	}

	if n != 5 {
		t.Fatalf("expected 5 labels but got %d:\n%#v", n, offs[:n])
	}

	eOffs := []int{0, 4, 8, 14, 19}
	for i, off := range eOffs {
		if offs[i] != off {
			t.Fatalf("expected offsets\n\t%#v\nbut got\n\t%#v", eOffs, offs[:n])
		}
	}
}

func TestLabelMarkLabelsWithEndingDot(t *testing.T) {
	s := "one.two.three.four.five."

	var offs [5]int
	n, err := markLabels(s, offs[:])
	if err != nil {
		t.Fatal(err)
	}

	if n != 5 {
		t.Fatalf("expected 5 labels but got %d:\n%#v", n, offs[:n])
	}

	eOffs := []int{0, 4, 8, 14, 19}
	for i, off := range eOffs {
		if offs[i] != off {
			t.Fatalf("expected offsets\n\t%#v\nbut got\n\t%#v", eOffs, offs[:n])
		}
	}
}

func TestLabelMarkLabelsWithEscapes(t *testing.T) {
	s := "one\\.two.th\\\\ree.four.five"

	var offs [4]int
	n, err := markLabels(s, offs[:])
	if err != nil {
		t.Fatal(err)
	}

	if n != 4 {
		t.Fatalf("expected 4 labels but got %d:\n%#v", n, offs[:n])
	}

	eOffs := []int{0, 9, 17, 22}
	for i, off := range eOffs {
		if offs[i] != off {
			t.Fatalf("expected offsets\n\t%#v\nbut got\n\t%#v", eOffs, offs[:n])
		}
	}
}

func TestLabelMarkLabelsWithInvalidEscape(t *testing.T) {
	s := "invalid\\"

	var offs [MaxLabels]int
	n, err := markLabels(s, offs[:])
	if err == nil {
		t.Fatalf("expected error but got %d offsets:\n%#v", n, offs[:n])
	}

	if err != ErrInvalidEscape {
		t.Fatalf("expected ErrInvalidEscape but got %q (%T)", err, err)
	}
}

func TestLabelMarkLabelsWithEmptyLabel(t *testing.T) {
	s := "one.two..four.five"

	var offs [MaxLabels]int
	n, err := markLabels(s, offs[:])
	if err == nil {
		t.Fatalf("expected error but got %d offsets:\n%#v", n, offs[:n])
	}

	if err != ErrEmptyLabel {
		t.Fatalf("expected ErrEmptyLabel but got %q (%T)", err, err)
	}
}

func TestLabelMarkLabelsWithTooManyLabels(t *testing.T) {
	s := "0.1.2.3.4.5"

	var offs [5]int
	n, err := markLabels(s, offs[:])
	if err == nil {
		t.Fatalf("expected error but got %d offsets:\n%#v", n, offs[:n])
	}

	if err != ErrTooManyLabels {
		t.Fatalf("expected ErrTooManyLabels but got %q (%T)", err, err)
	}
}

func TestLabelMarkLabelsWithMoreTooManyLabels(t *testing.T) {
	s := "0.1.2.3.4.5.6.7.8.9"

	var offs [5]int
	n, err := markLabels(s, offs[:])
	if err == nil {
		t.Fatalf("expected error but got %d offsets:\n%#v", n, offs[:n])
	}

	if err != ErrTooManyLabels {
		t.Fatalf("expected ErrTooManyLabels but got %q (%T)", err, err)
	}
}

func TestLabelGetLabel(t *testing.T) {
	var label [MaxLabel + 1]byte

	n, err := getLabel("label", label[:])
	if err != nil {
		t.Fatalf("expected label but got error: %s", err)
	}

	assertLabel(t, label[:n], []byte{5, 'L', 'A', 'B', 'E', 'L'})
}

func TestLabelGetLabelWithEndingDot(t *testing.T) {
	var label [MaxLabel + 1]byte

	n, err := getLabel("label.", label[:])
	if err != nil {
		t.Fatalf("expected label but got error: %s", err)
	}

	assertLabel(t, label[:n], []byte{5, 'L', 'A', 'B', 'E', 'L'})
}

func TestLabelGetLabelWithDot(t *testing.T) {
	var label [MaxLabel + 1]byte

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				t.Fatalf("expected panic with error but got %T (%#v)", r, r)
			}

			if !strings.Contains(err.Error(), "unescaped dot") {
				t.Fatalf("expected \"unescaped dot\" error but got: %s", err)
			}
		} else {
			t.Fatalf("expected panic")
		}
	}()

	n, err := getLabel("la.bel.", label[:])
	if err != nil {
		t.Fatalf("expected panic but got error: %s", err)
	}

	t.Fatalf("expected panic but got label:\n%#v", label[:n])
}

func TestLabelGetLabelWithEscapes(t *testing.T) {
	var label [MaxLabel + 1]byte

	n, err := getLabel("\\l\\a\\.\\\\\\b\\e\\l.", label[:])
	if err != nil {
		t.Fatalf("expected label but got error: %s", err)
	}

	assertLabel(t, label[:n], []byte{7, 'L', 'A', '.', '\\', 'B', 'E', 'L'})
}

func TestLabelGetLabelWithInvalidEscape(t *testing.T) {
	var label [MaxLabel + 1]byte

	n, err := getLabel("\\l\\a\\.\\\\\\b\\e\\", label[:])
	if err == nil {
		t.Fatalf("expected error but got label:\n%#v", label[:n])
	}

	if err != ErrInvalidEscape {
		t.Fatalf("expected ErrInvalidEscape but got %q (%T)", err, err)
	}
}

func TestLabelGetLabelWithLabelTooLong(t *testing.T) {
	var label [MaxLabel + 1]byte

	n, err := getLabel("0123456789012345678901234567890123456789012345678901234567890123", label[:])
	if err == nil {
		t.Fatalf("expected error but got label:\n%#v", label[:n])
	}

	if err != ErrLabelTooLong {
		t.Fatalf("expected ErrLabelTooLong but got %q (%T)", err, err)
	}
}

func TestLabelGetLabelWithEscapedLabelTooLong(t *testing.T) {
	var label [MaxLabel + 1]byte

	n, err := getLabel("looooooooooooooooooooooooooooooooooooooooooooooooooooooooooooon\\g", label[:])
	if err == nil {
		t.Fatalf("expected error but got label:\n%#v", label[:n])
	}

	if err != ErrLabelTooLong {
		t.Fatalf("expected ErrLabelTooLong but got %q (%T)", err, err)
	}
}

func assertLabel(t *testing.T, v, e []byte) {
	if len(v) != len(e) {
		t.Errorf("expected %d bytes\n\t%#v\nbut got %d\n\t%#v", len(e), e, len(v), v)
		return
	}

	for i, b := range e {
		if v[i] != b {
			t.Errorf("expected label\n\t%#v\nbut got\n\t%#v", e, v)
			return
		}
	}
}
