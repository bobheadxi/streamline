package pipeline

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Run("active", func(t *testing.T) {
		p := Map(func(line []byte) []byte { return nil })
		assert.False(t, p.Inactive())

		line, err := p.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.Empty(t, line)
	})

	t.Run("inactive", func(t *testing.T) {
		p := Map(nil)
		assert.True(t, p.Inactive())

		line, err := p.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.Equal(t, string(line), "foo")
	})
}

func TestErrMap(t *testing.T) {
	t.Run("active", func(t *testing.T) {
		p := MapErr(func(line []byte) ([]byte, error) { return line, errors.New("foo") })
		assert.False(t, p.Inactive())

		line, err := p.ProcessLine([]byte("foo"))
		assert.Error(t, err)
		assert.Equal(t, string(line), "foo")
	})

	t.Run("inactive", func(t *testing.T) {
		p := MapErr(nil)
		assert.True(t, p.Inactive())

		line, err := p.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.Equal(t, string(line), "foo")
	})
}
