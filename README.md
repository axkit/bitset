# bitset package

[![Build Status](https://github.com/axkit/bitset/actions/workflows/go.yml/badge.svg)](https://github.com/axkit/bitset/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axkit/bitset)](https://goreportcard.com/report/github.com/axkit/bitset)
[![GoDoc](https://pkg.go.dev/badge/github.com/axkit/bitset)](https://pkg.go.dev/github.com/axkit/bitset)
[![Coverage Status](https://coveralls.io/repos/github/axkit/bitset/badge.svg?branch=main)](https://coveralls.io/github/axkit/bitset?branch=main)

The `bitset` package provides a simple and efficient implementation of a bitmap data structure in Go, allowing you to store and manipulate individual bits. It is ideal for use cases where you need to work with large sets of bits and perform operations such as checking, setting, and clearing bits.

## Installation

To install the package, use:

```bash
go get github.com/axkit/bitset
```

## Overview

The bitset package provides the BitSet interface and its implementation ByteBitSet, which stores bits as a slice of bytes. Each bit is represented by a single bit within a byte, with the leftmost bit in the first byte representing the first bit, and the least significant bit in the last byte representing the last bit.

## Features

- Create BitSet: Initialize a new bitset with a specified size or from a hexadecimal or binary string.
- Set/Unset Bits: Set or unset bits at specified positions.
- Check Bits: Check if specific bits are set, using both All and Any rules.
- Check Bits: Check if specific bits are set, using both All and Any rules.
- Convert to String: Convert the bitset to hexadecimal or binary string formats.
- Clone: Create a deep copy of an existing bitset.

## Examples

### Bit Representation
```
1000 0010  ->  Hexadecimal: "82" (bits 0 and 6 are set)
1000 0011  ->  Hexadecimal: "83" (bits 0, 6, and 7 are set)
1100 0001 1100 0011  ->  Hexadecimal: "c1c3" (bits 0, 1, 7, 8, 9, 14, 15 are set)
```

## Usage

### Creating a New BitSet

You can create a new BitSet by specifying the number of bits to allocate:

```go
import "github.com/axkit/bitset"

bs := bitset.New(16) // Creates a BitSet with space for 16 bits
```

Alternatively, you can initialize it from a hexadecimal string:
```go
bs, err := bitset.NewFromString("82")
if err != nil {
    fmt.Println("Error:", err)
}
```

### Setting and Checking Bits

Set specific bits:
```go
bs.Set(true, 0, 6) // Sets bits 0 and 6
```

Check if a bit is set:
```go
if bs.IsSet(0) {
    fmt.Println("Bit 0 is set")
}
```

### Checking Multiple Bits
You can check if all or any bits are set using All or Any rules:
```go
if bs.AreSet(bitset.All, 0, 6) {
    fmt.Println("Bits 0 and 6 are both set")
}

if bs.AreSet(bitset.Any, 0, 7) {
    fmt.Println("At least one of the bits 0 or 7 is set")
}
```

### Converting to String Representations
To get the hexadecimal or binary string representation of the bitset:

```go
hexStr := bs.String()
fmt.Println("Hexadecimal:", hexStr)

binStr := bs.BinaryString()
fmt.Println("Binary:", binStr)
```

### Cloning a BitSet

You can create a deep copy of an existing BitSet:
```go
clone := bitset.Clone(bs)
```
## Errors

`ErrParseFailed`

Returned by parsing functions when an invalid character is encountered.

`ErrInvalidSourceString`

Indicates an invalid source string, such as an odd-length string in hexadecimal input.

## License

This package is open-source and distributed under the MIT License. Contributions and feedback are welcome!

