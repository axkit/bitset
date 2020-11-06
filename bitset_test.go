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
		ecnt       uint // expected max cnt
		byteIndex  int  // slice index where bitpos should belong to
		eByteValue byte // expected byte value
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
			eByteValue: 3,
		},
		{
			bitpos:     9,
			eslicelen:  2,
			ecnt:       10,
			byteIndex:  1,
			eByteValue: 2,
		},
	}

	a := BitSet{}

	for i := range tc {
		if a.InSet(tc[i].bitpos) {
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

		if a.Len() != tc[i].ecnt {
			t.Errorf("test case %d. Failed check #5", i)
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
