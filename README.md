[![Build Status](https://github.com/azazeal/dchash/actions/workflows/build.yml/badge.svg)](https://github.com/azazeal/dchash/actions/workflows/build.yml)
[![Coverage Report](https://coveralls.io/repos/github/azazeal/dchash/badge.svg?branch=master)](https://coveralls.io/github/azazeal/dchash?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/azazeal/dchash.svg)](https://pkg.go.dev/github.com/azazeal/dchash)

# dchash

Package `dchash` implements a configurable variation of Dropbox's
[Content Hash](https://www.dropbox.com/developers/reference/content-hash)
algorithm in idiomatic Go.

## Usage

```go
package checksum

import (
	"crypto/sha512"
	"io"

	"github.com/azazeal/dchash"
)

// Dropbox hashes the data in src as per Dropbox's Content Hash algorithm, using
// SHA256 and over blocks of 4 MiB in size.
func Dropbox(src io.Reader) (sum []byte, err error) {
	h := dchash.New(nil, -1)

	if _, err = io.Copy(h, src); err == nil {
		sum = h.Sum(nil)
	}

	return
}

// SHA512m1 hashes the data in src as per Dropbox's Content Hash algorithm but
// with SHA512 (instead of SHA256) and over blocks of 1 (instead  of 4) MiB in
// size.
func SHA512m1(src io.Reader) (sum []byte, err error) {
	h := dchash.New(sha512.New, 1<<20)

	if _, err = io.Copy(h, src); err == nil {
		sum = h.Sum(nil)
	}

	return
}
```
