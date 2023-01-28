package jq

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"go.bobheadxi.dev/streamline/pipeline"
)

// Pipeline builds a JQ query for a pipeline that runs the query against each line and
// maps the result to the output.
//
// Internally, it uses github.com/itchyny/gojq
func Pipeline(query string) pipeline.Pipeline {
	return PipelineContext(context.Background(), query)
}

// PipelineContext is the same as Pipeline, but runs the generated JQ code in the given
// context.
func PipelineContext(ctx context.Context, query string) pipeline.Pipeline {
	jqCode, err := buildJQ(query)
	if err != nil {
		return pipeline.MapErr(func(line []byte) ([]byte, error) { return nil, err })
	}

	return pipeline.MapErr(func(line []byte) ([]byte, error) {
		if len(line) == 0 {
			return line, nil
		}
		result, err := execJQ(ctx, jqCode, line)
		if err != nil {
			// Embed the consumed content
			return nil, fmt.Errorf("%w: %s", err, string(line))
		}
		return result, nil
	})
}

// Query is a utility for building and executing a JQ query against some data, such as a
// streamline.Stream instance.
//
// Internally, it uses github.com/itchyny/gojq
func Query(data io.Reader, query string) ([]byte, error) {
	return QueryContext(context.Background(), data, query)
}

// QueryContext is the same as Query, but runs the generated JQ code in the given context.
func QueryContext(ctx context.Context, data io.Reader, query string) ([]byte, error) {
	jqCode, err := buildJQ(query)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	if _, err := io.Copy(&b, data); err != nil {
		return nil, err
	}

	return execJQ(ctx, jqCode, b.Bytes())
}
