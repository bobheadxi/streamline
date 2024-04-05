package pipeline_test

import (
	"strings"
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.bobheadxi.dev/streamline"
	"go.bobheadxi.dev/streamline/pipeline"
)

func TestSample(t *testing.T) {
	t.Run("sample stream", func(t *testing.T) {
		t.Parallel()

		stream := streamline.New(strings.NewReader("foo bar baz\nbaz bar\nhello world\ngoodbye world")).
			WithPipeline(pipeline.Sample(2))

		lines, err := stream.Lines()
		require.NoError(t, err)
		autogold.Expect([]string{"baz bar", "goodbye world"}).Equal(t, lines)
	})

	t.Run("never sample", func(t *testing.T) {
		t.Parallel()

		s := pipeline.Sample(0)
		line, err := s.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.Nil(t, line)

	})

	t.Run("always sample", func(t *testing.T) {
		t.Parallel()

		s := pipeline.Sample(1)
		line, err := s.ProcessLine([]byte("foo"))
		assert.NoError(t, err)
		assert.NotNil(t, line)
	})

	t.Run("invalid sampling", func(t *testing.T) {
		t.Parallel()

		s := pipeline.Sample(-1)
		line, err := s.ProcessLine([]byte("foo"))
		assert.Error(t, err)
		assert.Nil(t, line)
	})

	t.Run("sample 3rd", func(t *testing.T) {
		t.Parallel()

		s := pipeline.Sample(3)

		line, err := s.ProcessLine([]byte("foo")) // 1
		assert.NoError(t, err)
		assert.Nil(t, line)

		line, err = s.ProcessLine([]byte("bar")) // 2
		assert.NoError(t, err)
		assert.Nil(t, line)

		line, err = s.ProcessLine([]byte("baz")) // 3
		assert.NoError(t, err)
		assert.Equal(t, "baz", string(line))

		line, err = s.ProcessLine([]byte("dropped")) // 4
		assert.NoError(t, err)
		assert.Nil(t, line)
	})
}
