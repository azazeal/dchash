# dchash

Package `dchash` implements a configurable variation of Dropbox's [Content Hash](https://www.dropbox.com/developers/reference/content-hash) in idiomatic Go.

## Usage

```go

import (
    _ "crypto/sha512"

    "github.com/azazeal/dchash"
)

// dropbox hashes the data in the given Reader as per
// Dropbox's Content Hash (4 MiB blocks, SHA256).
func dropbox(src io.Reader) (sum []byte, err error) {
    h := dchash.New(nil, -1)

    if _, err = io.Copy(h, src); err == nil {
        sum = h.Sum(nil)
    }
    
    return
}

// sha512m1 hashes the data in the given Reader as per
// Dropbox's Content Hash  but using SHA512 (instead of SHA256)
// and 1 MiB block size (instead of 4MiB).
func sha512m1(src io.Reader) (sum []byte, err error) {
    h := dchash.New(sha512.New, 1<<20)

    if _, err = io.Copy(h, src); err == nil {
        sum = h.Sum(nil)
    }
    
    return
}
```