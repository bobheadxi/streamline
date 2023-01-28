package pipe

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.bobheadxi.dev/streamline/internal/testdata"
)

func TestStream(t *testing.T) {
	t.Parallel()

	w, s := NewStream()

	input, size, _ := testdata.GenerateLargeInput(10)
	go func() {
		_, err := io.Copy(w, input)
		require.NoError(t, err)
		w.CloseWithError(nil)
	}()

	var read strings.Builder
	err := s.Stream(func(line string) { read.WriteString(line + "\n") })
	assert.NoError(t, err)
	assert.Equal(t, size, read.Len()-1) // we append an extra newline
}

func TestBoundedStream(t *testing.T) {
	t.Parallel()

	w, s := NewBoundedStream()

	input, size, _ := testdata.GenerateLargeInput(1)
	go func() {
		_, err := io.Copy(w, input)
		require.NoError(t, err)
		w.CloseWithError(nil)
	}()

	var read strings.Builder
	err := s.Stream(func(line string) { read.WriteString(line + "\n") })
	assert.NoError(t, err)
	assert.Equal(t, size, read.Len()-1) // we append an extra newline
}

func TestWriterErrorCloser(t *testing.T) {
	t.Parallel()

	w, s := NewStream()
	_, err := w.Write([]byte("foo\nbar"))
	assert.NoError(t, err)

	_ = w.CloseWithError(errors.New("oh no!"))
	v, err := s.String()
	assert.Equal(t, "foo\nbar", v)
	require.Error(t, err)
	assert.Equal(t, "oh no!", err.Error())
}
