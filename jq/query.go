package jq

import (
	"bytes"
	"context"
	"io"
)

// Query is a utility for building and executing a JQ query against some data, such as a
// streamline.Stream instance.
//
// Internally, Query uses github.com/itchyny/gojq to build and run the query.
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
