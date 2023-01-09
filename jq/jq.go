package jq

import (
	"bytes"
	"fmt"
	"io"

	"go.bobheadxi.dev/streamline/pipeline"
)

// Pipeline builds a JQ query for a pipeline that runs the query against each line and
// maps the result to the output.
func Pipeline(query string) pipeline.Pipeline {
	jqCode, err := buildJQ(query)
	if err != nil {
		return pipeline.Map(func(line []byte) ([]byte, error) { return nil, err })
	}

	return pipeline.Map(func(line []byte) ([]byte, error) {
		if len(line) == 0 {
			return line, nil
		}
		result, err := execJQ(jqCode, bytes.NewReader(line))
		if err != nil {
			// Embed the consumed content
			return nil, fmt.Errorf("%w: %s", err, string(line))
		}
		return result, nil
	})
}

// Query is a utility for building and executing a JQ query against some data, such as a
// streamline.Stream instance.
func Query(data io.Reader, query string) ([]byte, error) {
	jqCode, err := buildJQ(query)
	if err != nil {
		return nil, err
	}

	return execJQ(jqCode, data)
}
