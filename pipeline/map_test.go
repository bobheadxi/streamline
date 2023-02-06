package pipeline

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	t.Run("active", func(t *testing.T) {
		t.Parallel()

		p := Map(func(line []byte) []byte { return nil })

		line, err := p.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.Empty(t, line)
	})
}

func TestErrMap(t *testing.T) {
	t.Run("active", func(t *testing.T) {
		t.Parallel()

		p := MapErr(func(line []byte) ([]byte, error) { return line, errors.New("foo") })

		line, err := p.ProcessLine([]byte("foo"))
		assert.Error(t, err)
		assert.Equal(t, string(line), "foo")
	})
}
