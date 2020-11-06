package bitset

import (
	"database/sql/driver"
	"fmt"
)

// BitSet holds set of bits using slice of bytes.
type BitSet struct {
	mask []uint8
	cnt  uint
}

// Set sets bit with position bitpos to 1. Extends internal storage if it's required.
func (bs *BitSet) Set(bitpos uint) {
	bs.set(bitpos, true)
}

// Reset sets bit with position bitpos to 0. Does not resize internal storage.
func (bs *BitSet) Reset(bitpos uint) {
	bs.set(bitpos, false)
}

func (bs *BitSet) set(bitpos uint, isSet bool) {

	idx := bitpos / 8
	pos := bitpos % 8
	n := idx
	l := uint(len(bs.mask))
	switch {
	case n == 0 && l == 0:
		// mask is empty and set comes to the bit 0..7
		(*bs).mask = []byte{0}
	case n >= l:
		(*bs).mask = append((*bs).mask, make([]byte, n-l+1)...)
	}

	if bitpos >= bs.cnt {
		bs.cnt = bitpos + 1
	}

	if isSet {
		bs.mask[idx] |= (1 << pos)
	} else {
		bs.mask[idx] ^= (1 << pos)
	}

}

// IsSet returns true if bit with position bitpos is 1.
// Returns false if bitpos above maximal setted bitpos.
func (bs *BitSet) IsSet(bitpos uint) bool {
	if bitpos > bs.cnt {
		return false
	}
	idx := bitpos / 8
	pos := bitpos % 8
	return bs.mask[idx]&(1<<pos) == 1<<pos
}

// InSet returns true if space for the bit with position bitpos is allocated.
func (bs *BitSet) InSet(bitpos uint) bool {
	return len(bs.mask) > 0 && bitpos < bs.cnt
}

// Len return max bitpos stored in the set.
func (bs *BitSet) Len() uint {
	return bs.cnt
}

// AreSet returns true if every bit with position bit pos is equal 1.
// Return false if bitops or bs are empty.
func (bs *BitSet) AreSet(bitpos ...uint) bool {
	if len(bitpos) == 0 {
		return false
	}
	if bs.cnt == 0 {
		return false
	}

	for _, pos := range bitpos {
		if !bs.IsSet(pos) {
			return false
		}
	}
	return true
}

// Value implements database/sql Valuer.
func (bs BitSet) Value() (driver.Value, error) {
	if bs.mask == nil {
		return nil, nil
	}

	b := make([]byte, bs.cnt)
	for i := bs.cnt; i >= 0; i-- {
		b[i-1] = map[bool]byte{true: '1', false: '0'}[bs.IsSet(i-1)]
	}

	return b, nil
}

// Valid returns true if at least one bit is set.
func (bs *BitSet) Valid() bool {
	return bs.cnt > 0
}

// Scan implements database/sql Scanner. It's expected that
// PostgresSQL type BIT VARYING is used.
func (bs *BitSet) Scan(value interface{}) error {
	if value == nil {
		(*bs).mask = nil
		return nil
	}

	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("BitSet.Scan: expected []byte or string, got %T (%q)", value, value)
	}

	for i := len(b); i >= 0; i-- {
		bs.Set(uint(i))
	}

	return nil
}
