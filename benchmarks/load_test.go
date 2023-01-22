package benchmarks

import (
	"io"
	"os"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.bobheadxi.dev/streamline"
	"go.bobheadxi.dev/streamline/internal/testdata"
	"go.bobheadxi.dev/streamline/pipeline"
)

func TestStreamReadAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping io.ReadAll(*streamline.Stream) load test")
	}

	input, size, _ := testdata.GenerateLargeInput(100)
	s := streamline.New(input).
		WithPipeline(pipeline.Map(func(line []byte) []byte { return line }))

	start := time.Now()

	if *profile {
		f, err := os.OpenFile(t.Name(), os.O_RDWR|os.O_CREATE, 0600)
		require.NoError(t, err)
		pprof.StartCPUProfile(f)
	}

	n, err := io.ReadAll(s)

	if *profile {
		pprof.StopCPUProfile()
	}

	assert.NoError(t, err)
	assert.Equal(t, size+1, len(n)) // we currently add trailing newline
	t.Logf("read %d bytes in %s", len(n), time.Since(start))
}
