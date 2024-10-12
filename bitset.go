// The package provides typical bitmap storage and manipulation functions.
// BitSet is a data structure that holds a set of bits. Each bit is represented by a single bit in the byte slice.
// The first bit is stored in the left bit of the first byte,
// and the last bit is stored in the least significant bit of the last byte.
// Examples:
//
// 1000 0010 where bits 0,6 is set in the string representation will look like hex string "82"
// 1000 0011 where bits 0, 6, 7 is set in the string representation will look like hes string "83"
// 1100 0001 1100 0011 where bits 0,1,7, 8,9, 14,15 us set in the string representation will look like "c1c3"
package bitset

import (
	"errors"
	"fmt"
)

// ErrParseFailed is returned by the Parse function when an invalid character is found
// in the input string. Valid characters are 0-9, a-f, A-F.
var (
	ErrParseFailed         = errors.New("invalid character found")
	ErrInvalidSourceString = errors.New("invalid source string")
)

var (
	invalidSymbol = byte('_')
	nibbleMapping = [16]byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
	}
)

type CompareRule int8

const (
	All CompareRule = iota
	Any
)

type BitSet interface {
	IsSet(pos uint) bool
	AreSet(rule CompareRule, bitpos ...uint) bool
	Set(val bool, bitpos ...uint) BitSet
	Empty() bool
	String() string
	Bytes() []byte
}

// ByteBitSet holds a set of bits using a slice of bytes.
// Each bit in the bitset is represented by one bit in the byte slice.
type ByteBitSet struct {
	mask []uint8
}

// New creates a BitSet with allocated space for the specified number of bits.
// It takes the bit size as input, calculates how many bytes are needed to store
// that many bits, and returns a pointer to a new BitSet.
func New(bitSize int) *ByteBitSet {
	byteSize := bitSize / 8
	if bitSize%8 > 0 {
		byteSize++
	}
	return &ByteBitSet{mask: make([]uint8, byteSize)}
}

func NewFromString(src string) (ByteBitSet, error) {
	if len(src) != len([]rune(src)) {
		return ByteBitSet{}, ErrInvalidSourceString
	}
	return newFromBytes([]byte(src))
}

// NewFromBytes converts a hexadecimal byte slice into a BitSet.
// It returns an error if any of the input characters are not valid hexadecimal digits.
func NewFromBytes(buf []byte) (ByteBitSet, error) {
	return newFromBytes(buf)
}

func newFromBytes(buf []byte) (ByteBitSet, error) {
	var bs ByteBitSet
	if len(buf) == 0 {
		return bs, nil
	}

	if len(buf)%2 == 1 {
		return bs, ErrInvalidSourceString
	}

	bs.mask = make([]byte, len(buf)/2)

	for i := 0; i < len(buf); i += 2 {

		bm, err := parsePair([2]byte(buf[i : i+2]))
		if err != nil {
			return bs, err
		}

		bs.mask[i/2] = bm
	}
	return bs, nil

}

// Clone returns a deep copy of the given BitSet.
// It creates a new BitSet and copies the mask from the source BitSet.
func Clone(bs BitSet) BitSet {
	return &ByteBitSet{mask: append([]byte{}, bs.Bytes()...)}
}

// Set sets the specified bits to 1. If the bit position is out of range, the internal
// byte slice is automatically extended to accommodate the bit.
func (bs *ByteBitSet) Set(val bool, bitpos ...uint) BitSet {
	for _, u := range bitpos {
		bs.set(val, u)
	}
	return bs
}

// set is a helper function that sets or resets the bit at the specified position.
// It extends the internal byte slice if necessary, and modifies the specified bit.
func (bs *ByteBitSet) set(isSet bool, bitpos uint) {
	idx := bitpos / 8
	pos := uint8(bitpos % 8)
	n := idx
	l := uint(len(bs.mask))

	// Extend internal storage if needed
	switch {
	case n == 0 && l == 0:
		(*bs).mask = []byte{0}
	case n >= l:
		(*bs).mask = append((*bs).mask, make([]byte, n-l+1)...)
	}

	// Set or reset the bit at the specified position
	if isSet {
		bs.mask[idx] |= 1 << (7 - pos)
	} else {
		bs.mask[idx] ^= 1 << (7 - pos)
	}
}

