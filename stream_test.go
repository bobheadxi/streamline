package streamline_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/hexops/autogold"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.bobheadxi.dev/streamline"
	"go.bobheadxi.dev/streamline/pipeline"
)

func TestStream(t *testing.T) {
	s := streamline.New(strings.NewReader("foo\nbar\nbaz"))

	lines, err := s.Lines()
	require.NoError(t, err)
	autogold.Want("collected lines", []string{"foo", "bar", "baz"}).Equal(t, lines)
}

func TestStreamReader(t *testing.T) {
	t.Run("io.Read", func(t *testing.T) {
		stream := streamline.New(strings.NewReader("foo bar baz\nbaz bar\nhello world")).
			WithPipeline(pipeline.Map(func(line []byte) []byte {
				return bytes.ReplaceAll(line, []byte{' '}, []byte{'-'})
			}))

		p := make([]byte, 5)
		n, err := stream.Read(p)
		assert.NoError(t, err)
		assert.NotZero(t, n)
		autogold.Want("chunk 1", "foo-b").Equal(t, string(p[:n]))

		p = make([]byte, 5)
		n, err = stream.Read(p)
		assert.NoError(t, err)
		assert.NotZero(t, n)
		autogold.Want("chunk 2", "ar-ba").Equal(t, string(p[:n]))

		p = make([]byte, 5)
		n, err = stream.Read(p)
		assert.NoError(t, err)
		assert.NotZero(t, n)
		autogold.Want("chunk 3", "z\n").Equal(t, string(p[:n]))

		all, err := io.ReadAll(stream)
		assert.NoError(t, err)
		autogold.Want("all remaining", "baz-bar\nhello-world\n").Equal(t, string(all))
	})

	t.Run("io.ReadAll", func(t *testing.T) {
		stream := streamline.New(strings.NewReader("foo bar baz\nbaz bar\nhello world")).
			WithPipeline(pipeline.Map(func(line []byte) []byte {
				return bytes.ReplaceAll(line, []byte{' '}, []byte{'-'})
			}))

		all, err := io.ReadAll(stream)
		assert.NoError(t, err)
		autogold.Want("collected all", "foo-bar-baz\nbaz-bar\nhello-world\n").Equal(t, string(all))
	})
}
