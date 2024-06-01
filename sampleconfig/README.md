sampleconfig
============

[![Build Status](https://github.com/decred/dcrd/workflows/Build%20and%20Test/badge.svg)](https://github.com/decred/dcrd/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![Doc](https://img.shields.io/badge/doc-reference-blue.svg)](https://pkg.go.dev/github.com/decred/dcrd/sampleconfig)

Package sampleconfig provides functions that return the contents of sample
configuration files for dcrd and dcrctl.  This is provided for tools that
perform automatic configuration and would like to ensure the generated
configuration files not only include the specifically configured values, but
also provides samples of other configuration options.

## Installation and Updating

This package is part of the `github.com/decred/dcrd` module.  Use the standard
go tooling for working with modules to incorporate it.

## License

Package sampleconfig is licensed under the [copyfree](http://copyfree.org) ISC
License.
