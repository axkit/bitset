# BitSet


---

## Installation

```bash
go get github.com/yourusername/bitset
```

---

## Motivation


## Features

- Compact representation using `[]byte`
- MSB-first bit ordering
- Efficient bit manipulation (`Set`, `IsSet`, etc.)
- Logical comparison of multiple bits using `All` or `Any` rules
- Conversions to/from:
  - Hexadecimal strings (`String()`)
  - Binary strings (`BinaryString()`)
  - Raw byte slices (`Bytes()`)
- Robust parsing and validation functions
- Automatically resizes when setting bits beyond current capacity

---

## ✨ Usage Example

```go
package main

import (
    "fmt"
    "github.com/yourusername/bitset"
)

func main() {
    bs := bitset.New(16)
    bs.Set(true, 0, 6, 15)

    fmt.Println("Is bit 6 set?", bs.IsSet(6))               // true
    fmt.Println("Binary representation:", bs.BinaryString()) // 10000010...
    fmt.Println("Hex representation:", bs.String())         // "82" or similar
    fmt.Println("Are bits 0 and 6 set?", bs.AreSet(bitset.All, 0, 6)) // true
}
```

---

## 🧠 API Overview

### Types

```go
type CompareRule int8

const (
    All CompareRule = iota // All specified bits must be set
    Any                    // At least one of the specified bits must be set
)
```

---

### BitSet Interface

```go
type BitSet interface {
    IsSet(pos uint) bool
    AreSet(rule CompareRule, bitpos ...uint) bool
    Set(val bool, bitpos ...uint) BitSet
    Len() uint
    String() string      // Returns hex string
    Bytes() []byte       // Returns raw byte slice
}
```

---

### Core Functions

| Function | Description |
|---------|-------------|
| `New(size int)` | Creates a new `ByteBitSet` with enough space for `size` bits |
| `ParseString(s string)` | Parses a bitset from a hex string |
| `ParseBinaryString(s string)` | Parses a bitset from a binary `0`/`1` string |
| `ParseBytes([]byte)` | Parses a bitset from a hex-encoded byte slice |
| `Clone(bs ByteBitSet)` | Returns a deep copy of a bitset |
| `AreSet(hexStr string, rule CompareRule, bits ...uint)` | Compares bits directly on a hex string |
| `Validate(buf []byte)` | Validates whether a byte slice is a correct hex representation |

---

## 🧪 Testing

A comprehensive test suite is included and covers:

- Basic functionality (`Set`, `IsSet`, `AreSet`, `Clone`)
- Edge cases (out-of-bounds, empty inputs, invalid characters)
- String representations (hex and binary)
- Validation and parsing logic

Run tests using:

```bash
go test -v ./...
```

---

## 📜 License

MIT — use freely for any purpose.
