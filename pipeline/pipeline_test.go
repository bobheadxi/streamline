package pipeline

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiPipeline(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		p := MultiPipeline{}
		line, err := p.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.NotNil(t, line)
	})

	t.Run("all active", func(t *testing.T) {
		t.Parallel()

		p := MultiPipeline{
			Filter(func(line []byte) bool { return false }),
			Filter(func(line []byte) bool { return true }),
		}

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

		line, err := p.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.Empty(t, line)
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

		line, err := p.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.Nil(t, line)

		assert.True(t, map2called)
		assert.False(t, map3called)
	})

	t.Run("error handling", func(t *testing.T) {
		t.Parallel()

		var map2called bool
		p := MultiPipeline{
			MapErr(func(line []byte) ([]byte, error) { return line, errors.New("oh no") }),
			Map(func(line []byte) []byte {
				map2called = true
				return nil
			}),
		}

		line, err := p.ProcessLine([]byte("foo"))
		assert.Error(t, err)
		assert.NotEmpty(t, line)

		assert.False(t, map2called)
	})
}
