package testdata

import (
	"bytes"
	"io"
	"strings"
)

var testData = []string{
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	"Short.",
	"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
	"Also short.",
	"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
	"Brief.",
	"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
}

const inputLineCount = 100000

// GenerateInput creates arbitrary input with inputLineCount lines.
func GenerateInput() (io.Reader, func()) {
	inputLines := make([]string, inputLineCount)
	for l := 0; l < inputLineCount; l++ {
		inputLines[l] = testData[l%len(testData)]
	}
	r := bytes.NewReader([]byte(strings.Join(inputLines, "\n")))
	return r, func() { r.Seek(0, 0) }
}

// GenerateLargeInput creates arbitrary input with inputLineCount * x lines.
func GenerateLargeInput(x int) (io.Reader, int, func()) {
	wantCount := inputLineCount * x
	inputLines := make([]string, wantCount)
	for l := 0; l < wantCount; l++ {
		inputLines[l] = testData[l%len(testData)]
	}
	data := strings.Join(inputLines, "\n")
	r := bytes.NewReader([]byte(data))
	return r, len(data), func() { r.Seek(0, 0) }
}

// GenerateWideInput creates a single-line input with inputLineCount * x characters.
func GenerateWideInput(x int) (io.Reader, int, func()) {
	wantWidth := inputLineCount * x
	var inputData string
	for len(inputData) < wantWidth {
		inputData += (" " + testData[len(inputData)%len(testData)])
	}
	data := strings.Join([]string{inputData}, "\n")
	r := strings.NewReader(data)
	return r, len(inputData), func() { r.Seek(0, 0) }
}
