# bitset [![GoDoc](https://godoc.org/github.com/axkit/bitset?status.svg)](https://godoc.org/github.com/axkit/bitset) [![Build Status](https://travis-ci.org/axkit/bitset.svg?branch=master)](https://travis-ci.org/axkit/bitset) [![Coverage Status](https://coveralls.io/repos/github/bitset/gonfig/badge.svg)](https://coveralls.io/github/axkit/bitset) [![Go Report Card](https://goreportcard.com/badge/github.com/axkit/bitset)](https://goreportcard.com/report/github.com/axkit/bitset)

A simple bit set with JSON support

# Motivation

The package built specially to be used in package [github.com/axkit/aaa](https://github.com/axkit/aaa) as a JWT permissions holder but can be
used independently.

## Concepts

- Application functionality can be limited by using permissions.
- Permission (access right) represented by unique string code.
- Application supports many permissions.
- A user has a role.
- A role is set of allowed permission, it's subset of all permissions supported by application.
- As a result of successful sign in, a backend provides access and refresh tokens.
- The payload of access token have list of allowed permissions.
- A single permission code looks like "Customers.Create", "Customer.AttachDocuments", "Customer.Edit", etc.
- Store allowed permission codes could increase token size.
- Bitset comes here.
- Every permission shall be associated with a single bit in the set.
- Bitset adds to the token as hexadecimal string. Every 8 permissions represented by 2 characters.

## Usage Examples

Sign In

```
    var perms bitset.Bitset
    perms.Set(1)                    // 0000_0010
    perms.Set(2)                    // 0000_0110
    perms.Set(8, 10)                // 0000_0110 0000_0101
    tokenPerms := perms.String()    // returns "0605" as hex repsesentation of 0000_0110 0000_0101
```

Check allowed permission in auth middleware

```
    ...
    tokenPerms := accessToken.Payload.Perms     // "0605"
    bs, err := bitset.Parse(tokenPerms)         // returns 0000_0110 0000_0101
    if bs.AreSet(2,8) {
        // the permission allowed
    }
```

# Further Improvements

- [ ] Finalize integration BitSet with database/sql
- [ ] Add benchmarks
- [ ] Reduce memory allocations

Prague 2020
