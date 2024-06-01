indexers
========

[![Build Status](https://github.com/decred/dcrd/workflows/Build%20and%20Test/badge.svg)](https://github.com/decred/dcrd/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![Doc](https://img.shields.io/badge/doc-reference-blue.svg)](https://pkg.go.dev/github.com/decred/dcrd/internal/blockchain/indexers)

Package indexers implements optional block chain indexes.

These indexes are typically used to enhance the amount of information available
via an RPC interface.

## Supported Indexers

- Transaction-by-hash (txbyhashidx) Index
  - Creates a mapping from the hash of each transaction to the block that
    contains it along with its offset and length within the serialized block
- Address-ever-seen (existsaddridx) Index
  - Stores a key with an empty value for every address that has ever existed
    and was seen by the client

## Removed Legacy Indexers

- Committed Filter (cfindexparentbucket) Index
  - Stores all committed filters and committed filter headers for all blocks in
    the main chain
- Transaction-by-address (txbyaddridx) Index
  - Creates a mapping from every address to all transactions which either credit
    or debit the address
  - Requires the transaction-by-hash index

## Installation and Updating

This package is internal and therefore is neither directly installed nor needs
to be manually updated.

## License

Package indexers is licensed under the [copyfree](http://copyfree.org) ISC
License.
