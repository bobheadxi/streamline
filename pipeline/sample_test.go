package pipeline_test

import (
	"strings"
	"testing"

	"github.com/hexops/autogold"
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
		autogold.Want("sampled 2", []string{"baz bar", "goodbye world"}).Equal(t, lines)
	})

	t.Run("active", func(t *testing.T) {
		t.Parallel()

		s := pipeline.Sample(2)
		assert.False(t, s.Inactive())
	})

	t.Run("inactive", func(t *testing.T) {
		t.Parallel()

		s := pipeline.Sample(0)
		assert.True(t, s.Inactive())

		s = pipeline.Sample(1)
		assert.True(t, s.Inactive())

		s = pipeline.Sample(-1)
		assert.True(t, s.Inactive())
	})
}
