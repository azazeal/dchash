package dchash

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"io"
	mand "math/rand"
	"os"
	"path/filepath"
	"testing"

	_ "crypto/sha256"
)

var sums = map[string][]byte{
	"milky-way-nasa.jpg": decode("485291fa0ee50c016982abbfa943957bcd231aae0492ccbaa22c58e3997b35e0"),
}

func TestConstants(t *testing.T) {
	if exp := 1 << 22; DefaultBlockSize != exp {
		t.Errorf("expected DefaultBlockSize to be %d, got %d", exp, DefaultBlockSize)
	}
}

func Test(t *testing.T) {
	h := New(nil, -1)
	s := make([]byte, h.Size())

	_, _ = h.Write([]byte("does Reset work?"))

	for name, exp := range sums {
		h.Reset()

		readTestDataFile(t, h, name)

		if got := h.Sum(s[:0]); !bytes.Equal(got, exp) {
			t.Errorf("exp %x, got %x (file: %s)", exp, got, name)
		}
	}
}

func readTestDataFile(t *testing.T, w io.Writer, name string) {
	t.Helper()

	path := filepath.Join("testdata", name)
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed opening %s: %v", path, err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			t.Fatalf("failed closing %s: %v", name, err)
		}
	}()

	if _, err := io.CopyBuffer(w, f, makeBuf(t)); err != nil {
		t.Fatalf("failed copying %s: %v", name, err)
	}
}

func makeBuf(tb testing.TB) []byte {
	tb.Helper()

	src := mand.NewSource(seed(tb))
	rng := mand.New(src)

	return make([]byte, 1+rng.Intn(1<<8-1))
}

func seed(tb testing.TB) int64 {
	tb.Helper()

	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		tb.Fatalf("failed reading seed: %v", err)
	}

	return int64(binary.BigEndian.Uint64(b))
}

func decode(s string) []byte {
	v, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return v
}
