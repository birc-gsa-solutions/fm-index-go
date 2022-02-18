package gsa

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"

	"birc.au.dk/gsa/test"
)

func Test_Ctab(t *testing.T) {
	x := "aab"
	alpha := NewAlphabet(x)
	xb, _ := alpha.MapToBytesWithSentinel(x)
	ctab := NewCTab(xb, alpha.Size())

	if len(ctab.CumSum) != alpha.Size() {
		t.Fatal("The ctable's cumsum has the wrong length")
	}

	if !reflect.DeepEqual(ctab.CumSum, []int{0, 1, 3}) {
		t.Fatal("We have the wrong cumsum")
	}
}

func Test_Otab(t *testing.T) {
	x := "aab"
	alpha := NewAlphabet(x)
	sa, _ := SaisWithAlphabet(x, alpha)

	xb, _ := alpha.MapToBytesWithSentinel(x)
	bwt := Bwt(xb, sa)
	otab := NewOTab(bwt, alpha.Size())

	expectedBwt := []byte{2, 0, 1, 1}
	if !reflect.DeepEqual(bwt, expectedBwt) {
		t.Fatalf("Expected bwt %v, got %v", expectedBwt, bwt)
	}

	expectedA := []int{0, 0, 0, 1, 2}
	expectedB := []int{0, 1, 1, 1, 1}

	var (
		a byte = 1
		b byte = 2
	)

	for i := range expectedA {
		if otab.Rank(a, i) != expectedA[i] {
			t.Errorf("Unexpected value at Rank(%b,%d) = %d\n", a, i, otab.Rank(a, i))
		}
	}

	for i := range expectedB {
		if otab.Rank(b, i) != expectedB[i] {
			t.Errorf("Unexpected value at Rank(%b,%d) = %d\n", b, i, otab.Rank(b, i))
		}
	}
}

func Test_MississippiBWT(t *testing.T) {
	xs := "mississippi"
	ps := "is"

	preproc := FMIndexExactPreprocess(xs)
	preproc(ps, func(i int32) {
		test.CheckOccurrenceAt(t, xs, ps, int(i))
	})
}

func TestOTabEncoding(t *testing.T) {
	x := "aab"
	alpha := NewAlphabet(x)
	sa, _ := SaisWithAlphabet(x, alpha)
	xb, _ := alpha.MapToBytesWithSentinel(x)
	bwt := Bwt(xb, sa)

	otab1 := NewOTab(bwt, alpha.Size())
	otab2 := &OTab{}

	if reflect.DeepEqual(otab1, otab2) {
		t.Fatalf("The two otables should not be equal yet")
	}

	var (
		buf bytes.Buffer
		enc = gob.NewEncoder(&buf)
		dec = gob.NewDecoder(&buf)
	)

	if err := enc.Encode(&otab1); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if err := dec.Decode(&otab2); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if !reflect.DeepEqual(otab1, otab2) {
		t.Errorf("These two otables should be equal now")
	}
}
