package streamline

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.bobheadxi.dev/streamline/internal/testdata"
)

func TestNew(t *testing.T) {
	t.Run("with LineReader implementation", func(t *testing.T) {
		t.Parallel()

		lr := bufio.NewReader(strings.NewReader("hello world"))
		s := New(lr)
		assert.Equal(t, lr, s.reader)
	})

	t.Run("without LineReader implementation", func(t *testing.T) {
		t.Parallel()

		r := strings.NewReader("hello world")
		s := New(r)
		assert.NotEqual(t, r, s.reader)
	})
}

func TestStreamLoad(t *testing.T) {
	t.Run("large input", func(t *testing.T) {
		t.Parallel()

		input, size, _ := testdata.GenerateLargeInput(10)
		s := New(input)

		data, err := s.Bytes()
		assert.NoError(t, err)
		assert.Equal(t, size, len(data))
	})

	t.Run("wide input", func(t *testing.T) {
		t.Parallel()

		input, size, _ := testdata.GenerateWideInput(10)
		s := New(input)

		data, err := s.Bytes()
		assert.NoError(t, err)
		assert.Equal(t, size, len(data))
	})
}
