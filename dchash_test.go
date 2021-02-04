package dchash

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"io"
	mand "math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	_ "crypto/sha256"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var wellKnownSums = []struct {
	file   string
	offset int64
	size   int64
	exp    string
}{
	0: {
		file: "milky-way-nasa.jpg",
		size: DefaultBlockSize,
		exp:  "5fe03e2b2d4bbc75804fab4423cffe766f1b2d2080ad38c43d505e6e8ffc1344",
	},
	1: {
		file:   "milky-way-nasa.jpg",
		offset: DefaultBlockSize,
		size:   DefaultBlockSize,
		exp:    "299b77bc0ba00eba12dbf899ddec86d53312aca8849e7121dab23cbfddff3fed",
	},
	2: {
		file:   "milky-way-nasa.jpg",
		offset: 8388608,
		size:   1322815,
		exp:    "0e7baa366f72ac449f8904c45151bf6a2fc13dd07be47c7da6cbc3cebb437111",
	},
	3: {
		file: "milky-way-nasa.jpg",
		size: 9711423,
		exp:  "485291fa0ee50c016982abbfa943957bcd231aae0492ccbaa22c58e3997b35e0",
	},
}

func TestWellKnownSums(t *testing.T) {
	h := New(nil, -1)
	got := make([]byte, h.Size())

	_, _ = h.Write([]byte("does Reset work?"))

	for i := range wellKnownSums {
		kase := wellKnownSums[i]

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			h.Reset()

			readPart(t, h, kase.file, kase.offset, kase.size)
			_ = h.Sum(got[:0])

			assert.Equal(t, dec(kase.exp), got)
		})
	}
}

func readPart(t *testing.T, w io.Writer, name string, offset, size int64) {
	t.Helper()

	path := filepath.Join("testdata", name)
	f, err := os.Open(path)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, f.Close())
	}()

	_, err = f.Seek(offset, io.SeekStart)
	require.NoError(t, err)

	src := io.LimitReader(f, size)
	_, err = io.CopyBuffer(w, src, makeBuf(t))
	require.NoError(t, err)
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
	_, err := rand.Read(b)
	require.NoError(tb, err)

	return int64(binary.BigEndian.Uint64(b))
}

func dec(s string) []byte {
	v, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return v
}
