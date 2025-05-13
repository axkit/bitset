// The package provides typical bitmap storage and manipulation functions.
// BitSet is a data structure that holds a set of bits, where each bit is represented by a single bit in a byte slice.
// The first bit is stored in the most significant bit of the first byte,
// and the last bit is stored in the least significant bit of the last byte.
// Examples:
//
// 1000 0010: Bits 0 and 6 are set. The hexadecimal representation is "82".
// 1000 0011: Bits 0, 6, and 7 are set. The hexadecimal representation is "83".
// 1100 0001 1100 0011: Bits 0, 1, 7, 8, 9, 14, and 15 are set. The hexadecimal representation is "c1c3".
package bitset

import (
	"errors"
	"unsafe"
)

// ErrParseFailed is returned by the Parse function when an invalid character is found
// in the input string. Valid characters are 0-9, a-f, A-F.
var (
	ErrParseFailed         = errors.New("invalid character found")
	ErrInvalidSourceString = errors.New("invalid source string")
)

var (
	invalidNumber = byte('_')
	nibbleMapping = [16]byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
	}
)

// CompareRule defines the rule for comparing bits in a BitSet.
// It is used to determine whether all or any of the specified bits are set.
//
// The following constants are available:
//   - All: All specified bits must be set.
//   - Any: At least one of the specified bits must be set.
type CompareRule int8

const (
	All CompareRule = iota
	Any
)

// BitSet defines an interface for manipulating a set of bits.
// It provides methods to check, set, retrieve, and convert the bit set.
type BitSet interface {
	IsSet(pos uint) bool
	AreSet(rule CompareRule, bitpos ...uint) bool
	Set(val bool, bitpos ...uint)
	Len() uint
	String() string
	BinaryString() string
	Bytes() []byte
}

// ByteBitSet is a BitSet implementation that stores bits in a byte slice.
// Each bit is represented by a single bit in the byte array, starting from the most significant bit of the first byte.
type ByteBitSet struct {
	mask []uint8
}

var _ BitSet = (*ByteBitSet)(nil)

// New returns a new ByteBitSet with enough space to store the specified number of bits.
func New(size int) ByteBitSet {
	n := size / 8
	if size%8 > 0 {
		n++
	}
	return ByteBitSet{mask: make([]uint8, n)}
}

// ParseHexString creates a ByteBitSet from a hexadecimal string representation.
// Returns an error if the input string is invalid.
func ParseHexString(hexStr string) (ByteBitSet, error) {

	if len(hexStr) != len([]rune(hexStr)) {
		return ByteBitSet{}, ErrInvalidSourceString
	}
	buf := unsafe.Slice(unsafe.StringData(hexStr), len(hexStr))
	return parseHexBytes(buf)
}

// ParseHexBytes is a sugar function that creates a ByteBitSet from a byte slice.
func ParseHexBytes(hexStr []byte) (ByteBitSet, error) {
	return parseHexBytes(hexStr)
}

func parseHexBytes(src []byte) (ByteBitSet, error) {

	if len(src) == 0 {
		return ByteBitSet{}, nil
	}

	if len(src)%2 == 1 {
		return ByteBitSet{}, ErrInvalidSourceString
	}

	bbs := New(len(src) / 2 * 8)

	for i := 0; i < len(src); i += 2 {
		bm, err := parsePair(src[i], src[i+1])
		if err != nil {
			return ByteBitSet{}, err
		}
		bbs.mask[i/2] = bm
	}
	return bbs, nil
}

// ParseBinaryString creates a BitSet from a binary string of '0' and '1' characters.
// Returns an error if any characters other than '0' or '1' are found.
func ParseBinaryString(src string) (ByteBitSet, error) {

	if len(src) == 0 {
		return ByteBitSet{}, nil
	}

	bbs := New(len(src))

	for i, c := range []rune(src) {
		if !(c == '0' || c == '1') {
			return ByteBitSet{}, ErrParseFailed
		}
		bbs.Set(c == '1', uint(i))
	}
	return bbs, nil
}

// Clone returns a deep copy of the provided ByteBitSet.
func Clone(src ByteBitSet) ByteBitSet {
	dst := ByteBitSet{
		mask: make([]uint8, len(src.mask)),
	}
	copy(dst.mask, src.mask)
	return dst
}

// Set updates the bits at the specified positions to the given value (true to set, false to clear).
// Automatically expands the internal byte slice if necessary.
func (bbs *ByteBitSet) Set(val bool, bits ...uint) {
	for _, bit := range bits {
		bbs.set(val, bit)
	}
}

