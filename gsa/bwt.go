package gsa

import (
	"bytes"
	"encoding/gob"
)

// Bwt gives you the Burrows-Wheeler transform of a string,
// computed using the suffix array for the string. The
// string should have a sentinel
func Bwt(x []byte, sa []int32) []byte {
	bwt := make([]byte, len(x))

	for i := 0; i < len(sa); i++ {
		j := sa[i]
		if j == 0 {
			bwt[i] = 0
		} else {
			bwt[i] = x[j-1]
		}
	}

	return bwt
}

// CTab Structur holding the C-table for BWT search.
// This is a map from letters in the alphabet to the
// cumulative sum of how often we see letters in the
// BWT
type CTab struct {
	CumSum []int
}

// Rank How many times does the BWT hold a letter smaller
// than a? Undefined behaviour if a isn't in the table.
func (ctab *CTab) Rank(a byte) int {
	return ctab.CumSum[a]
}

// NewCTab builds the c-table from a string.
func NewCTab(bwt []byte, asize int) *CTab {
	// First, count how often we see each character
	counts := make([]int, asize)
	for _, b := range bwt {
		counts[b]++
	}
	// Then get the accumulative sum
	var n int
	for i, count := range counts {
		counts[i] = n
		n += count
	}

	return &CTab{counts}
}

// OTab Holds the o-table (rank table) from a BWT string
type OTab struct {
	nrow, ncol int
	table      []int
}

func (otab *OTab) offset(a byte, i int) int {
	// -1 to a because we don't store the sentinel
	// and -1 to i because we don't store the first
	// row (which is always zero)
	return otab.ncol*(int(a)-1) + (i - 1)
}

func (otab *OTab) get(a byte, i int) int {
	return otab.table[otab.offset(a, i)]
}

func (otab *OTab) set(a byte, i, val int) {
	otab.table[otab.offset(a, i)] = val
}

// Rank How many times do we see letter a before index i
// in the BWT string?
func (otab *OTab) Rank(a byte, i int) int {
	// We don't explicitly store the first column,
	// since it is always empty anyway.
	if i == 0 {
		return 0
	}

	return otab.get(a, i)
}

// NewOTab builds the o-table from a string. It uses
// the suffix array to get the BWT and a c-table
// to handle the alphabet.
func NewOTab(bwt []byte, asize int) *OTab {
	// We index for all characters except $, so
	// nrow is the alphabet size minus one.
	// We index all indices [0,len(sa)], but we emulate
	// row 0, since it is always zero, so we only need
	// len(sa) columns.
	nrow, ncol := asize-1, len(bwt)
	table := make([]int, nrow*ncol)
	otab := OTab{nrow, ncol, table}

	// The character at the beginning of bwt gets a count
	// of one at row one.
	otab.set(bwt[0], 1, 1)

	// The remaining entries either copies or increment from
	// the previous column. We count a from 1 to alpha size
	// to skip the sentinel, then -1 for the index
	for a := 1; a < asize; a++ {
		ba := byte(a) // get the right type for accessing otab
		for i := 2; i <= len(bwt); i++ {
			val := otab.get(ba, i-1)
			if bwt[i-1] == ba {
				val++
			}

			otab.set(ba, i, val)
		}
	}

	return &otab
}

// FMIndexTables contains the preprocessed tables used for FM-index
// searching
type FMIndexTables struct {
	Alpha *Alphabet
	Sa    []int32
	Ctab  *CTab
	Otab  *OTab
}

// BuildFMIndexExactTables builds the preprocessing tables for exact FM-index
// searching.
func BuildFMIndexExactTables(x string) *FMIndexTables {
	xb, alpha := MapStringWithSentinel(x)
	sa, _ := SaisWithAlphabet(x, alpha)
	bwt := Bwt(xb, sa)
	ctab := NewCTab(bwt, alpha.Size())
	otab := NewOTab(bwt, alpha.Size())

	return &FMIndexTables{
		Alpha: alpha,
		Sa:    sa,
		Ctab:  ctab,
		Otab:  otab,
	}
}

// FMIndexExactFromTables returns a search function based
// on the preprocessed tables
func FMIndexExactFromTables(tbls *FMIndexTables) func(p string, cb func(i int32)) {
	return func(p string, cb func(i int32)) {
		pb, err := tbls.Alpha.MapToBytes(p)
		if err != nil {
			return // p doesn't fit the alphabet, so we can't match
		}

		left, right := 0, len(tbls.Sa)

		for i := len(pb) - 1; i >= 0; i-- {
			a := pb[i]
			left = tbls.Ctab.Rank(a) + tbls.Otab.Rank(a, left)
			right = tbls.Ctab.Rank(a) + tbls.Otab.Rank(a, right)

			if left >= right {
				return // no match
			}
		}

		for i := left; i < right; i++ {
			cb(tbls.Sa[i])
		}
	}
}

// FMIndexExactPreprocess preprocesses the string x and returns a function
// that you can use to efficiently search in x.
func FMIndexExactPreprocess(x string) func(p string, cb func(i int32)) {
	return FMIndexExactFromTables(BuildFMIndexExactTables(x))
}

// GobEncode implements the encoder interface for serialising to a stream of bytes
func (otab OTab) GobEncode() (res []byte, err error) {
	defer catchError(&err)

	var (
		buf bytes.Buffer
		enc = gob.NewEncoder(&buf)
	)

	checkError(enc.Encode(otab.nrow))
	checkError(enc.Encode(otab.ncol))
	checkError(enc.Encode(otab.table))

	return buf.Bytes(), nil
}

// GobDecode implements the decoder interface for serialising to a stream of bytes
func (otab *OTab) GobDecode(b []byte) (err error) {
	defer catchError(&err)

	r := bytes.NewReader(b)
	dec := gob.NewDecoder(r)

	checkError(dec.Decode(&otab.nrow))
	checkError(dec.Decode(&otab.ncol))

	return dec.Decode(&otab.table)
}
