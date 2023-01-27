package pipeline_test

import (
	"strings"
	"testing"

	"github.com/hexops/autogold"
	"github.com/stretchr/testify/require"
	"go.bobheadxi.dev/streamline"
	"go.bobheadxi.dev/streamline/pipeline"
)

func TestSample(t *testing.T) {
	stream := streamline.New(strings.NewReader("foo bar baz\nbaz bar\nhello world\ngoodbye world")).
		WithPipeline(pipeline.Sample(2))

	lines, err := stream.Lines()
	require.NoError(t, err)
	autogold.Want("sampled 2", []string{"baz bar", "goodbye world"}).Equal(t, lines)
}
