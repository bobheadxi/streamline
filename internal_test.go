package streamline

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
