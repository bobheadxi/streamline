package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	t.Run("active", func(t *testing.T) {
		p := Filter(func(line []byte) bool { return false })
		assert.False(t, p.Inactive())

		line, err := p.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.Empty(t, line)
	})

	t.Run("inactive", func(t *testing.T) {
		p := Filter(nil)
		assert.True(t, p.Inactive())

		line, err := p.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.Equal(t, string(line), "foo")
	})
}