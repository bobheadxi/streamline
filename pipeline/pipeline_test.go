package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiPipeline(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		p := MultiPipeline{}
		assert.True(t, p.Inactive())

		p = nil
		assert.True(t, p.Inactive())
	})

	t.Run("all active", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

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
		t.Parallel()

		p := MultiPipeline{
			Filter(nil),
			Filter(nil),
		}
		assert.True(t, p.Inactive())

		line, err := p.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.Equal(t, string(line), "foo")
	})

	t.Run("line skipping", func(t *testing.T) {
		t.Parallel()

		var map2called, map3called bool
		p := MultiPipeline{
			// Return empty slice, next pipeline should run
			Map(func(line []byte) []byte { return []byte{} }),
			// Return nil slice, next pipeline should not run
			Map(func(line []byte) []byte {
				map2called = true
				return nil
			}),
			Map(func(line []byte) []byte {
				map3called = true
				return line
			}),
		}
		assert.False(t, p.Inactive())

		line, err := p.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.Nil(t, line)

		assert.True(t, map2called)
		assert.False(t, map3called)
	})
}
