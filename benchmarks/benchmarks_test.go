package benchmarks

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.bobheadxi.dev/streamline"
	"go.bobheadxi.dev/streamline/pipeline"
)

var testData = []string{
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
	"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
	"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
}

const inputLineCount = 100000

func generateInput() io.Reader {
	inputLines := make([]string, inputLineCount)
	for l := 0; l < inputLineCount; l++ {
		inputLines[l] = testData[l%len(testData)]
	}
	return strings.NewReader(strings.Join(inputLines, "\n"))
}

func BenchmarkBufioReader(b *testing.B) {
	input := generateInput()

	for i := 0; i < b.N; i++ {
		r := bufio.NewReader(input)
		for {
			_, err := r.ReadBytes('\n')
			if errors.Is(err, io.EOF) {
				break
			}
			assert.NoError(b, err)
		}
	}
}

func BenchmarkStreamBytes(b *testing.B) {
	input := generateInput()

	for i := 0; i < b.N; i++ {
		s := streamline.New(input)
		err := s.StreamBytes(func(_ []byte) error { return nil })
		assert.NoError(b, err)
	}
}

func BenchmarkStreamBytesWithPipeline(b *testing.B) {
	input := generateInput()

	for i := 0; i < b.N; i++ {
		s := streamline.New(input).
			WithPipeline(pipeline.Map(func(line []byte) []byte { return line }))
		err := s.StreamBytes(func(_ []byte) error { return nil })
		assert.NoError(b, err)
	}
}
