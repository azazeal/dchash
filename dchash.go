// Package dchash implements a configurable variation of Dropbox's Content Hash.
package dchash

import (
	"crypto"
	_ "crypto/sha256" // DefaultHash
	"hash"
)

const (
	// DefaultHash denotes the default Hash.
	DefaultHash = crypto.SHA256

	// DefaultBlockSize denotes the default block size.
	DefaultBlockSize = 1 << 22
)

// New returns a Hash which yields checksums according to Dropbox's Content Hash
// algorithm, more details on which may be found here:
//
// https://www.dropbox.com/developers/reference/content-hash
//
// Should new be nil, DefaultHash.New will be used in its place.
//
// Should blockSize be less than 1, DefaultBlockSize will be used instead.
//
// Hashes returned by New are NOT safe for concurrent use by multiple
// goroutines.
func New(new func() hash.Hash, blockSize int) hash.Hash {
	if new == nil {
		new = DefaultHash.New
	}
	if blockSize < 1 {
		blockSize = DefaultBlockSize
	}

	h := new()
	return &wrapper{
		blk:       h,
		sum:       new(),
		buf:       make([]byte, h.Size()),
		blockSize: blockSize,
		rem:       blockSize,
	}
}

type wrapper struct {
	buf       []byte
	blk       hash.Hash // sums blocks
	sum       hash.Hash // sums blocks' sums
	blockSize int       // denotes the size of individual blocks
	rem       int       // room in current block
}

func (w *wrapper) BlockSize() int {
	return w.blockSize
}

func (w *wrapper) Size() int {
	return w.blk.Size()
}

func (w *wrapper) Sum(b []byte) []byte {
	if w.rem != w.BlockSize() {
		w.sumBlock()
	}

	return w.sum.Sum(b)
}

func (w *wrapper) Reset() {
	w.blk.Reset()
	w.sum.Reset()
	w.rem = w.BlockSize()
}

func (w *wrapper) Write(b []byte) (n int, err error) {
	for err == nil && len(b) > 0 {
		// read as many bytes as are missing from the current block
		// (or until the end of b)
		var nn int
		nn, err = w.blk.Write(b[:min(len(b), w.rem)])
		n += nn
		b = b[nn:]

		if w.rem -= nn; w.rem == 0 {
			w.sumBlock()

			w.rem = w.BlockSize()
		}
	}

	return
}

// sumBlock sums w's blk hash into w's buf.
func (w *wrapper) sumBlock() {
	defer w.blk.Reset()

	_ = w.blk.Sum(w.buf[:0])
	_, _ = w.sum.Write(w.buf)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
