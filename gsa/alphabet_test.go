package gsa

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"strings"
	"testing"
)

func checkAlphabetMaps(t *testing.T, x string, letters []byte) {
	t.Helper()

	alpha := NewAlphabet(x)
	if alpha.Size() != len(letters)+1 { // +1 for sentinel
		t.Fatalf("Expected size %d, got %d", len(letters)+1, alpha.Size())
	}

	for i := 0; i < 256; i++ {
		a := byte(i)
		b := alpha._revmap[alpha._map[a]]

		if b != 0 && a != b {
			t.Fatalf("Mapping %d and then back gave us %d", a, b)
		}
	}

	for i := 0; i < 256; i++ {
		a := byte(i)
		b := alpha._map[alpha._revmap[a]]

		if b != 0 && a != b {
			t.Fatalf("Reverse mapping %d and then back gave us %d", a, b)
		}
	}
}

func Test_Alphabet_Maps(t *testing.T) {
	checkAlphabetMaps(t, "foo", []byte{'f', 'o'})
	checkAlphabetMaps(t, "foobar", []byte{'a', 'b', 'f', 'o', 'r'})

	// including the sentinel. Doesn't affect the alphabet, but
	// it will be included in the mapped string
	checkAlphabetMaps(t, "foobar\x00", []byte{'a', 'b', 'f', 'o', 'r'})
	checkAlphabetMaps(t, "foo\x00bar", []byte{'a', 'b', 'f', 'o', 'r'})
}

func checkAlphabet(
	t *testing.T, x string,
	letters []byte, expectedMapped []byte) {
	t.Helper()

	alpha := NewAlphabet(x)
	if alpha.Size() != len(letters)+1 { // +1 for sentinel
		t.Fatalf("Expected size %d, got %d", len(letters)+1, alpha.Size())
	}

	for i := 0; i < len(x); i++ {
		a := x[i]
		if !alpha.Contains(a) {
			t.Fatalf("Expected alphabet to contain %c.", a)
		}
	}

	bs, _ := alpha.MapToBytes(string(letters))
	for i, b := range bs {
		if b != byte(i+1) {
			t.Fatalf("Expected %c to map to %d", b, i+1)
		}
	}

	y, err := alpha.MapToBytesWithSentinel(x)
	if err != nil {
		t.Fatalf("Error mapping string %s", x)
	}

	if !reflect.DeepEqual(y, expectedMapped) {
		t.Fatalf("We expected the mapped string to be %v but it is %v", expectedMapped, y)
	}

	xx := strings.ReplaceAll(x, string(Sentinel), string(SentinelSymbol))

	z := alpha.RevmapBytesStripSentinel(y)
	if !reflect.DeepEqual(xx, z) {
		t.Fatalf(`Mapping back we expected "%s" but got "%s"`, x, z)
	}

	zz := alpha.RevmapBytes(y)
	if !reflect.DeepEqual(zz[:len(z)], z) {
		t.Fatalf(`Strings "%s" and "%s" should be equal`, z, zz)
	}

	if zz != z+string(SentinelSymbol) {
		t.Fatalf(`The last character in "%s" should be sentinel`, zz)
	}

	checkAlphabetInt(t, x, alpha, y, expectedMapped)
}

func checkAlphabetInt(t *testing.T, x string, alpha *Alphabet, y, expectedMapped []byte) {
	t.Helper()

	expectedInts := make([]int32, len(expectedMapped))
	for i, b := range expectedMapped {
		expectedInts[i] = int32(b)
	}

	yy, err2 := alpha.MapToIntsWithSentinel(x)
	if err2 != nil {
		t.Fatalf("Error mapping string %s", x)
	}

	if !reflect.DeepEqual(yy, expectedInts) {
		t.Fatalf("We expected the mapped string to be %v but it is %v", expectedMapped, y)
	}
}

func Test_Alphabet(t *testing.T) {
	checkAlphabet(t, "foo", []byte{'f', 'o'}, []byte{1, 2, 2, 0})

	checkAlphabet(t, "foobar", []byte{'a', 'b', 'f', 'o', 'r'}, []byte{3, 4, 4, 2, 1, 5, 0})

	// including the sentinel. Doesn't affect the alphabet, but
	// it will be included in the mapped string
	checkAlphabet(t, "foobar\x00", []byte{'a', 'b', 'f', 'o', 'r'}, []byte{3, 4, 4, 2, 1, 5, 0, 0})
	checkAlphabet(t, "foo\x00bar", []byte{'a', 'b', 'f', 'o', 'r'}, []byte{3, 4, 4, 0, 2, 1, 5, 0})

	// See what happens if we try to map a string with the wrong alphabet.
	alpha := NewAlphabet("foo")
	if _, err := alpha.MapToBytes("foobar"); err == nil {
		t.Fatalf("Expected an error when mapping to an alphabet that doesn't match")
	}

	if _, err := alpha.MapToBytesWithSentinel("foobar"); err == nil {
		t.Fatalf("Expected an error when mapping to an alphabet that doesn't match")
	}

	if _, err := alpha.MapToInts("foobar"); err == nil {
		t.Fatalf("Expected an error when mapping to an alphabet that doesn't match")
	}

	if _, err := alpha.MapToIntsWithSentinel("foobar"); err == nil {
		t.Fatalf("Expected an error when mapping to an alphabet that doesn't match")
	}
}

func TestAlphabetEncoding(t *testing.T) {
	x := "foobar"
	alpha := NewAlphabet(x)
	beta := NewAlphabet("blah")

	if reflect.DeepEqual(alpha, beta) {
		t.Fatalf("These two alphabets should *not* be equal")
	}

	var (
		buf bytes.Buffer
		enc = gob.NewEncoder(&buf)
		dec = gob.NewDecoder(&buf)
	)

	if err := enc.Encode(&alpha); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if err := dec.Decode(&beta); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if !reflect.DeepEqual(alpha, beta) {
		t.Fatalf("These two alphabets should be equal now")
	}
}
