package bitset

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

// ErrParseFailed returns if Parse function found character out of ranges
// 0..9, a..f, A..F
var ErrParseFailed = errors.New("invalid character found")

// BitSet holds set of bits using slice of bytes.
type BitSet struct {
	mask []uint8
}

// New creates BitSet with allocated space for size amount of bits.
func New(size uint) *BitSet {
	i := size / 8
	if size%8 > 0 {
		i++
	}
	return &BitSet{mask: make([]byte, i)}
}

// Clone clones bs.
func Clone(bs BitSet) BitSet {
	res := BitSet{}
	res.mask = make([]byte, len(bs.mask))
	copy(res.mask, bs.mask)
	return res
}

// Set sets bits to 1. Extends internal storage if it's required.
func (bs *BitSet) Set(bitpos ...uint) *BitSet {
	for _, u := range bitpos {
		bs.set(u, true)
	}
	return bs
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
		(*bs).mask = []byte{0}
	case n >= l:
		(*bs).mask = append((*bs).mask, make([]byte, n-l+1)...)
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
	if bitpos > bs.Len() {
		return false
	}
	idx := bitpos / 8
	pos := bitpos % 8
	return bs.mask[idx]&(1<<pos) == 1<<pos
}

// IsAllocated returns true if space for the bit with position bitpos is allocated already.
func (bs *BitSet) IsAllocated(bitpos uint) bool {
	return len(bs.mask) > 0 && bitpos < bs.Len()
}

// Len returns len of allocated byte slice.
func (bs *BitSet) Len() uint {
	return uint(len(bs.mask) * 8)
}

// AreSet returns true if every bit with position bit pos is equal 1.
// Return false if bitops or bs are empty.
func (bs *BitSet) AreSet(bitpos ...uint) bool {
	if len(bitpos) == 0 {
		return false
	}

	if len(bs.mask) == 0 {
		return false
	}

	for _, pos := range bitpos {
		if !bs.IsSet(pos) {
			return false
		}
	}
	return true
}

// String returns hex representation of bit array. Every 8 bits as 2 hex digits.
func (bs *BitSet) String() string {
	if len(bs.mask) == 0 {
		return ""
	}
	arr := [16]byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
	}
	res := make([]byte, len(bs.mask)*2)
	for i := range bs.mask {
		left := (bs.mask[i] >> 4)
		right := (bs.mask[i] &^ 0b11110000)
		res[i*2] = arr[left]
		res[i*2+1] = arr[right]
	}
	return string(res) //  string(fmt.Sprintf("%x", bs.mask))
}

// Parse converts []byte to BitSet.
func Parse(buf []byte) (BitSet, error) {
	var bs BitSet
	if len(buf) == 0 {
		return bs, nil
	}

	bs.mask = make([]byte, len(buf)/2)

	calc := func(c byte) byte {
		switch {
		case '0' <= c && c <= '9':
			return c - '0'
		case 'a' <= c && c <= 'f':
			return c - 'a' + 10
		case c >= 'A' && c <= 'F':
			return c - 'A' + 10
		}
		return 'Z'
	}

	for i := 0; i < len(buf); i += 2 {

		x := calc(buf[i])
		if x == 'Z' {
			return bs, ErrParseFailed

		}
		b := byte(x * 16)

		x = calc(buf[i+1])
		if x == 'Z' {
			return bs, ErrParseFailed
		}

		b |= (x)

		bs.mask[i/2] = b
	}
	return bs, nil
}

// AreSet recieves string representation of BitSet and returns true if
// every bit with position bitpos is equal 1.
func AreSet(buf []byte, bitpos ...uint) (bool, error) {
	bs, err := Parse(buf)
	if err != nil {
		return false, err
	}
	return bs.AreSet(bitpos...), nil
}

// Value implements database/sql Valuer.
func (bs BitSet) Value() (driver.Value, error) {
	if bs.mask == nil {
		return nil, nil
	}

	b := make([]byte, len(bs.mask)*8)
	for i := len(bs.mask) - 1; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			b[i] = map[bool]byte{true: '1', false: '0'}[bs.IsSet(uint(i*8+j))]
		}
	}

	return b, nil
}

// Valid returns true if at least one bit is set.
func (bs *BitSet) Valid() bool {
	return len(bs.mask) > 0
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
