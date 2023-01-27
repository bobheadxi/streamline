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

// TestLoad should test the exact matching of large inputs against the outputs.
func TestLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping load test")
	}

	input, _, reset := testdata.GenerateLargeInput(100)
	inputData, _ := io.ReadAll(input)
	reset()

	b, err := streamline.New(input).
		WithPipeline(pipeline.Map(func(line []byte) []byte { return line })).
		Bytes()

	assert.NoError(t, err)
	assert.Equal(t, string(inputData), string(b))
}

// TestLoadProfile should only do basic checks, focusing on generating profiles
// of large input loads.
func TestLoadProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping load tests")
	}

	if *profile {
		t.Logf("generating profiles in %q - view with 'go tool pprof -png ./benchmarks/%s/TEST'",
			t.Name(), t.Name())
		_ = os.MkdirAll(t.Name(), os.ModePerm)
	}

	for _, tc := range []struct {
		name string
		test func(s *streamline.Stream) (int, error)
	}{
		{
			name: "io.ReadAll",
			test: func(s *streamline.Stream) (int, error) {
				b, err := io.ReadAll(s)
				return len(b), err
			},
		},
		{
			name: "Stream",
			test: func(s *streamline.Stream) (int, error) {
				var n int
				return n, s.Stream(func(line string) { n += len(line) + 1 })
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			input, size, _ := testdata.GenerateLargeInput(100)
			s := streamline.New(input).
				WithPipeline(pipeline.Map(func(line []byte) []byte { return line }))

			start := time.Now()

			if *profile {
				f, err := os.OpenFile(t.Name(), os.O_RDWR|os.O_CREATE, 0600)
				require.NoError(t, err)
				defer f.Close()
				pprof.StartCPUProfile(f)
			}

			n, err := tc.test(s)

			if *profile {
				pprof.StopCPUProfile()
			}

			assert.NoError(t, err)
			assert.Equal(t, size+1, n) // we currently add trailing newline
			t.Logf("read %d bytes in %s", n, time.Since(start))
		})
	}
}
