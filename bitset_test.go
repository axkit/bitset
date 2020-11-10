package bitset

import (
	"testing"
)

func TestBitSet_SetSingle(t *testing.T) {

	var tc = []struct {
		bitpos     uint // target bit pos
		eslicelen  uint // expected slice len
		ecnt       uint // expected max cnt
		byteIndex  int  // slice index where bitpos belogs to
		eByteValue byte // int8 value of slice item
	}{
		{
			bitpos:     0,
			eslicelen:  1,
			ecnt:       1,
			byteIndex:  0,
			eByteValue: 1,
		},
		{
			bitpos:     1,
			eslicelen:  1,
			ecnt:       2,
			byteIndex:  0,
			eByteValue: 2,
		},
		{
			bitpos:     9,
			eslicelen:  2,
			ecnt:       10,
			byteIndex:  1,
			eByteValue: 2,
		},
	}

	for i := range tc {
		a := BitSet{}
		a.Set(tc[i].bitpos)
		if !a.IsSet(tc[i].bitpos) {
			t.Errorf("test case %d. Failed check #1", i)
		}
		if a.IsSet(tc[i].bitpos + 100) {
			t.Errorf("test case %d. Failed check #1", i)
		}

		if uint(len(a.mask)) != tc[i].eslicelen {
			t.Errorf("test case %d. Failed check #2", i)
		}

		if a.mask[tc[i].byteIndex] != tc[i].eByteValue {
			t.Errorf("test case %d. Failed check #3", i)
		}

		a.Reset(tc[i].bitpos)
		if a.IsSet(tc[i].bitpos) {
			t.Errorf("test case %d. Failed check #4", i)
		}
	}
}

func TestBitSet_SetMuilti(t *testing.T) {

	var tc = []struct {
		bitpos     uint // target bit pos
		eslicelen  uint // expected slice len
		byteIndex  int  // slice index where bitpos should belong to
		eByteValue byte // expected byte value
		allocated  bool
	}{
		{
			bitpos:     0,
			eslicelen:  1,
			byteIndex:  0,
			eByteValue: 1,
			allocated:  false,
		},
		{
			bitpos:     1,
			eslicelen:  1,
			byteIndex:  0,
			eByteValue: 3,
			allocated:  true,
		},
		{
			bitpos:     9,
			eslicelen:  2,
			byteIndex:  1,
			eByteValue: 2,
			allocated:  false,
		},
	}

	a := BitSet{}

	for i := range tc {
		if a.IsAllocated(tc[i].bitpos) != tc[i].allocated {
			t.Errorf("test case %d. Failed check #1", i)
		}

		a.Set(tc[i].bitpos)
		if !a.IsSet(tc[i].bitpos) {
			t.Errorf("test case %d. Failed check #2", i)
		}
		if uint(len(a.mask)) != tc[i].eslicelen {
			t.Errorf("test case %d. Failed check #3", i)
		}

		if a.mask[tc[i].byteIndex] != tc[i].eByteValue {
			t.Errorf("test case %d. Failed check #4", i)
		}

	}

}

func TestBitSet_AreSet(t *testing.T) {

	var tc = []struct {
		a       []uint // primary bit set
		b       []uint // slice to be compared with a
		eresult bool   // expected result
	}{
		{
			a:       []uint{0, 1, 2, 3, 4, 20},
			b:       []uint{0, 1, 4},
			eresult: true,
		},
		{
			a:       []uint{0, 1, 2, 3, 4, 20},
			b:       []uint{20},
			eresult: true,
		},
		{
			a:       []uint{0, 1, 2, 3, 4, 20},
			b:       []uint{0},
			eresult: true,
		},
		{
			a:       []uint{0, 1, 20, 40, 10},
			b:       []uint{},
			eresult: false,
		},
		{
			a:       []uint{},
			b:       []uint{0, 1, 4},
			eresult: false,
		},
		{
			a:       []uint{4, 5, 6, 21},
			b:       []uint{1, 5, 4},
			eresult: false,
		},
	}

	for i := range tc {

		a := BitSet{}
		for j := range tc[i].a {
			a.Set(tc[i].a[j])
		}

		if res := a.AreSet(tc[i].b...); res != tc[i].eresult {
			t.Errorf("test case %d. Failed check #1", i)

		}
	}
}
func TestBitSet_String(t *testing.T) {

	var tc = []struct {
		bitpos []uint // bit set
		exp    string // expected result
	}{
		{
			bitpos: []uint{0, 1, 2, 8},
			exp:    "0701",
		},
		{
			bitpos: []uint{},
			exp:    "",
		},
		{
			bitpos: []uint{0, 8, 16, 24},
			exp:    "01010101",
		},
		{
			bitpos: []uint{0, 1, 2, 3, 4, 5, 6, 7},
			exp:    "ff",
		},
	}

	for i := range tc {

		a := BitSet{}
		for j := range tc[i].bitpos {
			a.Set(tc[i].bitpos[j])
		}
		if s := a.String(); s != tc[i].exp {
			t.Errorf("test case %d. Failed check #1. Expected %s, got %s", i, tc[i].exp, s)
		}

	}
}
func TestParse(t *testing.T) {

	var tc = []struct {
		bitpos []uint
	}{
		{
			bitpos: []uint{1, 2, 8, 17},
		},
		{
			bitpos: []uint{81},
		},
		{
			bitpos: []uint{100},
		},
	}

	for i := range tc {
		a := BitSet{}
		for j := range tc[i].bitpos {
			a.Set(tc[i].bitpos[j])
		}
		s := a.String()
		b, err := Parse([]byte(s))
		if err != nil {
			t.Errorf("test case %d. Failed check #1. Parsing %s failed", i, s)
		}
		if s != b.String() {
			t.Errorf("test case %d. Failed check #2. Expected %s %v, got %s %v", i, s, a, b.String(), b)
		}

		if set, err := AreSet([]byte(s), tc[i].bitpos...); err != nil || !set {
			t.Errorf("test case %d. Failed check #3.", i)
		}
	}

	bs, err := Parse([]byte{})
	if bs.Len() > 0 || err != nil {
		t.Error("parsing empty string failed")
	}

}

func TestNew(t *testing.T) {
	var tc = []struct {
		size uint
		elen int // expected slice len in bytes
	}{
		{
			size: 1,
			elen: 1,
		},
		{
			size: 8,
			elen: 1,
		},
		{
			size: 15,
			elen: 2,
		},
		{
			size: 65,
			elen: 9,
		},
	}

	for i := range tc {
		a := New(tc[i].size)
		if len(a.mask) != tc[i].elen {
			t.Errorf("test case %d. Failed check #1", i)
		}
	}
}

// e65fff7f2feec3efbc7dffbfdcf3f7ff3f9ffdffff7f75bd01

func Test_AreSet(t *testing.T) {
	var tc = []struct {
		src    []byte
		bitpos uint
		exp    bool
	}{
		{
			src:    []byte(`e65f`),
			bitpos: 7,
			exp:    true,
		},
	}

	for i := range tc {
		res, err := AreSet(tc[i].src, tc[i].bitpos)
		if err != nil {
			t.Errorf("test case %d. Failed check #1. Parsing failed %s", i, string(tc[i].src))
		}
		if res != tc[i].exp {
			t.Errorf("test case %d. Failed check #1", i)
		}
	}
}