// set is a helper function that sets or resets the bit at the specified position.
// It extends the internal byte slice if necessary, and modifies the specified bit.
func (bbs *ByteBitSet) set(val bool, bit uint) {
	//bn := bit / 8
	//bitn := uint8(bit % 8)
	size := uint(len(bbs.mask))
	bn, bitn := offsets(bit)

	// Extend internal storage if needed
	switch {
	case bn == 0 && size == 0:
		(*bbs).mask = []byte{0}
	case bn >= size:
		(*bbs).mask = append(bbs.mask, make([]byte, bn-size+1)...)
	}

	// Set or reset the bit at the specified position
	if val {
		bbs.mask[bn] |= (1 << bitn)
	} else {
		bbs.mask[bn] &^= (1 << bitn)
	}
}

// IsSet returns true if the bit at the specified position is set to 1 or false if it is 0.
//
// If the position is out of bounds, it returns false. Use Len() to check the size of the BitSet.
func (bbs ByteBitSet) IsSet(bit uint) bool {
	if bit >= bbs.Len() {
		return false
	}

	bn, bitn := offsets(bit)
	return (uint(bbs.mask[bn]))>>bitn&1 == 1
}

// Len returns the total number of bits currently allocated in the bit set.
func (bbs ByteBitSet) Len() uint {
	return uint(len(bbs.mask) * 8)
}

// AreSet checks whether all or any of the specified bits are set, depending on the rule provided.
// Use IsSet to check a single bit.
func (bbs ByteBitSet) AreSet(rule CompareRule, bits ...uint) bool {
	if len(bits) == 0 {
		return false
	}

	if len(bbs.mask) == 0 {
		return false
	}

	if len(bits) == 1 {
		return bbs.IsSet(bits[0])
	}

	if rule == All {
		for _, bit := range bits {
			if !bbs.IsSet(bit) {
				return false
			}
		}
		return true
	}

	// if rule == Any
	for _, bit := range bits {
		if bbs.IsSet(bit) {
			return true
		}
	}
	return false
}

// Bytes returns the underlying byte slice representing the bitset.
func (bbs ByteBitSet) Bytes() []byte {
	return bbs.mask
}

// String returns the hexadecimal string representation of the bitset.
func (bbs ByteBitSet) String() string {
	if len(bbs.mask) == 0 {
		return ""
	}

	hex := make([]byte, len(bbs.mask)*2)
	for i := 0; i < len(bbs.mask); i++ {
		left := (bbs.mask[i] >> 4)
		right := (bbs.mask[i] &^ 0b11110000)
		hex[i*2] = nibbleMapping[left]
		hex[i*2+1] = nibbleMapping[right]
	}
	return *(*string)(unsafe.Pointer(&hex))
}

// BinaryString returns the binary string representation of the bitset.
// The leftmost bit corresponds to the lowest index.
func (bbs ByteBitSet) BinaryString() string {
	if len(bbs.mask) == 0 {
		return ""
	}

	bin := make([]byte, len(bbs.mask)*8)
	idx := 0
	for _, c := range bbs.mask {
		for i := 7; i >= 0; i-- {
			bin[idx] = '0' + byte((c&(1<<i))>>i) // 0 or 1
			idx++
		}
	}
	return *(*string)(unsafe.Pointer(&bin))
}

func hexCharToNumber(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	}
	return invalidNumber
}

func parsePair(first, second byte) (uint8, error) {

	var high, low byte

	if high = hexCharToNumber(first); high == invalidNumber {
		return 0, ErrParseFailed
	}

	if low = hexCharToNumber(second); low == invalidNumber {
		return 0, ErrParseFailed
	}

	return high<<4 | low, nil
}

// Validate quickly checks if the input byte slice is a valid hexadecimal representation of a BitSet.
func Validate(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}

	if len(buf)%2 == 1 {
		return ErrParseFailed
	}

	for i := 0; i < len(buf); i += 2 {
		if _, err := parsePair(buf[i], buf[i+1]); err != nil {
			return ErrParseFailed
		}
	}
	return nil
}

// AreSet evaluates whether all or any specified bits are set, based on the rule,
// using a hexadecimal string representation of the bitset.
func AreSet(hexStr string, rule CompareRule, bits ...uint) (bool, error) {

	n := len(hexStr)
	if n == 0 || len(bits) == 0 {
		return false, nil
	}

	if n%2 == 1 {
		return false, ErrInvalidSourceString
	}

	buf := unsafe.Slice(unsafe.StringData(hexStr), len(hexStr))
	hexBits := uint(len(buf) / 2 * 8)

	for _, bit := range bits {
		if bit >= hexBits {
			// if bit is out of bounds
			if rule == All {
				return false, nil
			}
			continue
		}

		bn, bitn := offsets(bit)
		byteVal, err := parsePair(buf[bn*2], buf[bn*2+1])
		if err != nil {
			return false, err
		}

		val := (byteVal >> bitn) & 1
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

	// we come here if rule equals to Any
	return false, nil
}

func offsets(bit uint) (uint, uint) {
	return bit / 8, 7 - bit%8
}
