package bitset_test

import (
	"testing"

	"github.com/axkit/bitset"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		size int
		want int
	}{
		{"Zero bits", 0, 0},
		{"One byte", 8, 1},
		{"Partial byte", 7, 1},
		{"Two bytes", 16, 2},
		{"More than a byte", 13, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := bitset.New(tt.size)
			if got := len(bs.Bytes()); got != tt.want {
				t.Errorf("New(%d) length = %d, want %d", tt.size, got, tt.want)
			}
		})
	}
}

func TestNewFromString(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantErr bool
	}{
		{"Empty string", "", false},
		{"Valid hex string", "f3", false},
		{"Invalid hex string", "fz", true},
		{"Uneven hex string", "abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := bitset.NewFromString(tt.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFromString(%q) error = %v, wantErr %v", tt.src, err, tt.wantErr)
			}
		})
	}
}

func TestSetAndIsSet(t *testing.T) {
	tests := []struct {
		name      string
		size      int
		setBits   []uint
		checkBits []uint
		expected  []bool
	}{
		{"Set bits in one byte", 8, []uint{0, 2, 7}, []uint{0, 2, 4, 7}, []bool{true, true, false, true}},
		{"Set bits across multiple bytes", 16, []uint{0, 9}, []uint{0, 9, 10}, []bool{true, true, false}},
		{"Set and unset bits", 16, []uint{2, 15}, []uint{2, 15}, []bool{true, true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := bitset.New(tt.size)
			for _, bit := range tt.setBits {
				bs.Set(true, bit)
			}
			for i, bit := range tt.checkBits {
				if got := bs.IsSet(bit); got != tt.expected[i] {
					t.Errorf("IsSet(%d) = %v, want %v", bit, got, tt.expected[i])
				}
			}
		})
	}
}

func TestAreSet(t *testing.T) {
	tests := []struct {
		name    string
		rule    bitset.CompareRule
		setBits []uint
		check   []uint
		want    bool
	}{
		{"All set", bitset.All, []uint{0, 2, 4}, []uint{0, 2, 4}, true},
		{"All not set", bitset.All, []uint{0, 2}, []uint{0, 1, 2}, false},
		{"Any set", bitset.Any, []uint{1, 3, 7}, []uint{0, 1, 4}, true},
		{"None set", bitset.Any, []uint{3, 6}, []uint{0, 2, 5}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := bitset.New(8)
			for _, bit := range tt.setBits {
				bs.Set(true, bit)
			}
			got, _ := bitset.AreSet(bs.Bytes(), tt.rule, tt.check...)
			if got != tt.want {
				t.Errorf("AreSet(%v, %v) = %v, want %v", tt.rule, tt.check, got, tt.want)
			}
		})
	}
}

func TestByteBitSet_StringAndBinaryString(t *testing.T) {
	tests := []struct {
		name       string
		size       int
		setBits    []uint
		wantHex    string
		wantBinary string
	}{

		{"Single byte:2 bits", 8, []uint{0, 6}, "82", "10000010"},
		{"Single byte:3 bits", 8, []uint{0, 6, 7}, "83", "10000011"},
		{"Single byte:4 bits", 8, []uint{0, 6, 7, 1}, "c3", "11000011"},
		{"Two byte:6 bits", 8, []uint{0, 1, 7, 8, 9, 14, 15}, "c1c3", "1100000111000011"},
		{"Empty", 0, nil, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := bitset.New(tt.size)
			bs.Set(true, tt.setBits...)
			if got := bs.String(); got != tt.wantHex {
				t.Errorf("String() = %v, want %v", got, tt.wantHex)
			}
			if got := bs.BinaryString(); got != tt.wantBinary {
				t.Errorf("BinaryString() = %v, want %v", got, tt.wantBinary)
			}
		})
	}
}

func TestByteBitSet_Set(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		setBits []uint
		wantHex string
	}{
		{"Single byte: empty", 0, []uint{0, 6}, "82"},
		{"Single byte:2 bits", 8, []uint{0, 6}, "82"},
		{"Single byte:3 bits", 8, []uint{0, 6, 7}, "83"},
		{"Single byte:4 bits", 8, []uint{0, 6, 7, 1}, "c3"},
		{"Two byte:6 bits", 8, []uint{0, 1, 7, 8, 9, 14, 15}, "c1c3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := bitset.New(tt.size).Set(true, tt.setBits...)
			if got := bs.String(); got != tt.wantHex {
				t.Errorf("String() = %v, want %v", got, tt.wantHex)
			}
		})
	}
}

func TestByteBitSet_Clone(t *testing.T) {
	bs := bitset.New(8).Set(true, 0, 1, 2)
	clone := bitset.Clone(bs)
	if &bs == &clone {
		t.Errorf("Clone() should return a new instance of BitSet")
	}
	if bs.String() != clone.String() {
		t.Errorf("Clone() should return a copy of the BitSet")
	}
}

func TestByteBitSet_AreSet(t *testing.T) {
	tests := []struct {
		name    string
		rule    bitset.CompareRule
		setBits []uint
		check   []uint
		want    bool
	}{
		{"All set", bitset.All, []uint{0, 2, 4}, []uint{0, 2, 4}, true},
		{"All not set", bitset.All, []uint{0, 2}, []uint{0, 1, 2}, false},
		{"Any set", bitset.Any, []uint{1, 3, 7}, []uint{0, 1, 4}, true},
		{"None set", bitset.Any, []uint{3, 6}, []uint{0, 2, 5}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := bitset.New(8)
			for _, bit := range tt.setBits {
				bs.Set(true, bit)
			}
			got := bs.AreSet(tt.rule, tt.check...)
			if got != tt.want {
				t.Errorf("AreSet(%v, %v) = %v, want %v", tt.rule, tt.check, got, tt.want)
			}
		})
	}
}

func TestByteBitSet_NewFromBinaryString(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantErr bool
	}{
		{"Empty string", "", false},
		{"Valid binary string", "10000010", false},
		{"Invalid binary string", "1000001z", true},
		{"Uneven binary string", "1000001", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, err := bitset.NewFromBinaryString(tt.src)
			if err == nil && len(tt.src) > 0 {
				if !bs.IsSet(0) {
					t.Errorf("NewFromBinaryString(%q) = %v, want %v", tt.src, bs.String(), tt.src)
				}
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFromBinaryString(%q) error = %v, wantErr %v", tt.src, err, tt.wantErr)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantErr bool
	}{
		{"Empty string", "", false},
		{"Valid hex string", "f3", false},
		{"Invalid hex string", "fz", true},
		{"Uneven hex string", "abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := bitset.Validate([]byte(tt.src)); (err != nil) != tt.wantErr {
				t.Errorf("Validate(%q) error = %v, wantErr %v", tt.src, err, tt.wantErr)
			}
		})
	}
}