// IsSet returns true if the bit at the specified position is set (i.e., 1).
// If the position is beyond the current length of the bitset, it returns false.
func (bs *ByteBitSet) IsSet(bitpos uint) bool {
	if bitpos >= bs.Len() {
		return false
	}

	idx := bitpos / 8
	pos := 7 - bitpos%8

	return (uint(bs.mask[idx]))>>pos&1 == 1
}

// IsAllocated returns true if the space for the specified bit is already allocated.
// It checks whether the internal storage has enough space for the given bit position.
func (bs *ByteBitSet) IsAllocated(bitpos uint) bool {
	return len(bs.mask) > 0 && bitpos < bs.Len()
}

// Len returns the length of the allocated bitset in bits.
// The length is calculated by multiplying the number of bytes in the mask by 8.
func (bs *ByteBitSet) Len() uint {
	return uint(len(bs.mask) * 8)
}

// AreSet checks if all or any of the specified bits are set to 1.
// If any of the bits is not set, or if the bitset is empty, it returns false.
func (bs *ByteBitSet) AreSet(rule CompareRule, bitpos ...uint) bool {
	if len(bitpos) == 0 {
		return false
	}

	if len(bs.mask) == 0 {
		return false
	}

	if rule == All {
		for _, pos := range bitpos {
			if !bs.IsSet(pos) {
				return false
			}
		}
		return true
	}

	// if rule == Any
	for _, pos := range bitpos {
		if bs.IsSet(pos) {
			return true
		}
	}
	return false
}

// Bytes returns the internal representation of the bitset as a byte slice.
func (bs *ByteBitSet) Bytes() []byte {
	return bs.mask
}

// String returns a hexadecimal representation of the bitset.
// Each byte in the bitset is converted to two hexadecimal characters.
func (bs *ByteBitSet) String() string {
	if len(bs.mask) == 0 {
		return ""
	}

	res := make([]byte, len(bs.mask)*2)
	for i := range bs.mask {
		left := (bs.mask[i] >> 4)
		right := (bs.mask[i] &^ 0b11110000)
		res[i*2] = nibbleMapping[left]
		res[i*2+1] = nibbleMapping[right]
	}
	return string(res)
}

// BinaryString returns binary representation of the bitset.
// Each byte in the bitset is converted to 8 characters (0 or 1).
// Zero bit of the bitset will be at the end of the string.
func (bs *ByteBitSet) BinaryString() string {
	if len(bs.mask) == 0 {
		return ""
	}

	var res string
	for _, c := range bs.mask {
		res += fmt.Sprintf("%08b", c)
	}
	return string(res)
}

func NewFromBinaryString(src string) (BitSet, error) {
	var bs ByteBitSet
	if len(src) == 0 {
		return &bs, nil
	}

	for i, c := range []rune(src) {
		if !(c == '0' || c == '1') {
			return nil, ErrParseFailed
		}
		bs.Set(src[i] == '1', uint(i))
	}
	return &bs, nil
}

func shiftByte(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	}
	return invalidSymbol
}

func parsePair(b [2]byte) (uint8, error) {

	var res uint8

	high := shiftByte(b[0])
	if high == invalidSymbol {
		return 0, ErrParseFailed
	}

	low := shiftByte(b[1])
	if low == invalidSymbol {
		return 0, ErrParseFailed
	}
	res = high*16 | low
	return res, nil
}

func Validate(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}

	if len(buf)%2 == 1 {
		return ErrParseFailed
	}

	for i := 0; i < len(buf); i += 2 {
		if _, err := parsePair([2]byte(buf[i : i+2])); err != nil {
			return ErrParseFailed
		}
	}
	return nil
}

// AreSet receives a hexadecimal string representation of a BitSet and checks
// if all specified bits are set to 1.
func AreSet(buf []byte, rule CompareRule, bitpos ...uint) (bool, error) {

	if len(buf) == 0 {
		return false, nil
	}

	for _, pos := range bitpos {
		bytePos := pos / 8
		inBytePos := 7 - pos%8
		// if inBytePos != 0 && pos >= 8 {
		// 	bytePos++
		// }

		val := buf[bytePos] >> byte(inBytePos) & 1
		if rule == All && val == 0 {
			return false, nil
		}

		if rule == Any && val == 1 {
			return true, nil
		}
	}

	if rule == All {
		return true, nil
	}

	// if rule == Any

	return false, nil
}

// Empty returns true if BitSet is not initialized.
func (bs *ByteBitSet) Empty() bool {
	return len(bs.mask) == 0
}
