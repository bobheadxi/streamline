package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiPipeline(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		p := MultiPipeline{}
		assert.True(t, p.Inactive())

		p = nil
		assert.True(t, p.Inactive())
	})

	t.Run("all active", func(t *testing.T) {
		p := MultiPipeline{
			Filter(func(line []byte) bool { return false }),
			Filter(func(line []byte) bool { return true }),
		}
		assert.False(t, p.Inactive())

		line, err := p.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.Empty(t, line)
	})

	t.Run("some active", func(t *testing.T) {
		p := MultiPipeline{
			Filter(func(line []byte) bool { return false }),
			Filter(nil),
			Filter(func(line []byte) bool { return true }),
		}
		assert.False(t, p.Inactive())

		line, err := p.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.Empty(t, line)
	})

	t.Run("all inactive", func(t *testing.T) {
		p := MultiPipeline{
			Filter(nil),
			Filter(nil),
		}
		assert.True(t, p.Inactive())

		line, err := p.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.Equal(t, string(line), "foo")
	})
}
