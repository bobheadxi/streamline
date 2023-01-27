package benchmarks

import (
	"bufio"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.bobheadxi.dev/streamline"
	"go.bobheadxi.dev/streamline/internal/testdata"
	"go.bobheadxi.dev/streamline/pipeline"
)

func BenchmarkBufioReader(b *testing.B) {
	input, reset := testdata.GenerateInput()

	for i := 0; i < b.N; i++ {
		reset()

		// assert we are benchmarking against a bare LineReader
		var r streamline.LineReader = bufio.NewReader(input)
		for {
			_, err := r.ReadSlice('\n')
			if errors.Is(err, io.EOF) {
				break
			}
			assert.NoError(b, err)
		}
	}
}

func BenchmarkBufioScannerBytes(b *testing.B) {
	input, reset := testdata.GenerateInput()

	for i := 0; i < b.N; i++ {
		reset()

		s := bufio.NewScanner(input)
		for s.Scan() {
			_ = s.Bytes()
		}
		assert.NoError(b, s.Err())
	}
}

func BenchmarkBufioScannerText(b *testing.B) {
	input, reset := testdata.GenerateInput()

	for i := 0; i < b.N; i++ {
		reset()

		s := bufio.NewScanner(input)
		for s.Scan() {
			_ = s.Text()
		}
		assert.NoError(b, s.Err())
	}
}

func BenchmarkStream(b *testing.B) {
	input, reset := testdata.GenerateInput()

	for i := 0; i < b.N; i++ {
		reset()

		s := streamline.New(input)
		err := s.Stream(func(_ string) {})
		assert.NoError(b, err)
	}
}

func BenchmarkStreamBytes(b *testing.B) {
	input, reset := testdata.GenerateInput()

	for i := 0; i < b.N; i++ {
		reset()

		s := streamline.New(input)
		err := s.StreamBytes(func(_ []byte) error { return nil })
		assert.NoError(b, err)
	}
}

func BenchmarkStreamBytesWithPipeline(b *testing.B) {
	input, reset := testdata.GenerateInput()

	for i := 0; i < b.N; i++ {
		reset()

		s := streamline.New(input).
			WithPipeline(pipeline.Map(func(line []byte) []byte { return line }))
		err := s.StreamBytes(func(_ []byte) error { return nil })
		assert.NoError(b, err)
	}
}

func BenchmarkStreamLines(b *testing.B) {
	input, reset := testdata.GenerateInput()

	for i := 0; i < b.N; i++ {
		reset()

		s := streamline.New(input)
		_, err := s.Lines()
		assert.NoError(b, err)
	}
}

func BenchmarkStreamLinesWithPipeline(b *testing.B) {
	input, reset := testdata.GenerateInput()

	for i := 0; i < b.N; i++ {
		reset()

		s := streamline.New(input).
			WithPipeline(pipeline.Map(func(line []byte) []byte { return line }))
		_, err := s.Lines()
		assert.NoError(b, err)
	}
}

func BenchmarkIoReadAll(b *testing.B) {
	input, reset := testdata.GenerateInput()

	for i := 0; i < b.N; i++ {
		reset()

		_, err := io.ReadAll(input)
		assert.NoError(b, err)
	}
}

func BenchmarkStreamReadAll(b *testing.B) {
	input, reset := testdata.GenerateInput()

	for i := 0; i < b.N; i++ {
		reset()

		s := streamline.New(input)
		_, err := io.ReadAll(s)
		assert.NoError(b, err)
	}
}

func BenchmarkStreamReadAllWithPipeline(b *testing.B) {
	input, reset := testdata.GenerateInput()

	for i := 0; i < b.N; i++ {
		reset()

		s := streamline.New(input).
			WithPipeline(pipeline.Map(func(line []byte) []byte { return line }))
		_, err := io.ReadAll(s)
		assert.NoError(b, err)
	}
}
