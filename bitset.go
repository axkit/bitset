package bitset

import (
	"database/sql/driver"
	"fmt"
)

type BitSet uint64

// Intersect проверяет, что все биты b входят в множество bs.
func (bs *BitSet) Intersect(b *BitSet) bool {
	return (*bs & *b) == *bs
}

func (bs *BitSet) IsSet(bitpos uint64) bool {
	return (*bs)&(1<<bitpos) == 1<<bitpos
}

// Set устанавливает в 1 бита с позицией bitpos.
func (bs *BitSet) Set(bitpos uint64) {
	var v uint64 = 1 << bitpos
	*bs = BitSet(uint64(*bs) | v)
}

func (d BitSet) Value() (driver.Value, error) {

	return d, nil
}

func (d *BitSet) Valid() bool {
	return true
}

func (d BitSet) Len() int {
	return 63
}

// Scan имплементирует database/sql Scanner итнтерфейс.
func (d *BitSet) Scan(value interface{}) error {
	if value == nil {
		*d = 0
		return nil
	}

	v, ok := value.(int64)
	if !ok {
		return fmt.Errorf("BitSet.Scan: expected int64, got %T (%q)", value, value)
	}

	*d = BitSet(v)
	return nil
}
