package pipe

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.bobheadxi.dev/streamline/internal/testdata"
)

func TestStream(t *testing.T) {
	w, s := NewStream()

	input, size, _ := testdata.GenerateLargeInput(10)
	go func() {
		_, err := io.Copy(w, input)
		require.NoError(t, err)
		w.CloseWithError(nil)
	}()

	var read strings.Builder
	err := s.Stream(func(line string) error {
		_, err := read.WriteString(line + "\n")
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, size, read.Len()-1) // we append an extra newline
}
