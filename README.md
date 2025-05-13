# bitset

[![Build Status](https://github.com/axkit/bitset/actions/workflows/go.yml/badge.svg)](https://github.com/axkit/bitset/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axkit/bitset)](https://goreportcard.com/report/github.com/axkit/bitset)
[![GoDoc](https://pkg.go.dev/badge/github.com/axkit/bitset)](https://pkg.go.dev/github.com/axkit/bitset)
[![Coverage Status](https://coveralls.io/repos/github/axkit/bitset/badge.svg?branch=main)](https://coveralls.io/github/axkit/bitset?branch=main)

`bitset` is a lightweight Go package that provides a typical bit set data structure using a `[]byte` slice. It allows fast and compact storage of bits, with utility methods for checking, setting, cloning, parsing, and serializing bitsets. Bits are stored from the most significant bit (MSB) of each byte.

## Motivation

The main motivation behind this package was to provide an efficient and compact representation of user permissions within JWT tokens, avoiding repeated calls to an external authorization service.

In a typical authorization system, each permission (like `read_user`, `edit_order`, etc.) is assigned a unique identifier that corresponds to a bit position in a bitmask. 

Let’s assume that perms and roles are stored in the database. 
All available permissions are stored in a `permissions` table:

```sql
CREATE TABLE permissions (
    id      TEXT    PRIMARY KEY,
    title   TEXT    NOT NULL,
    bit_pos INTEGER NOT NULL
);

INSERT INTO permissions (id, title, bit_pos) VALUES
('read_user',       'Read user data',            0),
('edit_user',       'Edit user profile',         1),
('delete_user',     'Delete user account',       2),
('create_order',    'Create an order',           3),
('edit_order',      'Modify existing order',     4),
('cancel_order',    'Cancel order',              5),
('manage_payments', 'Access payment settings',   6),
('admin_panel',     'Access admin panel',        7);
```

User roles (e.g., `user`, `admin`) are defined in a `roles` table, each containing an array of permission IDs:

```sql
CREATE TABLE roles (
    id              INT4    PRIMARY KEY,
    name            TEXT    NOT NULL,
    permission_ids  TEXT[]  NOT NULL
);

INSERT INTO roles (id, name, permission_ids) VALUES
(1, 'Regular User',        ARRAY['read_user', 'create_order']),
(2, 'Moderator',           ARRAY['read_user', 'edit_user', 'edit_order']),
(3, 'Administrator',       ARRAY['read_user', 'edit_user', 'delete_user', 'create_order', 'edit_order', 'cancel_order', 'manage_payments']),
(4, 'Super Admin',         ARRAY['read_user', 'edit_user', 'delete_user', 'create_order', 'edit_order', 'cancel_order', 'manage_payments', 'admin_panel']);
```

When a JWT token is generated, the list of permissions associated with a role is converted into a bitmask. This bitmask is then serialized as a **hexadecimal string** and embedded into the JWT payload.

**Example**:  
Permissions: `read_user`, `create_order` → Bit positions: 0 and 3 → Bitmask: `10010000` → Hex string: `"90"`

### Why this is useful 

- **Compact token size** — bitmask in hex is far smaller than an array of strings
- **Instant permission checks** — no need to call external services
- **Unified permission model** — easily shared across distributed services
- **Minimal HTTP overhead** — especially in the `Authorization` header

This package provides robust tools for encoding, decoding, and validating such bitmasks — making it ideal for use in permission-based access control scenarios.

## Installation

To install the package, run:

```bash
go get github.com/axkit/bitset
```

## Features

- **Create BitSet**: Initialize a new bitset with a specified size or from a hexadecimal or binary string.
- **Set/Unset Bits**: Set or unset bits at specified positions.
- **Check Bits**: Check if specific bits are set using either `All` or `Any` rules.
- **Convert to String**: Convert the bitset to hexadecimal or binary string formats.
- **Clone**: Create a deep copy of an existing bitset.

## Example

```go
package main

import (
    "fmt"
    "github.com/axkit/bitset"
)

func main() {
    bs := bitset.New(8)
    bs.Set(true, 0, 6)

    fmt.Println("Is bit 6 set?", bs.IsSet(6))                               // true
    fmt.Println("Binary representation:", bs.BinaryString())                // 10000010
    fmt.Println("Hex representation:", bs.String())                         // "82" 
    fmt.Println("Are bits 0 and 6 set?", bs.AreSet(bitset.All, 0, 6))       // true
    fmt.Println("Are bits 1, 6 or 7 set?", bs.AreSet(bitset.Any, 1, 6, 7))  // true
}
```

### BitSet Interface

> Using `uint64` instead of `byte` and leveraging SIMD instructions can speed up bit array operations, which is why the BitSet interface is designed to allow for future optimizations.

```go
type BitSet interface {
    IsSet(bit uint) bool
    AreSet(rule CompareRule, bits ...uint) bool
    Set(val bool, bits ...uint)
    Len() uint              
    String() string         
    BinaryString() string   
}
```

## Usage

### Creating a New BitSet

You can create a new `BitSet` by specifying the number of bits to allocate:

```go
bs := bitset.New(16) // Creates a BitSet with space for 16 bits
```

Alternatively, you can initialize it from a hexadecimal string:

```go
bs, err := bitset.ParseString("1с")
```

### Setting and Checking Bits

```go
bs.Set(true, 0, 6) // Sets bits 0 and 6
if bs.IsSet(6) {
    fmt.Println("Bit 6 is set")   
}
```

### Checking Multiple Bits

You can check if all or any bits are set using the `All` or `Any` rules:

```go
if bs.AreSet(bitset.All, 0, 6) {
    fmt.Println("Bits 0 and 6 are both set")
}

if bs.AreSet(bitset.Any, 0, 7) {
    fmt.Println("Bit 0 is set")
}
```

### Quick Checking

This feature allows you to validate bits directly from a hexadecimal bitmask without creating intermediate objects. It's designed for scenarios where performance and simplicity are critical, such as validating permissions in a compact format.

```go
// perms =  "a1ff90e428c7"

isAllowed, err := bitset.AreSet(perms, bitset.Any, 13, 42, 89)
```

This approach minimizes overhead and ensures rapid permission checks, making it ideal for high-performance applications.

### Converting to String Representations

To get the hexadecimal or binary string representation of the bitset:

```go
hexStr := bs.String()
fmt.Println("Hexadecimal:", hexStr)

binaryStr := bs.BinaryString()
fmt.Println("Binary:", binaryStr)

```

## Errors

- `ErrInvalidSourceString` - an invalid source string, such as an odd-length string in hexadecimal input.
- `ErrParseFailed` - Returned by parsing functions when an invalid character is encountered.

## Licensing

This package is open-source and distributed under the MIT License. Contributions and feedback are welcome!
