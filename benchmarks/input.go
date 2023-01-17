package benchmarks

import (
	"io"
	"strings"
)

var testData = []string{
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
	"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
	"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
}

const inputLineCount = 100000

func generateInput() (io.Reader, func()) {
	inputLines := make([]string, inputLineCount)
	for l := 0; l < inputLineCount; l++ {
		inputLines[l] = testData[l%len(testData)]
	}
	r := strings.NewReader(strings.Join(inputLines, "\n"))
	return r, func() { r.Seek(0, 0) }
}

func generateLargeInput() (io.Reader, func()) {
	wantCount := inputLineCount * 100
	inputLines := make([]string, wantCount)
	for l := 0; l < wantCount; l++ {
		inputLines[l] = testData[l%len(testData)]
	}
	r := strings.NewReader(strings.Join(inputLines, "\n"))
	return r, func() { r.Seek(0, 0) }
}
