package benchmarks

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.bobheadxi.dev/streamline"
	"go.bobheadxi.dev/streamline/pipeline"
)

func TestStreamReadAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping io.ReadAll(*streamline.Stream) load test")
	}

	input, _ := generateLargeInput()
	s := streamline.New(input).
		WithPipeline(pipeline.Map(func(line []byte) []byte { return line }))

	start := time.Now()
	n, err := io.ReadAll(s)
	assert.NoError(t, err)
	t.Logf("read %d bytes in %s", len(n), time.Since(start))
}
